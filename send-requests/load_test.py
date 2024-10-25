from locust import HttpUser, task, between, events
import random
import time
import logging

# Define the two target URLs
SERVERS = [
    "http://localhost:5500",  # Flask server
    "http://localhost:8090"    # Gin server
]

class LoadTestUser(HttpUser):
    wait_time = between(0.5, 2)  # Random wait time between requests (in seconds)

    @task
    def ping_servers(self):
        # Randomly select one of the servers to send a request to
        target_server = random.choice(SERVERS)
        url = f"{target_server}/ping"  # Adjust the endpoint as necessary
        
        with self.client.get(url, catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Request failed with status code: {response.status_code}")
            elif response.elapsed.total_seconds() > 1:  # Check if response is slow
                response.failure(f"Response too slow: {response.elapsed.total_seconds()} seconds")

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    logging.info("Load test starting...")

@events.quitting.add_listener
def on_quitting(environment, **kwargs):
    logging.info("Load test finished.")

if __name__ == "__main__":
    import os
    from locust import run_single_user
    
    # Optionally set environment variables for custom configurations
    os.environ["LOCUST_TARGET_HOST"] = "http://localhost"  # Base URL, not used in this script directly

    run_single_user(LoadTestUser)
