# thwInventoryMerge
Dieses kleine Tool ermöglicht es, erfasstes Equipment mit den THW-Inventurdaten aus THWin zu mergen. Das Equipment wird dabei mittels Barcode-Scannern erfasst, die die gescannten Codes als CSV-Dateien speichern.

## Installation
Das Tool lädt man am einfachsten aus der [Releases-Sektion](https://github.com/mvach/thwInventoryMerge/releases) herunter und legt es in ein beliebiges leeres Verzeichnis (das `working_dir`).

## Konfiguration
Die CSV-Dateien mit den erfassten Barcodes der Scanner sowie die (aus THWin exportierte) Inventur CSV Datei legt man am besten ebenfalls in das `working_dir`.

Zudem erstellt man eine Konfigurationsdatei (`config.json`), die man auch in das `working_dir` legt.

## Hier ist eine Beispielkonfiguration:

```
// config.json
{
    "inventory_csv_file_name": "20240101_Bestand_FGr_N.csv",
    "inventory_csv_config": {
        "equipment_id_column_name": "Inventar Nr",
        "equipment_available_column_name": "Verfügbar"
    }
}
```

### Verzeichnisstruktur

```
C:/
└── Users/
    └── DeinUser/
        └── MeinArbeitsverzeichnis/
            ├── config.json
            ├── 20240101_Bestand_FGr_N.csv
            ├── scanner1.csv
            ├── scanner2.csv
            ├── scanner3.csv
            └── thwInventoryMerge.exe
```

### CSV Beispiel
```csv
// scanner1.csv

0001-S001304
0509-002494
0509-002494
0591-S002360
0591-002781
0591-002781
0591-S002318
...
```

## Ausführung

Befinden sich alle Dateien gemeinsam im `working_dir`, kann die Datei `thwInventoryMerge.exe` einfach per Doppelklick ausgeführt werden.

Nach der Ausführung wird im `working_dir` ein `result`-Verzeichnis erstellt, in dem sich eine Datei namens `result_<timestamp>.csv` befindet. Diese Datei enthält die zusammengeführten Inventurdaten.

Jede weitere Ausführung erzeugt eine neue Datei `result_<timestamp>.csv`.