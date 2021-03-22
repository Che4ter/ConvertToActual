# ConvertToActual
Just a small convert script to handle the cumbersome .csv transaction exports from swiss banks

The utility tries to autodetect the bank of the provided export and convert it to the actual format.
It also applies some regex and string formatting to have nicer payer and notes.

**Usage:**

Provide files to convert as arguments. The script supports multiple arguments at once. Each input will be converted and saved 
as _converted.csv
```
convertToActual path/export1.csv path/export2.csv path/export3.csv
```

**Compile:**
```
make compile
```

## Supported Banks
Because exporting transactions isn't easy with swiss banks, below are the instructions how to export your transactions, so 
that the convert script works:

### Raiffeisen Switzerland
- Select Account
- Click the small download icon at the top right of the "Kontoauszug"
- Click "Kontodaten-Download"
- Select Account again
- Select Excel CSV
- Select time under "Vordefinierter Kontoauszug-Zeitraum"
- Select "Ohne Details"
- Click "Herunterladen"

### Revolut
Only CHF and EUR Accounts are supported. Be aware that Actual doesn't support multycurrency
in the same budget file.

- Select Account
- Click the three little dots
- Click "Kontoauszug abrufen"
- Select "Excel" and the time frame
- Copy the file to your computer

### Postfinance

- Go to: Zahlungen -> Bewegungen -> 
- Select duration under "Suchoptionen"
- Activate radio button "Buchungsdetails"
- Click on "Weitere anzeigen" as many times as it takes to load all the results.
- Click "Exportieren"