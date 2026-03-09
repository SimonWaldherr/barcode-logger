#!/usr/bin/env python3
"""
Barcode-Logger – Minimalistic Sync Backend (Python, stdlib only)

Receives scan entries as HTTP POST (application/json) from the Barcode-Logger
PWA webhook and appends them to a newline-delimited JSON log file.

Usage
-----
    python3 sync.py [--port 8080] [--host 0.0.0.0] [--log scans.ndjson]

Then point the Barcode-Logger Webhook URL to:
    http://<your-host>:8080/scan

The server responds with JSON {"ok": true, "code": "<value>"} on success.

Security note
-------------
This is a minimal development/demo server.  For production use, run behind a
reverse proxy (nginx / caddy) that handles TLS, and add authentication (e.g. a
shared secret checked in the Authorization header).
"""

import argparse
import datetime
import fcntl
import http.server
import json
import os
import sys
from typing import Any


# ---------- Configuration defaults -------------------------------------------
DEFAULT_HOST = "0.0.0.0"
DEFAULT_PORT = 8080
DEFAULT_LOG  = os.path.join(os.path.dirname(__file__), "scans.ndjson")


# ---------- Request handler --------------------------------------------------
class ScanHandler(http.server.BaseHTTPRequestHandler):
    log_path: str = DEFAULT_LOG  # set by main() before serving

    # Suppress default access log to stderr; replace with structured output
    def log_message(self, fmt: str, *args: Any) -> None:  # type: ignore[override]
        ts = datetime.datetime.now().isoformat(timespec="seconds")
        print(f"[{ts}] {self.address_string()} {fmt % args}", flush=True)

    def _send_json(self, status: int, body: dict) -> None:
        payload = json.dumps(body, ensure_ascii=False).encode()
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(payload)))
        self._cors_headers()
        self.end_headers()
        self.wfile.write(payload)

    def _cors_headers(self) -> None:
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization")

    # ---- OPTIONS (pre-flight) -----------------------------------------------
    def do_OPTIONS(self) -> None:
        self.send_response(204)
        self._cors_headers()
        self.end_headers()

    # ---- POST /scan ---------------------------------------------------------
    def do_POST(self) -> None:
        if self.path.rstrip("/") not in ("/scan", ""):
            self._send_json(404, {"error": "Not found"})
            return

        length = int(self.headers.get("Content-Length", 0))
        if length == 0:
            self._send_json(400, {"error": "Empty body"})
            return

        raw = self.rfile.read(length)
        try:
            entry: dict = json.loads(raw)
        except json.JSONDecodeError:
            self._send_json(422, {"error": "Invalid JSON"})
            return

        if not isinstance(entry, dict) or "code" not in entry:
            self._send_json(422, {"error": 'Missing required field "code"'})
            return

        entry["received_at"] = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")

        # Append atomically using a file lock
        try:
            with open(self.log_path, "a", encoding="utf-8") as fh:
                fcntl.flock(fh, fcntl.LOCK_EX)
                fh.write(json.dumps(entry, ensure_ascii=False) + "\n")
                fcntl.flock(fh, fcntl.LOCK_UN)
        except OSError as exc:
            self._send_json(500, {"error": f"Cannot write log: {exc}"})
            return

        self._send_json(200, {"ok": True, "code": entry["code"], "received_at": entry["received_at"]})


# ---------- Entry point ------------------------------------------------------
def main() -> None:
    parser = argparse.ArgumentParser(description="Barcode-Logger sync backend")
    parser.add_argument("--host", default=DEFAULT_HOST)
    parser.add_argument("--port", type=int, default=DEFAULT_PORT)
    parser.add_argument("--log",  default=DEFAULT_LOG, metavar="FILE",
                        help="NDJSON log file path (default: scans.ndjson)")
    args = parser.parse_args()

    ScanHandler.log_path = args.log

    server = http.server.HTTPServer((args.host, args.port), ScanHandler)
    print(f"Barcode-Logger sync backend listening on http://{args.host}:{args.port}/scan", flush=True)
    print(f"Logging to: {args.log}", flush=True)
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nStopped.", flush=True)
        sys.exit(0)


if __name__ == "__main__":
    main()
