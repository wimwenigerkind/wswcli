# Änderungsprotokoll

Alle bemerkenswerten Änderungen an diesem Projekt werden in dieser Datei dokumentiert.

Das Format basiert auf [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
und dieses Projekt folgt [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unveröffentlicht]

### Hinzugefügt
- Nightly-Release-Pipeline mit automatischen täglichen Builds
- Nightly-Docker-Images mit `nightly`, `nightly-amd64`, `nightly-arm64` Tags
- Manuelle Auslösung für Nightly-Releases über GitHub Actions
- Intelligente Nightly-Release-Logik (erstellt nur Release bei neuen Commits)
- Separate GoReleaser-Konfiguration für Nightly-Builds

## [2.1.0] - 2025-07-10

### Hinzugefügt
- Docker-Unterstützung mit Multi-Architektur-Images (AMD64, ARM64)
- Multi-Platform Docker-Manifeste für automatische Plattformauswahl
- Verbessertes Dockerfile mit Sicherheitsverbesserungen (Non-Root-User, CA-Zertifikate)
- GoReleaser-Konfiguration auf Version 2 aktualisiert

### Geändert
- Verbesserte Dokumentationsstruktur und -inhalte
- GitHub Actions Workflow auf GoReleaser v2 aktualisiert (goreleaser-action@v6)

### Behoben
- GitHub Actions Release-Workflow Kompatibilität mit GoReleaser Konfigurationsversion 2
- Fehlende `packages: write` Berechtigung für Docker-Image-Publishing zur GitHub Container Registry hinzugefügt
- Docker-Authentifizierungsproblem, das das Veröffentlichen von Images zur GitHub Container Registry verhinderte

## [2.0.0] - 2025-07-05

### Hinzugefügt
- Homebrew-Installation über benutzerdefinierten Tap
- Installationsanweisungen für mehrere Paketmanager

### Geändert
- **BREAKING**: `generateUnifiedDiff` refaktoriert, um System-Git-Diff für verbesserte Performance zu verwenden
- Verbesserte Fehlerbehandlung und Performance-Optimierungen

### Behoben
- Performance-Verbesserungen bei der Diff-Generierung

## [1.0.0] - 2025-07-04

### Hinzugefügt
- Erste stabile Version von wswcli
- PatchVendor-Befehl zur Generierung von Unified-Diff-Patches
- Verzeichnisverarbeitungsfunktionen
- Interaktiver Modus mit geführtem Workflow
- Intelligente Validierung mit umfassenden Fehlermeldungen
- Vendor-Pfad-Behandlung und -Normalisierung
- Unterstützung für Shopware-Vendor-Modifikationen

### Behoben
- Grep-Befehl-Kompatibilitätsproblem mit nicht erkannter Option behoben
- Verbesserte PHP-Klassendefinitionsformatierung in Testdateien

### Geändert
- Projekt von ursprünglichem Namen zu wswcli umbenannt
- Alle Konfigurationen und Dokumentationen aktualisiert, um den neuen Projektnamen zu reflektieren

---

## Release-Hinweise

### Docker-Verwendung
Ab der unveröffentlichten Version können Sie wswcli mit Docker ausführen:

```bash
# Neueste Version (Multi-Platform)
docker run --rm ghcr.io/wimwenigerkind/wswcli:latest --version

# Spezifische Architektur
docker run --rm ghcr.io/wimwenigerkind/wswcli:latest-arm64 --help
```

### Installationsmethoden
- **Go**: `go install github.com/wimwenigerkind/wswcli@latest`
- **Homebrew**: `brew install wimwenigerkind/tap/wswcli`
- **Docker**: `docker run ghcr.io/wimwenigerkind/wswcli:latest`
- **Binary**: Download von der [Releases-Seite](https://github.com/wimwenigerkind/wswcli/releases)

[Unveröffentlicht]: https://github.com/wimwenigerkind/wswcli/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/wimwenigerkind/wswcli/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/wimwenigerkind/wswcli/releases/tag/v1.0.0