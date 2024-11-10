# thwInventoryMerge
Das THW pflegt sein Equipment in THWin. Aus THWin können die entsprechenden Daten für Inventurzwecke als CSV exportiert werden.

Dieses kleine Tool ermöglicht es, das vorhandene Equipment mit den THW-Inventurdaten aus THWin zu mergen. Das Equipment kann dabei z. B. mittels Barcode-Scannern erfasst werden, die die gescannten Codes als CSV-Dateien speichern.

In den THW-Inventurdaten ist ein großer Teil des Equipments als geringwertiges Material (GWM) eingestuft. GWM-Equipment besitzt keine Inventarnummer und ist somit nur schwer automatisiert zu erfassen. Daher ermöglicht das Tool zusätzlich die Erzeugung sogenannter Pseudo-Inventarnummern für geringwertiges Material. Diese Pseudo-Inventarnummern setzen sich aus der Sachnummer und der nächsthöheren Inventarnummer der übergeordneten Ebene zusammen.

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
    "columns": {
        "equipment_layer": "Ebene",
        "equipment_part_number": "Sachnummer",
        "equipment_id": "Inventar Nr",
        "equipment_count_actual": "Bestand IST",
        "equipment_count_target": "Menge"
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

Zur Ausführung öffnet man ein Terminal im `working_dir` und startet dort das Tool.

Zu Beginn einer Inventur ist es notwendig, die CSV-Datei aus THWin einmalig zu initialisieren. Dabei werden Pseudo-Inventarnummern sowie die Spalte "Bestand IST" im CSV (siehe `inventory_csv_file_name`) erstellt.

```bash
?>thwInventoryMerge.exe -s init
```

Anschließend können die Inventurdaten durch die Daten der Scanner ergänzt werden. Dazu reicht es, das Tool entweder per Doppelklick oder im Terminal aufzurufen.

```bash
?>thwInventoryMerge.exe
```

Nach der Ausführung wird im `working_dir` ein Verzeichnis namens `result` erstellt, das eine Datei `result_<timestamp>.csv` enthält. Diese Datei beinhaltet die zusammengeführten Inventurdaten.

Jede weitere Ausführung erzeugt eine neue Datei `result_<timestamp>.csv`.