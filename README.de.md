# Barcode-Logger

Eine abhängigkeitsfreie, einzelne HTML-Datei als Progressive Web App (PWA) zum Scannen, Protokollieren und Exportieren von Barcodes und QR-Codes. Läuft vollständig im Browser – kein Server, kein Build-Schritt, keine Installation erforderlich.

---

## Übersicht

Der **Barcode-Logger** verwandelt jedes Smartphone oder jeden Desktop-Browser in einen Mehrfachmode-Barcode-Scanner mit persistentem lokalem Speicher, CSV/JSON-Export, optionaler Webhook-Synchronisation und vollständiger Offline-Unterstützung. Typische Anwendungsfälle:

- **Wareneingang & Versand** – Pakete mit Empfänger- und Referenzverfolgung scannen
- **Lagerbewegungen** – Artikel per EAN mit Menge und Lagerplatz einbuchen
- **Wartung & Asset-Management** – Arbeitsaufträge, Geräte-Tags und Technikerzuweisungen protokollieren
- **Allgemeine Inventur** – schnelle Erfassung beliebiger Barcodes oder QR-Codes mit optionalen Notizen

---

## Funktionen auf einen Blick

- 📷 **Live-Kamera-Scan** per ZXing (15 konfigurierbare Barcode-Formate)
- 📋 **Vier Betriebsmodi** mit kontextspezifischen Eingabefeldern und Live-Statistiken
- 📍 **Optionale GPS-Erfassung** – Koordinaten jedem Scan anhängen
- 🔔 **Audio- und Haptik-Feedback** bei erfolgreichem Scan
- 🔦 **Taschenlampensteuerung** auf unterstützten Geräten
- 📋 **Filterable Scan-Liste** – Textsuche und Modus-Filter
- 📤 **CSV-Export** (15 Spalten), Web Share API, Mailto-Fallback
- 💾 **JSON-Backup & Wiederherstellung**
- 🔗 **Webhook-Sync** – jeden abgeschlossenen Scan per HTTP POST an einen beliebigen Endpunkt senden
- 🌐 **Internationalisierung** – Deutsch (Standard) und Englisch; erweiterbar
- ♿ **Barrierefrei** – ARIA-Rollen, `aria-live`-Regionen, Tastatur-Tab-Navigation
- 🎨 **Auto/Hell/Dunkel-Theme** – reagiert auf `prefers-color-scheme`
- 📦 **PWA** – installierbar, offline-fähig via `manifest.webmanifest`

---

## Betriebsmodi

| Modus | Badge-Farbe | Auto-Einstellungen | Workflow-CTA |
|-------|------------|-------------------|-------------|
| 📋 Allgemein | Grau | – | 💾 Speichern |
| 📦 Lieferung | Blau | Standort AN, Dubletten AN | ✅ Zugestellt |
| 🗃️ Lager | Grün | – | 📥 Einbuchen |
| 🔧 Wartung | Amber | – | 🔧 Protokollieren |

### Allgemein
Erfasst beliebige Barcodes mit optionaler Freitextnotiz. Geeignet für Ad-hoc-Scans. Live-Statistiken zeigen die heutige Anzahl und die laufende Gesamtzahl.

### Lieferung
Konzipiert für die Paketzustellung. Aktiviert automatisch GPS-Standort und erlaubt Duplikat-Scans (mehrere Pakete mit dem gleichen Label). Felder: Menge, Empfängername, Referenz (Auftrag/Lieferschein), Notizen (z. B. Abstellgenehmigung). Statistiken zeigen, wie viele Scans GPS-Koordinaten enthalten.

### Lager
Verfolgt Lagerbewegungen per EAN oder Regalcode. Felder: Menge, Lagerplatz, Referenz (Lieferung/Bestellnummer), Notizen. Statistiken zeigen Buchungen heute, eindeutige SKU-Anzahl und Gesamtmenge.

### Wartung
Erfasst Geräte-Wartungsaufträge. Felder: Arbeitstyp (Dropdown: Inspektion / Reparatur / Austausch / Prüfung / Reinigung), Techniker-Name, Status (Dropdown: offen / in Bearbeitung / abgeschlossen), Notizen. Statistiken zeigen heutige Anzahl, offene und abgeschlossene Einträge.

---

## Scan-Panel

### Modus-Auswahl
Eine Schaltflächengruppe (`#modeSeg`) am oberen Rand des Scan-Panels ermöglicht den Wechsel zwischen Modi. Bei Moduswechsel können Standortverfolgung (Lieferung) oder Dubletten-Erlaubnis (Lieferung) automatisch aktiviert werden. Einstellungen werden sitzungsübergreifend gespeichert.

