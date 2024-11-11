from locust import HttpUser, task, events
from locust.env import Environment
import time
import logging
from locust.runners import LocalRunner
from locust.stats import stats_printer
import gevent

logging.basicConfig(level=logging.INFO)

def increasing_wait_time(user):
    current_time = time.time()
    elapsed_since_last_increase = current_time - user.last_increase

    if elapsed_since_last_increase >= user.increase_interval:
        user.current_wait_time = max(user.current_wait_time * 0.5, 0.1)  # Minimum wait time
        user.last_increase = current_time
    return user.current_wait_time

class JavaUser(HttpUser):
    host = "http://localhost:8091"
    current_wait_time = 1  # Initial wait time
    last_increase = 0
    increase_interval = 1

    def wait_time(self):
        return increasing_wait_time(self)

    @task
    def ping_java(self):
        print("Executing Java task")
        try:
            with self.client.get("/actuator/prometheus", name="Java Ping", catch_response=True) as response:
                if response.status_code != 200:
                    response.failure(f"Java ping failed with status code: {response.status_code}")
                elif response.elapsed.total_seconds() > 1:
                    response.failure(f"Java ping too slow: {response.elapsed.total_seconds()}s")
        except Exception as e:
            print(f"Java request failed: {str(e)}")

class GinUser(HttpUser):
    host = "http://localhost:8090"
    current_wait_time = 1
    last_increase = 0
    increase_interval = 1

    def wait_time(self):
        return increasing_wait_time(self)

    @task
    def ping_gin(self):
        print("Gin ping")
        try:
            with self.client.get("/metrics", name="Gin Ping", catch_response=True) as response:
                if response.status_code != 200:
                    response.failure(f"Gin ping failed with status code: {response.status_code}")
                elif response.elapsed.total_seconds() > 1:
                    response.failure(f"Gin ping too slow: {response.elapsed.total_seconds()}s")
        except Exception as e:
            print(f"Gin request failed: {str(e)}")

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print(f"Load test starting at {time.strftime('%Y-%m-%d %H:%M:%S')}")
    print("Press Ctrl+C to stop the test")

@events.request.add_listener
def on_request(request_type, name, response_time, response_length, response, **kwargs):
    if response_time > 1000:
        print(f"{time.strftime('%Y-%m-%d %H:%M:%S')} - Warning: Slow request to {name} completed in {response_time}ms")

@events.quitting.add_listener
def on_quitting(environment, **kwargs):
    print(f"\nLoad test finished at {time.strftime('%Y-%m-%d %H:%M:%S')}")

if __name__ == "__main__":
    # Create separate Environment instances for each user type
    java_env = Environment(user_classes=[JavaUser])
    gin_env = Environment(user_classes=[GinUser])

    # Create separate runners for each environment
    java_runner = LocalRunner(environment=java_env)
    gin_runner = LocalRunner(environment=gin_env)

    # Start greenlets for stats printing for each environment
    gevent.spawn(stats_printer(java_env.stats))
    gevent.spawn(stats_printer(gin_env.stats))

    # Start each runner with its own user count and spawn rate
    java_runner.start(5, spawn_rate=1)
    gin_runner.start(5, spawn_rate=1)

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nStopping load test...")
    finally:
        java_runner.quit()
        gin_runner.quit()
