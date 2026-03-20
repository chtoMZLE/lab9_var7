import http.client
import json
from typing import Any, Dict


def fetch_prime_count(host: str, port: int, limit: int, timeout: float = 5.0) -> int:
    """
    Calls Go primeservice:
    POST /primes
    Request:  {"limit": <int>}
    Response: {"limit": <int>, "prime_count": <int>}
    """
    conn = http.client.HTTPConnection(host, port, timeout=timeout)
    payload: Dict[str, Any] = {"limit": int(limit)}
    body = json.dumps(payload).encode("utf-8")

    headers = {"Content-Type": "application/json; charset=utf-8"}
    conn.request("POST", "/primes", body=body, headers=headers)
    resp = conn.getresponse()
    raw = resp.read().decode("utf-8")

    if resp.status != 200:
        raise RuntimeError(f"unexpected status {resp.status}: {raw}")

    data = json.loads(raw)
    return int(data["prime_count"])


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=9001)
    parser.add_argument("--limit", type=int, required=True)
    args = parser.parse_args()

    print(fetch_prime_count(args.host, args.port, args.limit))

