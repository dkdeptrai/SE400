from flask import Flask, Response
from prometheus_client import generate_latest, Counter, Gauge, Histogram, CONTENT_TYPE_LATEST
import psutil
import time

app = Flask(__name__)

# Prometheus metrics
REQUEST_COUNTER = Counter('flask_app_requests_total', 'Total number of requests to the Flask app')
REQUEST_LATENCY = Histogram('flask_app_request_latency_seconds', 'Latency of requests to the Flask app')
CPU_USAGE_GAUGE = Gauge('flask_app_cpu_usage_percent', 'CPU usage percentage of the Flask app')
RAM_USAGE_GAUGE = Gauge('flask_app_ram_usage_bytes', 'RAM usage in bytes of the Flask app')
RAM_USAGE_PERCENT_GAUGE = Gauge('flask_app_ram_usage_percent', 'RAM usage percentage of the Flask app')

@app.route('/ping')
@REQUEST_LATENCY.time()  # Automatically measures latency of this route
def ping():
    REQUEST_COUNTER.inc()  # Count requests to /ping
    return {'message': 'pong'}

@app.route('/metrics')
def metrics():
    # Update CPU and RAM metrics
    CPU_USAGE_GAUGE.set(psutil.cpu_percent())  # CPU percentage
    RAM_USAGE_GAUGE.set(psutil.virtual_memory().used)  # RAM in bytes
    RAM_USAGE_PERCENT_GAUGE.set(psutil.virtual_memory().percent)
    # Return metrics for Prometheus scraping
    return Response(generate_latest(), mimetype=CONTENT_TYPE_LATEST)

if __name__ == '__main__':
    app.run(port=5500)
