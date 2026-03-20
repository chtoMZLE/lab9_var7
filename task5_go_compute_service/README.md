# Task 5: Go heavy-compute microservice invoked from Python

HTTP service (Go) exposes:

- `POST /primes`
  - Request JSON: `{"limit": <int>}`
  - Response JSON: `{"limit": <int>, "prime_count": <int>}`

Python client:

- `task5_go_compute_service/python/client.py`

