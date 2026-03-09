<?php
/**
 * Barcode-Logger – Minimalistic Sync Backend (PHP)
 *
 * Receives scan entries as HTTP POST (application/json) from the Barcode-Logger
 * PWA webhook and appends them to a newline-delimited JSON log file.
 *
 * Usage
 * -----
 * 1. Copy this file to a PHP-enabled web server (e.g. Apache, Nginx + php-fpm).
 * 2. Make sure the directory is writable by the web-server process:
 *      chmod 755 /path/to/backend
 * 3. Point the Barcode-Logger Webhook URL to:
 *      https://your-server.example.com/backend/sync.php
 *
 * Storage
 * -------
 * Scans are appended to `scans.ndjson` (newline-delimited JSON) in the same
 * directory.  Each line is one complete JSON scan object.
 *
 * Security note
 * -------------
 * This is a minimal example.  For production use, add authentication
 * (e.g. a shared secret in the `Authorization` header), rate limiting,
 * and input size limits.
 */

declare(strict_types=1);

/* ---------- CORS (adjust origins for production) ---------- */
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(204);
    exit;
}

/* ---------- Only accept POST ---------- */
if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    http_response_code(405);
    header('Content-Type: application/json');
    echo json_encode(['error' => 'Method Not Allowed']);
    exit;
}

/* ---------- Read + validate body ---------- */
$raw = file_get_contents('php://input');
if ($raw === false || $raw === '') {
    http_response_code(400);
    header('Content-Type: application/json');
    echo json_encode(['error' => 'Empty body']);
    exit;
}

$entry = json_decode($raw, true);
if (!is_array($entry) || !isset($entry['code'])) {
    http_response_code(422);
    header('Content-Type: application/json');
    echo json_encode(['error' => 'Invalid JSON or missing "code" field']);
    exit;
}

/* ---------- Append to NDJSON log ---------- */
$logFile = __DIR__ . '/scans.ndjson';

// Ensure only one request writes at a time
$fp = fopen($logFile, 'a');
if ($fp === false) {
    http_response_code(500);
    header('Content-Type: application/json');
    echo json_encode(['error' => 'Cannot open log file']);
    exit;
}

if (flock($fp, LOCK_EX)) {
    $entry['received_at'] = date('c'); // server-side ISO timestamp
    fwrite($fp, json_encode($entry, JSON_UNESCAPED_UNICODE) . "\n");
    flock($fp, LOCK_UN);
}
fclose($fp);

/* ---------- Success response ---------- */
http_response_code(200);
header('Content-Type: application/json');
echo json_encode([
    'ok'          => true,
    'code'        => $entry['code'],
    'received_at' => $entry['received_at'],
]);