### Kontext-Banner (`#modeBanner`)
Unterhalb der Modus-Auswahl zeigt eine `aria-live`-Region:
- Den aktuellen Modusnamen
- Live-Statistik-Pills (nach jedem Scan aktualisiert)
- Eine kursive Aufforderung, was gescannt werden soll

### Kamera
Das `<video>`-Element belegt den oberen Bereich des Panels. Das Scannen startet automatisch beim Laden der Seite. Falls die Kamera nicht verfügbar ist, wird der Benutzer aufgefordert, HTTPS und Browser-Berechtigungen zu prüfen. Nach dem Start erscheint die Taschenlampen-Schaltfläche, wenn das Gerät diese unterstützt.

### Einstellungsschalter
| Schalter | Standard | Beschreibung |
|---------|---------|-------------|
| Standort mitschreiben | Aus (bei Lieferung: An) | Hängt GPS-Koordinaten an jeden Eintrag an |
| Dubletten erlauben | Aus (bei Lieferung: An) | Wenn deaktiviert, wird ein bereits in der Liste vorhandener Code lautlos übersprungen |
| Akustisches Feedback | An | Gibt einen kurzen Piepton aus und löst bei Scan eine Vibration aus |

### Start / Stopp / Liste leeren
- **Start** – öffnet die Kamera und beginnt mit der Dekodierung
- **Stopp** – gibt den Kamera-Track frei
- **Liste leeren** – entfernt alle Einträge aus der In-Memory-Liste und localStorage (Bestätigung erforderlich)

---

## Detail-Karte (`#lastScanCard`)

Nach jedem erfolgreichen Scan erscheint eine Detail-Karte mit dem gescannten Code und den moduspezifischen Eingabefeldern. Der Benutzer kann zusätzlichen Kontext eingeben und die CTA-Schaltfläche drücken, um zu speichern, oder **✕ Überspringen** drücken, um ohne Zusatzfelder zu verwerfen (der rohe Scan-Eintrag wird bereits beim Scannen gespeichert).

---

## Listenansicht

Wechseln Sie zur Registerkarte **Liste**, um alle erfassten Einträge anzuzeigen. Steuerelemente:

- **Textsuche** – filtert nach Barcode-Wert
- **Modus-Dropdown** – zeigt nur Einträge aus dem ausgewählten Modus an
- Jede Zeile zeigt: Code, Modus-Badge, Format, Zeitstempel, optionaler GPS und alle Zusatzfelder als kompakte Pills

Das `<ul>` hat `aria-live="polite"`, damit Screenreader neue Einträge ankündigen.

---

## Export

Zu finden auf der Registerkarte **Export**.

### CSV-Export
Alle Scans werden in ein 15-spaltiges CSV serialisiert:

| Spalte | Beschreibung |
|--------|-------------|
| `code` | Roher Barcode-/QR-Wert |
| `format` | ZXing-Format-Name (z. B. `QR_CODE`, `EAN_13`) |
| `mode` | Betriebsmodus-Schlüssel (`general`, `delivery`, `stock`, `maintenance`) |
| `qty` | Menge (Lager / Lieferung) |
| `shelf` | Lagerplatz (Lager) |
| `recipient` | Empfängername (Lieferung) |
| `reference` | Auftrags- / Lieferschein-Referenz |
| `work_type` | Wartungsarbeitstyp |
| `technician` | Technikernamme (Wartung) |
| `status` | Statuswert (Wartung) |
| `notes` | Freitextnotizen |
| `timestamp` | Für Menschen lesbare Ortszeit |
| `iso` | ISO 8601 UTC-Zeitstempel |
| `lat` | GPS-Breitengrad (Lieferung, wenn aktiviert) |
| `lon` | GPS-Längengrad |

**CSV exportieren** – löst einen Datei-Download aus.  
**Per Mail/Share senden** – verwendet die Web Share API (falls verfügbar), um die Datei zu teilen; fällt auf einen `mailto:`-Link mit dem CSV-Inhalt + parallelen Download zurück.

### JSON-Backup / Wiederherstellen
**Backup (JSON)** – lädt alle Einträge als formatiertes JSON-Array herunter.  
**Wiederherstellen (JSON)** – öffnet eine Dateiauswahl und fügt die geladenen Einträge zur aktuellen Liste hinzu.

