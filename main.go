package main

import (
	"ConvertToActual/parser"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

var HEADER_RAIFFEISEN = []string{"IBAN", "Booked At", "Text", "Credit/Debit Amount", "Balance", "Valuta Date"}
var HEADER_REVOLUT_CH = []string{"Abschlussdatum ", " Beschreibung ", " Ausgezahlt (CHF) ", " Eingezahlt (CHF) ", " Umtausch aus", " Umtausch in", " Kontostand (CHF)", " Kategorie", " Anmerkungen"}
var HEADER_REVOLUT_EUR = []string{"Abschlussdatum ", " Beschreibung ", " Ausgezahlt (EUR) ", " Eingezahlt (EUR) ", " Umtausch aus", " Umtausch in", " Kontostand (EUR)", " Kategorie", " Anmerkungen"}
var HEADER_POSTFINANCE_CHF = []string{"Buchungsdatum", "Avisierungstext", "Gutschrift in CHF", "Lastschrift in CHF", "Valuta", "Saldo in CHF"}

type BankType string

const (
	Raiffeisen      BankType = "Raiffeisen"
	Revolut_CHF     BankType = "Revolut_CHF"
	Revolut_EUR     BankType = "Revolut_EUR"
	Postfinance_CHF BankType = "Postfinance_CHF"
)

var Version = "Development Build"

func main() {
	fmt.Println("Welcome to ConverToActual version: ", Version)

	args := os.Args[1:]
	if len(args) > 0 {
		for _, arg := range args {

			results, err := readCSV(arg)
			if err != nil {
				log.Fatalln(err)
			}

			bankType, err := detectBank(results)
			if err != nil {
				log.Fatalln(err)
			}

			var output = [][]string{}
			fmt.Printf("Detected %s Format", bankType)

			if bankType == Raiffeisen {
				output, err = parser.ParseRaiffeisen(results)
				if err != nil {
					log.Fatalln(err)
				}
			} else if bankType == Revolut_CHF || bankType == Revolut_EUR {
				output, err = parser.ParseRevolut(results)
				if err != nil {
					log.Fatalln(err)
				}
			} else if bankType == Postfinance_CHF {
				output, err = parser.ParsePostfinance(results)
				if err != nil {
					log.Fatalln(err)
				}
			}

			err = writeToFile(output, arg)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}else{
		fmt.Println("Error: Please provides file to convert as arguments. E.g. convertoactual filename.csv")
	}
}

func readCSV(path string) ([][]string, error) {
	charset, err := detectCharset(path)
	csvFile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var reader *csv.Reader
	if charset != "UTF-8" {
		reader = csv.NewReader(charmap.ISO8859_15.NewDecoder().Reader(csvFile))
	} else {
		reader = csv.NewReader(csvFile)
	}

	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, err
}

func detectCharset(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 32<<10)
	size, _ := io.ReadFull(f, buffer)
	input := buffer[:size]
	var detector = chardet.NewTextDetector()

	result, err := detector.DetectBest(input)

	return result.Charset, nil
}

func detectBank(content [][]string) (BankType, error) {
	if len(content) == 0 {
		return "", errors.New("Empty file")
	}

	if Equal(content[0], HEADER_RAIFFEISEN) {
		return Raiffeisen, nil
	}

	if Equal(content[0], HEADER_REVOLUT_CH) {
		return Revolut_CHF, nil
	}

	if Equal(content[0], HEADER_REVOLUT_EUR) {
		return Revolut_EUR, nil
	}

	if content[0][0] == "Datum von:" && Equal(content[4], HEADER_POSTFINANCE_CHF) {
		return Postfinance_CHF, nil
	}
	return "", nil
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func writeToFile(data [][]string, fileName string) error {
	fileName = strings.TrimSuffix(fileName, path.Ext(fileName))

	file, err := os.Create(fileName + "_converted.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}
