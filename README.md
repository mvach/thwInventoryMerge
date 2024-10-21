# thwInventoryMerge
Dieses kleine Tool ermöglicht es, erfasstes Equipment mit den THW-Inventurdaten aus THWin zu mergen. Das Equipment wird dabei mittels Barcode-Scannern erfasst, die die gescannten Codes als CSV-Dateien speichern.

## Installation
Das Tool lädt man am einfachsten aus der [Releases-Sektion](https://github.com/mvach/thwInventoryMerge/releases) herunter und legt es in ein beliebiges leeres Verzeichnis (das `working_dir`).

## Konfiguration
Die CSV-Dateien mit den erfassten Barcodes der Scanner sowie die (aus THWin exportierte) Inventur-Excel-Datei legt man am besten ebenfalls in das `working_dir`.

Zudem erstellt man eine Konfigurationsdatei (`config.json`), die man am einfachsten auch in das `working_dir` legt.

## Hier ist eine Beispielkonfiguration:

```
// config.json
{
    "excel_file_name": "20240101_Bestand.xlsx",
    "excel_config": {
        "worksheet_name": "N",
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
            ├── 20240101_Bestand.xlsx
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

Liegen alle Dateien gemeinsam im `working_dir`, kann thwInventoryMerge.exe einfach per Doppelklick ausgeführt werden.
