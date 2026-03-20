import json
import threading
import unittest
from http.server import BaseHTTPRequestHandler, HTTPServer

from task5_go_compute_service.python.client import fetch_prime_count


class _Handler(BaseHTTPRequestHandler):
    expected_limit = 0

    def do_POST(self):  # noqa: N802
        length = int(self.headers.get("Content-Length", "0"))
        raw = self.rfile.read(length).decode("utf-8")
        data = json.loads(raw)
        limit = int(data["limit"])

        if limit != self.expected_limit:
            self.send_response(400)
            self.end_headers()
            self.wfile.write(b"bad limit")
            return

        resp = {"limit": limit, "prime_count": 7}
        payload = json.dumps(resp).encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def log_message(self, format, *args):  # noqa: A002
        # Silence noisy test output.
        return


def _serve(server):
    server.serve_forever()


class TestPrimeClient(unittest.TestCase):
    def test_fetch_prime_count(self):
        host = "127.0.0.1"
        # Bind to ephemeral port.
        httpd = HTTPServer((host, 0), _Handler)
        port = httpd.server_address[1]

        _Handler.expected_limit = 30

        t = threading.Thread(target=_serve, args=(httpd,), daemon=True)
        t.start()
        def _cleanup():
            httpd.shutdown()
            httpd.server_close()

        self.addCleanup(_cleanup)

        got = fetch_prime_count(host, port, 30, timeout=2.0)
        self.assertEqual(got, 7)


if __name__ == "__main__":
    unittest.main()

