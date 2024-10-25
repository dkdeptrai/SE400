from locust import HttpUser, task, constant, events
import time
import math
import gevent
from locust.env import Environment
from locust.stats import stats_printer, stats_history
from locust.log import setup_logging
import logging
from locust.exception import RescheduleTask

class LoadTestUser(HttpUser):
    abstract = True
    
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.start_time = time.time()
        self.initial_wait = 1.0
        self.max_rps = 400  # Maximum requests per second per user
    
    def wait_time(self):
        # Calculate minutes elapsed since start
        minutes_elapsed = (time.time() - self.start_time) / 60
        
        # Calculate the new wait time with a gentler increase (using 1.5 instead of 2)
        # Set a minimum wait time based on max_rps
        target_wait = self.initial_wait / math.pow(2, minutes_elapsed)
        min_wait = 1.0 / self.max_rps
        new_wait = max(target_wait, min_wait)
        
        # Log current rate every 15 seconds
        if int(time.time()) % 15 == 0:
            current_rate = 1/new_wait if new_wait > 0 else "∞"
            print(f"\nCurrent target rate per user: {current_rate:.2f} requests/second")
            
        return new_wait

class FlaskUser(LoadTestUser):
    host = "http://localhost:5500"
    weight = 1
    
    @task(1)
    def ping_flask(self):
        try:
            with self.client.get("/ping", name="Flask Ping", catch_response=True) as response:
                if response.status_code != 200:
                    response.failure(f"Flask ping failed with status code: {response.status_code}")
                elif response.elapsed.total_seconds() > 1:  # If response takes more than 1 second
                    response.failure(f"Flask ping too slow: {response.elapsed.total_seconds()}s")
        except Exception as e:
            print(f"Flask request failed: {str(e)}")
            raise RescheduleTask()

class GinUser(LoadTestUser):
    host = "http://localhost:8090"
    weight = 1
    
    @task(1)
    def ping_gin(self):
        try:
            with self.client.get("/ping", name="Gin Ping", catch_response=True) as response:
                if response.status_code != 200:
                    response.failure(f"Gin ping failed with status code: {response.status_code}")
                elif response.elapsed.total_seconds() > 1:  # If response takes more than 1 second
                    response.failure(f"Gin ping too slow: {response.elapsed.total_seconds()}s")
        except Exception as e:
            print(f"Gin request failed: {str(e)}")
            raise RescheduleTask()

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
    
    # Create a test environment
    env = Environment(user_classes=[FlaskUser, GinUser])
    
    # Start the test
    env.create_local_runner()
    
    # Start a greenlet that periodically outputs the current stats
    gevent.spawn(stats_printer(env.stats))
    
    # Start a greenlet that saves current stats to history
    gevent.spawn(stats_history, env.runner)
    
    # Start with fewer users but maintain the ratio
    env.runner.start(10, spawn_rate=1)
    
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nStopping load test...")
    finally:
        env.runner.quit()