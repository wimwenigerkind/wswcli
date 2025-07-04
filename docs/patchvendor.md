# PatchVendor Command Documentation

## Deutsch

### Überblick

Der `patchvendor` Command erstellt unified diff Patches für Shopware Vendor-Modifikationen. Das Tool ist speziell für die Arbeit mit Shopware-Projekten entwickelt und generiert Patches mit korrekten `a/` und `b/` Pfaden basierend auf der Vendor/Provider-Struktur.

### Verwendung

```bash
wswcli patchvendor [SOURCE] [PATCHED] [OUTPUT]
```

#### Parameter

- **SOURCE**: Pfad zur originalen, unmodifizierten Vendor-Datei oder -Verzeichnis
- **PATCHED**: Pfad zur modifizierten Version mit Ihren Änderungen
- **OUTPUT**: Pfad, wo die generierte Patch-Datei gespeichert werden soll

#### Modi

##### 1. Direkter Modus
Alle Parameter werden direkt angegeben:

```bash
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php \
                   custom/patches/PluginManager.php \
                   patches/shopware-plugin-manager.patch
```

##### 2. Interaktiver Modus
Ohne Parameter startet der interaktive Modus mit hilfreichen Eingabeaufforderungen:

```bash
wswcli patchvendor
```

Das Tool führt Sie durch die Eingabe und bietet:
- Erklärungen für jeden Parameter
- Beispiele für typische Shopware-Pfade
- Automatische Pfad-Vorschläge für die Ausgabe
- Bestätigungsdialog vor der Verarbeitung

##### 3. Teilweise Parameter
Sie können auch nur einige Parameter angeben:

```bash
# Nur SOURCE angeben
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php

# SOURCE und PATCHED angeben
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php \
                   custom/patches/PluginManager.php
```

### Automatische Pfad-Vorschläge

Wenn Sie den OUTPUT-Parameter weglassen, generiert das Tool automatisch einen strukturierten Pfad:

**Format:** `<Arbeitsverzeichnis>/artifacts/patches/<provider>/<package>/<timestamp>-patch.patch`

**Beispiel:**
- **Input:** `vendor/shopware/core/Framework/Plugin/PluginManager.php`
- **Vorschlag:** `./artifacts/patches/shopware/core/1704369600-patch.patch`

### Unterstützte Dateitypen

Das Tool unterstützt alle gängigen Dateitypen in Shopware-Projekten:
- **PHP**: `.php`
- **JavaScript**: `.js`, `.ts`, `.jsx`, `.tsx`
- **Styling**: `.css`, `.scss`, `.sass`, `.less`
- **Templates**: `.html`, `.twig`
- **Konfiguration**: `.xml`, `.json`, `.yml`, `.yaml`
- **Dokumentation**: `.md`, `.txt`
- **Andere**: `.sql`, `.sh`, `.vue`

### Beispiele

#### Einfache Datei-Modifikation
```bash
# Originale Shopware-Datei
vendor/shopware/core/Framework/Plugin/PluginManager.php

# Ihre modifizierte Version
custom/modifications/PluginManager.php

# Patch erstellen
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php \
                   custom/modifications/PluginManager.php \
                   patches/plugin-manager-modifications.patch
```

#### Verzeichnis-basierte Patches
```bash
# Ganzes Verzeichnis patchen
wswcli patchvendor vendor/shopware/administration/Resources/app/administration/src \
                   custom/administration-modifications \
                   patches/administration-changes.patch
```

#### Interaktiver Workflow
```bash
$ wswcli patchvendor

=== Shopware Vendor Patch Generator ===
This tool creates unified diff patches for Shopware vendor modifications.

SOURCE PATH:
   This is the original, unmodified vendor file or directory.
   Example: vendor/shopware/core/Framework/Plugin/PluginManager.php
   Enter source path: vendor/shopware/core/Framework/Plugin/PluginManager.php

PATCHED PATH:
   This is the modified version of the vendor file or directory.
   It contains your custom changes that should be preserved.
   Example: custom/plugins/MyPlugin/vendor-patches/PluginManager.php
   Enter patched path: custom/modifications/PluginManager.php

OUTPUT PATH:
   This is where the generated patch file will be saved.
   The patch can later be applied using 'git apply' or 'patch' command.
   Suggested: ./artifacts/patches/shopware/core/1704369600-patch.patch
   Enter output path (or press Enter for suggested): 

Summary:
   SOURCE:  vendor/shopware/core/Framework/Plugin/PluginManager.php
   PATCHED: custom/modifications/PluginManager.php
   OUTPUT:  ./artifacts/patches/shopware/core/1704369600-patch.patch

Proceed? (y/N): y
```

