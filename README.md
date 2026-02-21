# Barcode-Logger

A zero-dependency, single-file Progressive Web App (PWA) for scanning, logging, and exporting barcodes and QR codes. Runs entirely in the browser — no server, no build step, no install required.

---

## Overview

**Barcode-Logger** turns any smartphone or desktop browser into a multi-mode barcode scanner with persistent local storage, CSV/JSON export, optional webhook sync, and full offline support. Common use cases include:

- **Receiving & dispatch** — scan inbound parcels with recipient and reference tracking  
- **Warehouse stock movements** — book items in by EAN with quantity and shelf location  
- **Maintenance & asset management** — log work orders, equipment tags, and technician assignments  
- **General inventory** — quick capture of any barcode or QR code with optional notes  

---

## Features at a Glance

- 📷 **Live camera scanning** via ZXing (15 configurable barcode formats)  
- 📋 **Four operational modes** with context-specific input fields and live stats  
- 📍 **Optional GPS tracking** — attach coordinates to each scan  
- 🔔 **Audio + haptic feedback** on successful scan  
- 🔦 **Torch / flashlight** control on supported devices  
- 📋 **Filterable scan list** — text search + mode filter  
- 📤 **CSV export** (15 columns), Web Share API, mailto fallback  
- 💾 **JSON backup & restore**  
- 🔗 **Webhook sync** — HTTP POST each completed scan to any endpoint  
- 🌐 **Internationalization** — German (default) and English; extensible  
- ♿ **Accessible** — ARIA roles, `aria-live` regions, keyboard tab navigation  
- 🎨 **Auto/Light/Dark theme** — respects `prefers-color-scheme`  
- 📦 **PWA** — installable, offline-capable via `manifest.webmanifest`  

---

## Operational Modes

| Mode | Badge colour | Auto-settings | Workflow CTA |
|------|-------------|---------------|-------------|
| 📋 General | Grey | — | 💾 Save |
| 📦 Delivery | Blue | Location ON, duplicates ON | ✅ Delivered |
| 🗃️ Stock | Green | — | 📥 Book in |
| 🔧 Maintenance | Amber | — | 🔧 Log |

### General
Captures any barcode with an optional free-text note. Suitable for ad-hoc scanning. Live stats show today's count and the running total.

### Delivery
Designed for parcel receiving. Automatically enables GPS location and allows duplicate scans (multiple parcels with the same label). Fields: quantity, recipient name, reference (order/delivery note), notes (e.g. safe-place instructions). Stats include how many scans have GPS coordinates attached.

### Stock
Tracks inventory movements by EAN or shelf code. Fields: quantity, shelf location, reference (delivery/order number), notes. Stats display bookings today, unique SKU count, and total item quantity.

### Maintenance
Records equipment work orders. Fields: work type (dropdown: Inspection / Repair / Replacement / Check / Cleaning), technician name, status (dropdown: open / in progress / completed), notes. Stats show today's count, open items, and completed items.

---

## Scan Panel

### Mode Selector
A button group (`#modeSeg`) at the top of the scan panel lets you switch modes. Switching a mode may auto-enable location tracking (Delivery) or allow duplicates (Delivery). Settings are persisted across sessions.

### Context Banner (`#modeBanner`)
Below the mode selector an `aria-live` region shows:
- The current mode name  
- Live stats pills (updated after every scan)  
- An italic prompt describing what to scan  

### Camera
The `<video>` element occupies the top of the panel. Scanning starts automatically on page load. If the camera is unavailable, the user is prompted to check HTTPS and browser permissions. Once running, the torch button appears if the device supports it.

### Settings Toggles
| Toggle | Default | Description |
|--------|---------|-------------|
| Track location | Off (On for Delivery) | Attaches GPS coordinates to each entry |
| Allow duplicates | Off (On for Delivery) | If off, a code already in the list is silently skipped |
| Audio feedback | On | Plays a short beep and triggers vibration on scan |

