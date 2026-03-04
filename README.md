# Songbeamer Helper

Ein Go-Tool zur Verwaltung und Synchronisierung von Songbeamer-Liedern mit ChurchTools.

## Anforderungen

- **Go**: Version 1.21 oder höher
- **Songbeamer**: Eine lokale Songbeamer-Installation
- **ChurchTools**: Ein ChurchTools-Account mit API-Zugriffsrechten
- **Konfigurationsdatei**: `songbeamer-helper.yaml` (siehe Konfiguration)

## Installation

### Abhängigkeiten installieren

```bash
go get
```

### Projekt bauen

```bash
./build.sh
```

### Tests ausführen

```bash
./test.sh
```

## Konfiguration

Die Konfiguration erfolgt über die Datei `songbeamer-helper.yaml`, die sich im Home-Verzeichnis oder im aktuellen Arbeitsverzeichnis befinden muss.

### Beispiel-Konfiguration

```yaml
songspath: /pfad/zu/Songbeamer/Songs
duplicates: /pfad/zu/Songbeamer/SongDuplicates


# Weitere ChurchTools-Konfiguration nach Bedarf
```

### Erforderliche Konfigurationsfelder

- **songspath**: Pfad zum Songbeamer Songs-Verzeichnis (lokal)
- **duplicates**: Pfad zum Songbeamer Duplicates-Verzeichnis (lokal)

## Geheimnisse (lokal installiert)

Die folgenden sensiblen Informationen sollten lokal in der `songbeamer-helper.yaml` gespeichert werden:

- **ChurchTools API-Token**: Authentifizierung für ChurchTools-API
- **AWS-Zugriffsschlüssel**: Für AWS S3-Integration
- **API-Endpoint**: URL der ChurchTools-Installation

Diese Datei **sollte nicht in der Versionskontrolle committet werden**. Nutze stattdessen eine `.gitignore`-Regel:

```
songbeamer-helper.yaml
```

## Verfügbare Befehle

Das Tool stellt verschiedene Befehle über die CLI zur Verfügung:

- **download**: Lädt Lieder von ChurchTools herunter
- **upload**: Lädt Lieder zu ChurchTools hoch
- **update**: Aktualisiert vorhandene Lieder
- **auto-upload**: Automatisches Upload-System
- **churchsongid**: Verwaltet ChurchTools-Song-IDs

Verwende `songbeamer-helper --help` für weitere Informationen.

## Entwicklung

### Projektstruktur

- `churchtools/`: API-Bindings für ChurchTools
- `cmd/`: Befehlsimplementierungen
- `songbeamer/`: Songbeamer-Dateiverwaltung
- `log/`: Logging-Modul
- `util/`: Utility-Funktionen

### Neue Version veröffentlichen

1. Version im Build-Befehl aktualisieren
2. Projekt bauen: `go build ...`
3. Self-Update-Manifest generieren: `go-selfupdate songbeamer-helper VERSION`
4. Zu S3 synchronisieren:

```bash
./build.sh
```

## Lizenz

Siehe LICENSE-Datei für Details.
