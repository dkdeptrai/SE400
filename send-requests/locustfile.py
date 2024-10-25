from locust import HttpUser, task, constant, between
import random

class PingUser(HttpUser):
    wait_time = constant(0)

    # List of hosts to switch between
    hosts = ["http://localhost:5500", "http://localhost:8090"]

    @task
    def ping(self):
        # Randomly select a host for the request
        host = random.choice(self.hosts)
        self.client.get(f"{host}/ping")

        # Wait for a short time after the request
        self.wait_task(0.01)