### Start / Stop / Clear
- **Start** — opens the camera and begins decoding  
- **Stop** — releases the camera track  
- **Clear list** — removes all entries from the in-memory list and localStorage (prompts confirmation)

---

## Detail Card (`#lastScanCard`)

After each successful scan a detail card slides in showing the scanned code and the mode-specific input fields. The user can fill in extra context and press the CTA button to save, or press **✕ Skip** to dismiss without extra fields (the raw scan entry is already saved at scan time).

---

## List View

Switch to the **List** tab to see all captured entries. Controls:

- **Text search** — filters by barcode value  
- **Mode dropdown** — shows only entries from the selected mode  
- Each row shows: code, mode badge, format, timestamp, optional GPS, and any extra fields as compact pills  

The `<ul>` has `aria-live="polite"` so screen readers announce new entries.

---

## Export

Found on the **Export** tab.

### CSV Export
All scans are serialised to a 15-column CSV:

| Column | Description |
|--------|-------------|
| `code` | Raw barcode/QR value |
| `format` | ZXing format name (e.g. `QR_CODE`, `EAN_13`) |
| `mode` | Operational mode key (`general`, `delivery`, `stock`, `maintenance`) |
| `qty` | Quantity (Stock / Delivery) |
| `shelf` | Shelf location (Stock) |
| `recipient` | Recipient name (Delivery) |
| `reference` | Order / delivery note reference |
| `work_type` | Maintenance work type |
| `technician` | Technician name (Maintenance) |
| `status` | Status value (Maintenance) |
| `notes` | Free-text notes |
| `timestamp` | Human-readable local time |
| `iso` | ISO 8601 UTC timestamp |
| `lat` | GPS latitude (Delivery, when enabled) |
| `lon` | GPS longitude |

**CSV exportieren / Export CSV** — triggers a file download.  
**Per Mail/Share senden / Send via Mail/Share** — uses the Web Share API (if available) to share the file; falls back to a `mailto:` link with the CSV body + a parallel download.

### JSON Backup / Restore
**Backup (JSON)** — downloads all entries as a pretty-printed JSON array.  
**Restore (JSON)** — opens a file picker and merges the loaded entries into the current list.

### Delete All
**Alles löschen / Delete all** — clears all entries after a confirmation prompt.

---

## Sync Backend (Webhook)

Configure a URL in the **Config** tab under *Sync Backend*. Every time a scan detail card is saved or dismissed, the full entry object is sent as an HTTP POST:

```json
{
  "id": "uuid-or-timestamp",
  "code": "4006381333931",
  "format": "EAN_13",
  "mode": "stock",
  "qty": 2,
  "shelf": "A1-03",
  "reference": "LN-4711",
  "notes": "",
  "ts": 1718000000000,
  "lat": 48.8566,
  "lon": 2.3522
}
```

The request uses `keepalive: true` so it can complete even if the user navigates away or closes the tab.

**Test** sends a synthetic `TEST-123` QR entry to verify connectivity. The status line below the input field shows the HTTP response code or a connection error message.

Leave the URL field empty and press **Save** to disable the webhook.

---

## Configuration

The **Config** tab lists all barcode formats supported by the installed ZXing build. Check the formats you want to detect and press **Apply**. If the camera is running when you apply, it is stopped and restarted automatically with the new hint set. Defaults include QR_CODE, EAN_13, EAN_8, CODE_128, CODE_39, ITF, UPC_A, UPC_E, DATA_MATRIX, AZTEC, PDF_417.

---

## Internationalization

The app ships with **German (de, default)** and **English (en)** translations, switchable via the DE / EN button group in the header. The selected language is persisted in `localStorage` under the key `barcode_logger_lang`.

### How translations work

All translation strings live in a single `TRANSLATIONS` object at the top of the `<script>`:

```js
const TRANSLATIONS = {
  de: { tab_scan: 'Scan', ... },
  en: { tab_scan: 'Scan', ... }
};
let currentLang = localStorage.getItem(LS_KEY_LANG) || 'de';
function t(key){ return (TRANSLATIONS[currentLang]||TRANSLATIONS.de)[key] || TRANSLATIONS.de[key] || key; }
```