### Patch-Anwendung

Die generierten Patches können mit Standard-Tools angewendet werden:

```bash
# Mit Git
git apply patches/shopware-plugin-manager.patch

# Mit dem patch Command
patch -p1 < patches/shopware-plugin-manager.patch
```

### Fehlerbehandlung

Das Tool führt umfangreiche Validierungen durch:

- **Pfad-Existenz**: Überprüft, ob SOURCE und PATCHED existieren
- **Dateityp-Konsistenz**: Stellt sicher, dass beide Pfade Dateien oder Verzeichnisse sind
- **Dateierweiterungen**: Warnt vor unterschiedlichen Dateierweiterungen
- **Pfad-Konflikte**: Verhindert identische SOURCE und PATCHED Pfade
- **Output-Validierung**: Überprüft Output-Pfad auf Konflikte

### Tipps

1. **Backup erstellen**: Erstellen Sie immer Backups vor der Anwendung von Patches
2. **Testen**: Testen Sie Patches in einer Entwicklungsumgebung
3. **Versionierung**: Versionieren Sie Ihre Patches für bessere Nachverfolgung
4. **Dokumentation**: Dokumentieren Sie, welche Änderungen jeder Patch enthält

---

## English

### Overview

The `patchvendor` command creates unified diff patches for Shopware vendor modifications. This tool is specifically designed for working with Shopware projects and generates patches with proper `a/` and `b/` paths based on the vendor/provider structure.

### Usage

```bash
wswcli patchvendor [SOURCE] [PATCHED] [OUTPUT]
```

#### Parameters

- **SOURCE**: Path to the original, unmodified vendor file or directory
- **PATCHED**: Path to the modified version containing your changes
- **OUTPUT**: Path where the generated patch file should be saved

#### Modes

##### 1. Direct Mode
All parameters are provided directly:

```bash
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php \
                   custom/patches/PluginManager.php \
                   patches/shopware-plugin-manager.patch
```

##### 2. Interactive Mode
Without parameters, starts interactive mode with helpful prompts:

```bash
wswcli patchvendor
```

The tool guides you through input and provides:
- Explanations for each parameter
- Examples of typical Shopware paths
- Automatic output path suggestions
- Confirmation dialog before processing

##### 3. Partial Parameters
You can also provide only some parameters:

```bash
# Only SOURCE specified
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php

# SOURCE and PATCHED specified
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php \
                   custom/patches/PluginManager.php
```

### Automatic Path Suggestions

When you omit the OUTPUT parameter, the tool automatically generates a structured path:

**Format:** `<working-directory>/artifacts/patches/<provider>/<package>/<timestamp>-patch.patch`

**Example:**
- **Input:** `vendor/shopware/core/Framework/Plugin/PluginManager.php`
- **Suggestion:** `./artifacts/patches/shopware/core/1704369600-patch.patch`

### Supported File Types

The tool supports all common file types in Shopware projects:
- **PHP**: `.php`
- **JavaScript**: `.js`, `.ts`, `.jsx`, `.tsx`
- **Styling**: `.css`, `.scss`, `.sass`, `.less`
- **Templates**: `.html`, `.twig`
- **Configuration**: `.xml`, `.json`, `.yml`, `.yaml`
- **Documentation**: `.md`, `.txt`
- **Others**: `.sql`, `.sh`, `.vue`

### Examples

#### Simple File Modification
```bash
# Original Shopware file
vendor/shopware/core/Framework/Plugin/PluginManager.php

# Your modified version
custom/modifications/PluginManager.php

# Create patch
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php \
                   custom/modifications/PluginManager.php \
                   patches/plugin-manager-modifications.patch
```