### Alles löschen
**Alles löschen** – löscht alle Einträge nach einer Bestätigungsaufforderung.

---

## Sync-Backend (Webhook)

Konfigurieren Sie eine URL auf der Registerkarte **Konfig** unter *Sync-Backend*. Jedes Mal, wenn eine Scan-Detail-Karte gespeichert oder verworfen wird, wird das vollständige Eintragsobjekt als HTTP POST gesendet:

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

Die Anfrage verwendet `keepalive: true`, damit sie abgeschlossen werden kann, auch wenn der Benutzer navigiert oder den Tab schließt.

**Testen** sendet einen synthetischen `TEST-123` QR-Eintrag, um die Konnektivität zu überprüfen. Die Statuszeile unter dem Eingabefeld zeigt den HTTP-Antwortcode oder eine Verbindungsfehlermeldung.

Lassen Sie das URL-Feld leer und drücken Sie **Speichern**, um den Webhook zu deaktivieren.

### Einsatzbereite Backend-Beispiele

Der Ordner `backend/` enthält minimale, abhängigkeitsfreie Server-Implementierungen:

| Datei | Sprache | Startbefehl |
|-------|---------|-------------|
| `backend/sync.php` | PHP 7.4+ | Auf einem PHP-fähigen Webserver bereitstellen (Apache, Nginx + php-fpm). Endpunkt: `https://ihr-server/backend/sync.php` |
| `backend/sync.py` | Python 3.7+ (nur stdlib) | `python3 backend/sync.py --port 8080` → Endpunkt `http://localhost:8080/scan` |
| `backend/sync.go` | Go 1.18+ | `go run backend/sync.go -addr :8080` → Endpunkt `http://localhost:8080/scan` |

Alle drei Backends:
- Akzeptieren HTTP POST mit `Content-Type: application/json`
- Prüfen, dass das Pflichtfeld `code` vorhanden ist
- Hängen jeden Eintrag als newline-delimitiertes JSON an `scans.ndjson` an
- Antworten mit `{"ok": true, "code": "...", "received_at": "..."}` bei Erfolg
- Setzen permissive CORS-Header (für Produktion anpassen)
- Behandeln OPTIONS-Preflight-Anfragen

Für Produktionsumgebungen: Authentifizierung hinzufügen (z. B. ein gemeinsames Geheimnis im `Authorization`-Header) und hinter einem TLS-terminierenden Reverse-Proxy (nginx / caddy) betreiben.

---

## Konfiguration

Die Registerkarte **Konfig** listet alle Barcode-Formate auf, die vom installierten ZXing-Build unterstützt werden. Aktivieren Sie die Formate, die erkannt werden sollen, und drücken Sie **Übernehmen**. Wenn die Kamera beim Übernehmen läuft, wird sie automatisch gestoppt und mit dem neuen Hinweis-Set neu gestartet. Standardformate: QR_CODE, EAN_13, EAN_8, CODE_128, CODE_39, ITF, UPC_A, UPC_E, DATA_MATRIX, AZTEC, PDF_417.

---

## Internationalisierung

Die App enthält **Deutsch (de, Standard)** und **Englisch (en)** als Übersetzungen, umschaltbar über die DE / EN-Schaltflächengruppe im Header. Die ausgewählte Sprache wird in `localStorage` unter dem Schlüssel `barcode_logger_lang` gespeichert.

### Funktionsweise der Übersetzungen

Alle Übersetzungsstrings befinden sich in einem einzigen `TRANSLATIONS`-Objekt am Anfang des `<script>`:

```js
const TRANSLATIONS = {
  de: { tab_scan: 'Scan', ... },
  en: { tab_scan: 'Scan', ... }
};
let currentLang = localStorage.getItem(LS_KEY_LANG) || 'de';
function t(key){ return (TRANSLATIONS[currentLang]||TRANSLATIONS.de)[key] || TRANSLATIONS.de[key] || key; }
```

`applyTranslations()` durchläuft alle `[data-i18n]`-, `[data-i18n-placeholder]`-, `[data-i18n-aria-label]`-, `[data-i18n-title]`- und `[data-i18n-mode-label]`-Elemente und setzt das entsprechende Attribut.

### Eine neue Sprache hinzufügen (z. B. Französisch)

1. Fügen Sie einen `fr: { ... }`-Schlüssel zu `TRANSLATIONS` mit denselben Schlüsseln wie `de` und `en` hinzu.
2. Fügen Sie `<button id="lang-fr" aria-pressed="false">FR</button>` in `#langSwitcher` ein.
3. Fügen Sie in `initLang()` einen Click-Listener für `lang-fr` hinzu, der `currentLang = 'fr'` setzt und `applyTranslations()` aufruft.

