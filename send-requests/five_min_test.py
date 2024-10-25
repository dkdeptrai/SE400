from locust import HttpUser, task, constant, events
import time
import math
import gevent
from locust.env import Environment
from locust.stats import stats_printer, stats_history
from locust.log import setup_logging
import logging

class LoadTestUser(HttpUser):
    abstract = True
    
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.start_time = time.time()
        self.initial_wait = 1.0
    
    def wait_time(self):
        # Calculate minutes elapsed since start
        minutes_elapsed = (time.time() - self.start_time) / 60
        
        # Calculate the new wait time with 50% reduction each minute
        # Minimum wait time of 0.01 seconds to prevent overwhelming the system
        new_wait = max(self.initial_wait / math.pow(2, minutes_elapsed), 0.01)
        
        # Log current rate every 30 seconds
        if int(time.time()) % 30 == 0:
            current_rate = 1/new_wait if new_wait > 0 else "∞"
            print(f"\nCurrent target rate per user: {current_rate:.2f} requests/second")
            
        return new_wait

# Separate the services into different user classes with fixed ratios
class FlaskUser(LoadTestUser):
    host = "http://localhost:5500"
    weight = 1  # Equal weight with GinUser
    
    @task(1)
    def ping_flask(self):
        with self.client.get("/ping", name="Flask Ping", catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Flask ping failed with status code: {response.status_code}")

class GinUser(LoadTestUser):
    host = "http://localhost:8090"
    weight = 1  # Equal weight with FlaskUser
    
    @task(1)
    def ping_gin(self):
        with self.client.get("/ping", name="Gin Ping", catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Gin ping failed with status code: {response.status_code}")

# Custom event handlers
@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print(f"Load test starting at {time.strftime('%Y-%m-%d %H:%M:%S')}")
    print("Press Ctrl+C to stop the test")

@events.request.add_listener
def on_request(request_type, name, response_time, response_length, response, **kwargs):
    if response_time > 1000:  # Log only slow requests (>1000ms)
        print(f"{time.strftime('%Y-%m-%d %H:%M:%S')} - Warning: Slow request to {name} completed in {response_time}ms")

@events.quitting.add_listener
def on_quitting(environment, **kwargs):
    print(f"\nLoad test finished at {time.strftime('%Y-%m-%d %H:%M:%S')}")

if __name__ == "__main__":
    setup_logging("INFO", None)
    
    # Create a test environment with both user types
    env = Environment(user_classes=[FlaskUser, GinUser])
    
    # Start the test
    env.create_local_runner()
    
    # Start a greenlet that periodically outputs the current stats
    gevent.spawn(stats_printer(env.stats))
    
    # Start a greenlet that saves current stats to history
    gevent.spawn(stats_history, env.runner)
    
    # Start with 20 users (10 for each service) with spawn rate of 2 per second
    env.runner.start(20, spawn_rate=2)
    
    try:
        # Run indefinitely until interrupted
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nStopping load test...")
    finally:
        # Stop the runner
        env.runner.quit()