#### Directory-based Patches
```bash
# Patch entire directory
wswcli patchvendor vendor/shopware/administration/Resources/app/administration/src \
                   custom/administration-modifications \
                   patches/administration-changes.patch
```

#### Interactive Workflow
```bash
$ wswcli patchvendor

=== Shopware Vendor Patch Generator ===
This tool creates unified diff patches for Shopware vendor modifications.

SOURCE PATH:
   This is the original, unmodified vendor file or directory.
   Example: vendor/shopware/core/Framework/Plugin/PluginManager.php
   Enter source path: vendor/shopware/core/Framework/Plugin/PluginManager.php

PATCHED PATH:
   This is the modified version of the vendor file or directory.
   It contains your custom changes that should be preserved.
   Example: custom/plugins/MyPlugin/vendor-patches/PluginManager.php
   Enter patched path: custom/modifications/PluginManager.php

OUTPUT PATH:
   This is where the generated patch file will be saved.
   The patch can later be applied using 'git apply' or 'patch' command.
   Suggested: ./artifacts/patches/shopware/core/1704369600-patch.patch
   Enter output path (or press Enter for suggested): 

Summary:
   SOURCE:  vendor/shopware/core/Framework/Plugin/PluginManager.php
   PATCHED: custom/modifications/PluginManager.php
   OUTPUT:  ./artifacts/patches/shopware/core/1704369600-patch.patch

Proceed? (y/N): y
```

### Applying Patches

Generated patches can be applied using standard tools:

```bash
# With Git
git apply patches/shopware-plugin-manager.patch

# With patch command
patch -p1 < patches/shopware-plugin-manager.patch
```

### Error Handling

The tool performs comprehensive validations:

- **Path Existence**: Checks if SOURCE and PATCHED exist
- **Type Consistency**: Ensures both paths are files or directories
- **File Extensions**: Warns about different file extensions
- **Path Conflicts**: Prevents identical SOURCE and PATCHED paths
- **Output Validation**: Checks output path for conflicts

### Best Practices

1. **Create Backups**: Always create backups before applying patches
2. **Test First**: Test patches in a development environment
3. **Version Control**: Version your patches for better tracking
4. **Documentation**: Document what changes each patch contains
5. **Atomic Changes**: Keep patches focused on specific changes
6. **Review**: Always review generated patches before applying

### Advanced Usage

#### Batch Processing
```bash
# Process multiple files
for file in vendor/shopware/core/Framework/Plugin/*.php; do
    if [ -f "custom/modifications/$(basename "$file")" ]; then
        wswcli patchvendor "$file" \
                          "custom/modifications/$(basename "$file")" \
                          "patches/$(basename "$file" .php).patch"
    fi
done
```

#### Integration with CI/CD
```bash
# Generate patches in CI pipeline
wswcli patchvendor vendor/shopware/core \
                   custom/shopware-modifications \
                   artifacts/shopware-core-modifications.patch

# Validate patches can be applied
git apply --check artifacts/shopware-core-modifications.patch
```

### Troubleshooting

#### Common Issues

1. **"source path does not exist"**
   - Verify the source file/directory exists
   - Check for typos in the path

2. **"different extensions"**
   - Ensure source and patched files have the same extension
   - Use `--force` flag if intentional (future feature)

3. **"output path exists and is a directory"**
   - Specify a file path for output, not a directory
   - Use a `.patch` extension for clarity

4. **Empty patches**
   - Verify the files actually differ
   - Check file permissions and readability

#### Debug Mode
```bash
# Enable verbose output (future feature)
wswcli patchvendor --verbose source.php patched.php output.patch
```

### Integration Examples

#### With Shopware Plugin Development
```bash
# Create patches for plugin-specific vendor modifications
wswcli patchvendor vendor/shopware/core/Framework/Plugin/PluginManager.php \
                   plugins/MyPlugin/patches/PluginManager.php \
                   plugins/MyPlugin/patches/core-plugin-manager.patch
```

This documentation provides comprehensive guidance for using the `patchvendor` command effectively in Shopware development workflows.