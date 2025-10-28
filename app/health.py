import logging
from http.server import BaseHTTPRequestHandler, HTTPServer
from threading import Thread

logger = logging.getLogger(__name__)

class _Handler(BaseHTTPRequestHandler):
    def log_message(self, format, *args):
        return

    def do_GET(self):
        if self.path == "/healthz":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"ok")
        else:
            self.send_response(404)
            self.end_headers()


def start_health_server(host: str = "0.0.0.0", port: int = 8080):
    server = HTTPServer((host, port), _Handler)
    thread = Thread(target=server.serve_forever, daemon=True)
    thread.start()
    logger.info("health_server_started", extra={"port": port})
    return server


def stop_health_server(server) -> None:
    try:
        server.shutdown()
    except Exception:
        pass
