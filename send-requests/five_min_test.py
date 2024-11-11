from locust import HttpUser, task, constant_throughput, events
from locust.env import Environment
import time
import logging
from locust.runners import Runner, LocalRunner
from locust.stats import stats_printer, stats_history
import gevent

logging.basicConfig(level=logging.INFO)
class IncreasingRateUser(HttpUser):
    abstract = True
    # host = "http://localhost"
    increase_interval = 1

    def __init__(self, environment):
        super().__init__(environment)
        self.current_wait_time = 1  # Initial wait time
        self.last_increase = 0

    def wait_time(self):
        current_time = time.time()
        elapsed_since_last_increase = current_time - self.last_increase

        if elapsed_since_last_increase >= self.increase_interval:
            self.current_wait_time = self.current_wait_time * 0.1
            self.last_increase = current_time
        return self.current_wait_time

    # @task(1)
    # def ping_gin(self):
    #     try:
    #         with self.client.get(":8080/ping", name="Gin Ping", catch_response=True) as response:
    #             if response.status_code != 200:
    #                 response.failure(f"Gin ping failed with status code: {response.status_code}")
    #             elif response.elapsed.total_seconds() > 1:  # If response takes more than 1 second
    #                 response.failure(f"Gin ping too slow: {response.elapsed.total_seconds()}s")
    #     except Exception as e:
    #         print(f"Gin request failed: {str(e)}")
    
    # @task(0)
    # def ping_flask(self):
    #     try:
    #         with self.client.get(":5500/ping", name="Flask Ping", catch_response=True) as response:
    #             if response.status_code != 200:
    #                 response.failure(f"Flask ping failed with status code: {response.status_code}")
    #             elif response.elapsed.total_seconds() > 1:
    #                 response.failure(f"Flask ping too slow: {response.elapsed.total_seconds()}s")
    #     except Exception as e:
    #         print(f"Flask request failed: {str(e)}")
    

class GinUser(IncreasingRateUser):
    host = "http://localhost:8090"

    @task
    def ping_gin(self):
        try:
            with self.client.get("/ping", name="Gin Ping", catch_response=True) as response:
                if response.status_code != 200:
                    response.failure(f"Gin ping failed with status code: {response.status_code}")
                elif response.elapsed.total_seconds() > 1:  # If response takes more than 1 second
                    response.failure(f"Gin ping too slow: {response.elapsed.total_seconds()}s")
        except Exception as e:
            print(f"Gin request failed: {str(e)}")

class FlaskUser(IncreasingRateUser):
    host = "http://localhost:5500"

    @task
    def ping_flask(self):
        try:
            with self.client.get("/ping", name="Flask Ping", catch_response=True) as response:
                if response.status_code != 200:
                    response.failure(f"Flask ping failed with status code: {response.status_code}")
                elif response.elapsed.total_seconds() > 1:  # If response takes more than 1 second
                    response.failure(f"Flask ping too slow: {response.elapsed.total_seconds()}s")
        except Exception as e:
            print(f"Flask request failed: {str(e)}")
class JavaUser(IncreasingRateUser):
    host = "http://localhost:8091"
    
    @task
    def ping_java(self):
        try:
            with self.client.get("/actuator/prometheus", name="Java Ping", catch_response=True) as response:
                if response.status_code != 200:
                    response.failure(f"Java ping failed with status code: {response.status_code}")
                elif response.elapsed.total_seconds() > 1:  # If response takes more than 1 second
                    response.failure(f"Java ping too slow: {response.elapsed.total_seconds()}s")
        except Exception as e:
            print(f"Java request failed: {str(e)}")


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
    # Create an Environment instance
    env = Environment(user_classes=[GinUser])

    # Create a LocalRunner instance
    runner = LocalRunner(env)
    env.runner = runner

    # Start a greenlet that periodically outputs the current stats
    gevent.spawn(stats_printer(env.stats))

    runner.start(10, 1)

    # Allow test to run
    try:
        while runner.state != "stopped":
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nStopping load test...")
    finally:
        runner.quit()