# Änderungsprotokoll

Alle bemerkenswerten Änderungen an diesem Projekt werden in dieser Datei dokumentiert.

Das Format basiert auf [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
und dieses Projekt folgt [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unveröffentlicht]

### Entfernt
- Emojis aus der Ausgabe des Befehls `bs-4-to-5` entfernt für übersichtlicheren Text
- Twig-Kommentare in HTML-Dateien und umgekehrt für den `bs-4-to-5`-Befehl handhaben

## [v2.4.1] - 2025-07-21

### Entfernt
- Entferne das `media`-Replacement aus dem `bs-4-to-5`-Befehl, da es `media` ersetzt und andere Klassen beschädigt

## [v2.4.0] - 2025-07-21

### Hinzugefügt
- **Bootstrap 4 zu 5 Migrations-Befehl** (`bs-4-to-5`) für automatisierte Migration von HTML- und Twig-Templates

## [v2.3.0] - 2025-07-17

### Hinzugefügt
- Konfigurationsunterstützung für den `patchvendor`-Befehl, ermöglicht flexible Patch-Erstellung und verbesserte Diff-Pfad-Behandlung

### Geändert

- Verbesserte Validierung des Vendor-Pfads in Patch-Tests für korrekte relative Pfade und richtiges `a/b`-Präfix-Format
- Patchvendor-Testfälle für genauere Pfaderwartungen und robuste Validierung refaktoriert

## [v2.2.0] - 2025-07-10

### Hinzugefügt
- **TwigBlocks-Befehl** zum Finden von duplizierten Twig-Blöcken in Shopware/Symfony-Projekten
- Rekursives Scannen von `*.html.twig` Dateien mit intelligenter Verzeichnisfilterung
- Erkennung von doppelten Block-Definitionen innerhalb derselben Datei (verhindert Template-Konflikte)
- Mehrere Ausgabeformate: Menschenlesbar, JSON und JUnit XML für CI/CD
- Bitbucket Code Insights Integration mit Reports API und Annotationen
- Datei-Links in Test-Reports für einfache Navigation zu problematischen Dateien
- Intelligente Verzeichnisfilterung (ignoriert `node_modules`, `vendor`, `cache`, etc.)
- CI/CD-Integration mit korrekten Exit-Codes für Automatisierung
- Umfassende Testsuite für TwigBlocks-Funktionalität
- Unterstützung für relative (`.`) und absolute Pfad-Scans
- Nightly-Release-Pipeline mit automatischen täglichen Builds
- Nightly-Docker-Images mit `nightly`, `nightly-amd64`, `nightly-arm64` Tags
- Manuelle Auslösung für Nightly-Releases über GitHub Actions
- Intelligente Nightly-Release-Logik (erstellt nur Release bei neuen Commits)
- Separate GoReleaser-Konfiguration für Nightly-Builds

### Geändert
- Verbesserte Verzeichnis-Scan-Logik für korrekte Behandlung von Root-Verzeichnis (`.`) Pfaden
- Erweiterte JUnit XML Ausgabe mit klickbaren Datei-Links für Bitbucket-Integration

### Behoben
- Problem behoben, bei dem Scannen mit relativem Pfad (`.`) keine Dateien fand
- Docker-Authentifizierung für GitHub Container Registry in CI/CD behoben
- GoReleaser-Konfiguration Kompatibilität mit Version 2 behoben

## [v2.1.0] - 2025-07-10

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

## [v2.0.0] - 2025-07-05

### Hinzugefügt
- Homebrew-Installation über benutzerdefinierten Tap
- Installationsanweisungen für mehrere Paketmanager

### Geändert
- **BREAKING**: `generateUnifiedDiff` refaktoriert, um System-Git-Diff für verbesserte Performance zu verwenden
- Verbesserte Fehlerbehandlung und Performance-Optimierungen

### Behoben
- Performance-Verbesserungen bei der Diff-Generierung

## [v1.0.0] - 2025-07-04

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