---

## Barrierefreiheit

- **Tastaturnavigation** – Pfeiltasten (← →) wechseln zwischen den vier Header-Tabs, wenn der Fokus auf der Tab-Leiste liegt.
- **Tab-Panels** haben `tabindex="0"`, damit sie Tastaturfokus erhalten können.
- **`aria-live="polite"`** auf `#modeBanner` (Region), `#scanList`, `#supportNote` und `#webhookStatus` kündigt Aktualisierungen für Screenreader an, ohne zu unterbrechen.
- **`role="status"`** auf `#count` kommuniziert die Scan-Anzahl.
- **`role="region"` + `aria-labelledby`** auf der Detail-Karte `#lastScanCard`.
- **`aria-pressed`** auf Modus-Schaltflächen, Theme-Schaltflächen und Sprach-Schaltflächen.
- Alle `aria-label`-Attribute mit sichtbarem Text werden durch `data-i18n-aria-label` gesichert, damit sie sich bei Sprachwechsel aktualisieren.

---

## Erweiterbarkeit

### Einen neuen Modus hinzufügen

1. **MODES** – neuen Eintrag in `const MODES` mit `get label()`, `get prompt()`, `get saveCta()`, `autoLocation`, `autoAllowDupes` und `stats(s)` hinzufügen.
2. **getModeFields()** – entsprechenden Schlüssel mit den Felddefinitionen hinzufügen.
3. **HTML-Schaltflächen** – `<button data-mode="neuerModus" aria-pressed="false" data-i18n-mode-label="neuerModus">` in `#modeSeg` einfügen.
4. **Übersetzungen** – Schlüssel `mode_neuerModus`, `prompt_neuerModus`, `cta_neuerModus` (und alle Feld-/Platzhalter-Schlüssel) zu `de` und `en` in `TRANSLATIONS` hinzufügen.
5. **CSS-Badge** – `.mode-badge-neuerModus { background: <Farbe>; }` und `.mode-seg button[data-mode="neuerModus"][aria-pressed="true"]` Regeln hinzufügen.
6. **Filter-Option** – `<option value="neuerModus" data-i18n-mode-label="neuerModus">` in `#filterMode` einfügen.

---

## PWA

Die App referenziert `manifest.webmanifest` und zwei Icon-Größen (`icon-192.png`, `icon-512.png`). Wenn diese Dateien auf demselben Origin vorhanden sind, bieten Browser eine "Zum Startbildschirm hinzufügen" / "Installieren"-Aufforderung an. Nach der Installation läuft die App aus dem Browser-Cache und funktioniert vollständig offline (keine Netzwerkanfragen außer für den optionalen Webhook und das CDN-ZXing-Laden).

---

## Theme

Drei Modi, umschaltbar über die Header-Schaltflächengruppe:

| Schaltfläche | Verhalten |
|-------------|----------|
| Auto | Folgt `prefers-color-scheme` (Hell/Dunkel-OS-Einstellung) |
| Hell | Erzwingt helle Farbpalette |
| Dunkel | Erzwingt dunkle Farbpalette |

Die Theme-Präferenz wird in `localStorage` unter `barcode_logger_theme` gespeichert.

---

## Technische Hinweise

| Thema | Detail |
|-------|--------|
| **Kein Build-Schritt** | Einzelne HTML-Datei; direkt im Browser öffnen oder von einem beliebigen statischen Host bereitstellen |
| **ZXing** | Geladen von `unpkg.com/@zxing/library@latest/umd/index.min.js` |
| **localStorage-Schlüssel** | `barcode_logger_scans_v3`, `barcode_logger_theme`, `barcode_logger_formats`, `barcode_logger_mode`, `barcode_logger_webhook`, `barcode_logger_lang` |
| **Kodierung** | Durchgehend UTF-8; alle String-Ersetzungen müssen UTF-8-bewusste Tools verwenden |
| **Kein externes CSS/JS** | Alle Styles und Logik sind inline in der einzelnen Datei |
| **Duplikaterkennung** | In der App-Logik (nicht über ZXing-Hints); schnelle Wiederholungsunterdrückung von 800 ms verhindert Doppelscans |
| **Geo-Watching** | `navigator.geolocation.watchPosition` – gestartet, wenn der Standort-Schalter aktiviert wird oder beim Wechsel in den Lieferungsmodus |
