import socket
import sys


def main() -> None:
    host = sys.argv[1] if len(sys.argv) >= 3 else "127.0.0.1"
    port = int(sys.argv[2]) if len(sys.argv) >= 3 else 9000
    msg = sys.argv[3] if len(sys.argv) >= 4 else "hello"

    with socket.create_connection((host, port), timeout=5.0) as sock:
        sock.sendall((msg + "\n").encode("utf-8"))
        data = b""
        while b"\n" not in data:
            chunk = sock.recv(4096)
            if not chunk:
                break
            data += chunk

    # Keep output single-line.
    sys.stdout.write(data.decode("utf-8").splitlines()[0] + "\n")


if __name__ == "__main__":
    main()