`applyTranslations()` walks every `[data-i18n]`, `[data-i18n-placeholder]`, `[data-i18n-aria-label]`, `[data-i18n-title]`, and `[data-i18n-mode-label]` element and sets the appropriate attribute.

### Adding a new language (e.g. French)

1. Add a `fr: { ... }` key to `TRANSLATIONS` with all the same keys as `de` and `en`.  
2. Add a `<button id="lang-fr" aria-pressed="false">FR</button>` inside `#langSwitcher`.  
3. In `initLang()`, add a click listener for `lang-fr` that sets `currentLang = 'fr'` and calls `applyTranslations()`.

---

## Accessibility

- **Keyboard navigation** — arrow keys (← →) cycle through the four header tabs when focus is on the tab bar.  
- **Tab panels** have `tabindex="0"` so they can receive keyboard focus.  
- **`aria-live="polite"`** on `#modeBanner` (region), `#scanList`, `#supportNote`, and `#webhookStatus` announces updates to screen readers without interrupting.  
- **`role="status"`** on `#count` communicates the scan count.  
- **`role="region"` + `aria-labelledby`** on the detail card `#lastScanCard`.  
- **`aria-pressed`** on mode buttons, theme buttons, and language buttons.  
- All `aria-label` attributes that contain visible text are backed by `data-i18n-aria-label` so they update when the language changes.  

---

## Extensibility

### Adding a New Mode

1. **MODES** — add a new entry in `const MODES` with `get label()`, `get prompt()`, `get saveCta()`, `autoLocation`, `autoAllowDupes`, and `stats(s)`.  
2. **getModeFields()** — add a corresponding key with the field definitions.  
3. **HTML buttons** — add a `<button data-mode="newmode" aria-pressed="false" data-i18n-mode-label="newmode">` in `#modeSeg`.  
4. **Translations** — add `mode_newmode`, `prompt_newmode`, `cta_newmode` keys (and any field/placeholder keys) to both `de` and `en` in `TRANSLATIONS`.  
5. **CSS badge** — add `.mode-badge-newmode { background: <colour>; }` and `.mode-seg button[data-mode="newmode"][aria-pressed="true"]` rules.  
6. **Filter option** — add `<option value="newmode" data-i18n-mode-label="newmode">` in `#filterMode`.

---

## PWA

The app references `manifest.webmanifest` and two icon sizes (`icon-192.png`, `icon-512.png`). When these files are present on the same origin, browsers offer an "Add to Home Screen" / "Install" prompt. Once installed, the app runs from the browser cache and works fully offline (no network requests except for the optional webhook and CDN ZXing load).

---

## Theme

Three modes toggled via the header button group:

| Button | Behaviour |
|--------|-----------|
| Auto | Follows `prefers-color-scheme` (light/dark OS setting) |
| Light | Forces light colour palette |
| Dark | Forces dark colour palette |

Theme preference is persisted in `localStorage` under `barcode_logger_theme`.

---

## Technical Notes

| Topic | Detail |
|-------|--------|
| **No build step** | Single HTML file; open directly in a browser or serve from any static host |
| **ZXing** | Loaded from `unpkg.com/@zxing/library@latest/umd/index.min.js` |
| **localStorage keys** | `barcode_logger_scans_v3`, `barcode_logger_theme`, `barcode_logger_formats`, `barcode_logger_mode`, `barcode_logger_webhook`, `barcode_logger_lang` |
| **Encoding** | UTF-8 throughout; all string replacements must use UTF-8-aware tools |
| **No external CSS/JS** | All styles and logic are inline in the single file |
| **Duplicate detection** | Handled in app logic (not via ZXing hints); quick-repeat suppression of 800 ms prevents double-scans |
| **Geo-watching** | `navigator.geolocation.watchPosition` — started when location toggle is enabled or when switching to Delivery mode |
