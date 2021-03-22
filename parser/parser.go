package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var HEADER_ACTUAL = []string{"Date", "Payee", "Notes", "Amount"}

var TIME_LAYOUT_ACTUAL = "2006-01-02"
var TIME_LAYOUT_RAIFFEISEN = "2006-01-02 15:04:05.0"
var TIME_LAYOUT_REVOLUT = "02.01.2006 "
var TIME_LAYOUT_POSTFINANCE = "2006-01-02"

func ParsePostfinance(data [][]string) ([][]string, error) {
	var output = [][]string{}
	output = append(output, HEADER_ACTUAL)
	kaufRegex := regexp.MustCompile(`KAUF\/DIENSTLEISTUNG VOM \d\d\.\d\d\.\d\d\d\d KARTEN NR\. XXXX\d{4} `)
	senderReferenz := regexp.MustCompile(`SENDER REFERENZ:`)
	mitteilung := regexp.MustCompile(`MITTEILUNGEN:`)

	lastRowNumber := len(data) - 3
	for _, element := range data[5:lastRowNumber] {
		valuta_date, err := time.Parse(TIME_LAYOUT_POSTFINANCE, element[0])
		if err != nil {
			return nil, err
		}

		note := ""
		amount := ""
		payee := element[1]
		payee = strings.Replace(payee, "GIRO BANK ", "", 1)
		payee = strings.Replace(payee, "GIRO POST ", "", 1)
		payee = strings.Replace(payee, "ESR ", "", 1)
		payee = strings.Replace(payee, "GUTSCHRIFT VON FREMDBANK AUFTRAGGEBER:", "", 1)

		if strings.HasPrefix(payee, "(Kontoübertrag)") {
			note = "Kontoübertrag"
			payee = strings.Replace(payee, "(Kontoübertrag) ", "", 1)
		}

		index := kaufRegex.FindStringIndex(payee)
		if index != nil {
			payee = payee[index[1]:]
		}

		index = senderReferenz.FindStringIndex(payee)
		if index != nil {
			note = payee[index[0]:]
			payee = payee[:index[0]]
		}

		index = mitteilung.FindStringIndex(payee)
		if index != nil {
			note = payee[index[0]:]
			payee = payee[:index[0]]
		}

		if strings.HasPrefix(payee, "Migros ") {
			note = payee
			payee = "Migros"
		}

		if strings.HasPrefix(payee, "Coop") {
			note = payee
			payee = "Coop"
		}

		if element[2] != "" {
			amount = element[2]
		} else {
			amount = element[3]
		}

		note = strings.Trim(note, " ")
		payee = strings.Trim(payee, " ")
		output = append(output, []string{valuta_date.Format(TIME_LAYOUT_ACTUAL), payee, note, amount})
	}

	return output, nil
}

func ParseRaiffeisen(data [][]string) ([][]string, error) {
	var output = [][]string{}
	output = append(output, HEADER_ACTUAL)
	einkaufRegex := regexp.MustCompile(`\s\d\d.\d\d.\d\d\d\d,`)
	sepaRegex := regexp.MustCompile("\\(SEPA\\) [A-Z]{3} [+-]?([0-9]*[.])?[0-9]+, Umrechnungskurs [+-]?([0-9]*[.])?[0-9]+ ")
	twintRegex := regexp.MustCompile(` TWINT Nr. [+-]?([0-9]*[.])?[0-9]+`)
	bancomatRegex := regexp.MustCompile(`Bancomat Bezug `)

	for _, element := range data[1:] {
		valuta_date, err := time.Parse(TIME_LAYOUT_RAIFFEISEN, element[5])
		if err != nil {
			return nil, err
		}

		note := ""
		payee := element[2]
		payee = strings.Replace(payee, "E-Banking Auftrag (eBill) ", "", 1)
		payee = strings.Replace(payee, "E-Banking Auftrag an ", "", 1)
		payee = strings.Replace(payee, "E-Banking Auftrag ", "", 1)
		payee = strings.Replace(payee, "Gutschrift TWINT von ", "", 1)
		payee = strings.Replace(payee, "Einkauf TWINT, ", "", 1)
		payee = strings.Replace(payee, "Überweisung TWINT an , ", "", 1)
		payee = strings.Replace(payee, "Überweisung TWINT an ", "", 1)
		payee = strings.Replace(payee, "Gutschrift TWINT ", "", 1)
		payee = strings.Replace(payee, "Gutschrift ", "", 1)
		payee = strings.Replace(payee, "Postvergütung von ", "", 1)

		if strings.HasPrefix(payee, "Einkauf ") {
			index := einkaufRegex.FindStringIndex(payee)
			payee = payee[8:index[0]] //Einkauf Unternehmen 12.01.2021,
		}

		if strings.HasPrefix(payee, "E-Banking Dauerauftrag an ") {
			note = "Dauerauftrag"
			payee = strings.Replace(payee, "E-Banking Dauerauftrag an ", "", 1)
		}

		if strings.HasPrefix(payee, "(Kontoübertrag)") {
			note = "Kontoübertrag"
			payee = strings.Replace(payee, "(Kontoübertrag) ", "", 1)
		}

		index := sepaRegex.FindStringIndex(payee)
		if index != nil {
			note = payee[0:index[1]]
			payee = payee[index[1]:]
		}

		index = twintRegex.FindStringIndex(payee)
		if index != nil {
			note = payee[index[0]:]
			payee = payee[0:index[0]]
		}

		if strings.HasPrefix(payee, "Migros ") {
			note = payee
			payee = "Migros"
		}

		if strings.HasPrefix(payee, "Coop") {
			note = payee
			payee = "Coop"
		}

		if strings.HasPrefix(payee, "Post CH AG") {
			note = payee
			payee = "Post CH AG"
		}

		index = bancomatRegex.FindStringIndex(payee)
		if index != nil {
			note = payee[index[1]:]
			payee = "Bancomat Bezug"
		}

		amount := element[3]
		note = strings.Trim(note, " ")
		payee = strings.Trim(payee, " ")
		output = append(output, []string{valuta_date.Format(TIME_LAYOUT_ACTUAL), payee, note, amount})
	}

	return output, nil
}

func ParseRevolut(data [][]string) ([][]string, error) {
	var output = [][]string{}
	output = append(output, HEADER_ACTUAL)
	gekauftverkauftRegex := regexp.MustCompile("\\w{3} (mit|an) \\w{3} (gekauft|verkauft)")
	wechselkursRegex := regexp.MustCompile(" Wechselkurs.*")

	for _, element := range data[1:] {
		date, err := time.Parse(TIME_LAYOUT_REVOLUT, element[0])
		if err != nil {
			return nil, err
		}

		note := ""
		element[2] = strings.Trim(element[2], " ")
		element[3] = strings.Trim(element[3], " ")
		element[4] = strings.Trim(element[3], " ")
		element[5] = strings.Trim(element[3], " ")
		element[2] = strings.Replace(element[2], "’", "", -1)
		element[3] = strings.Replace(element[3], "’", "", -1)
		element[4] = strings.Replace(element[2], "’", "", -1)
		element[5] = strings.Replace(element[3], "’", "", -1)

		payee := strings.Trim(element[1], " ")
		payee = strings.Replace(payee, "Einkauf ", "", 1)
		payee = strings.Replace(payee, "Zahlung von ", "", 1)
		payee = strings.Replace(payee, "Von ", "", 1)
		payee = strings.Replace(payee, "From ", "", 1)
		payee = strings.Replace(payee, "To ", "", 1)
		payee = strings.Replace(payee, "An ", "", 1)
		payee = strings.Replace(payee, "Payment from ", "", 1)

		if strings.HasPrefix(payee, "Bargeld am ") {
			note = payee
			payee = "Bancomat Bezug"
		}

		if strings.HasPrefix(payee, "Rückerstattung von ") {
			payee = strings.Replace(payee, "Rückerstattung von ", "", 1)
			note = "Rückerstattung"
		}

		var amount float64

		if element[2] != "" {
			amount, err = strconv.ParseFloat(element[2], 64)
			if err != nil {
				return nil, err
			}

			amount = amount * -1
		} else {
			amount, err = strconv.ParseFloat(element[3], 64)
			if err != nil {
				return nil, err
			}
		}

		index := gekauftverkauftRegex.FindStringIndex(payee)
		if index != nil {
			note = payee[index[1]:]
			payee = payee[:index[1]]
		}

		index = wechselkursRegex.FindStringIndex(payee)
		if index != nil {
			note = payee[index[0]:]
			payee = payee[:index[0]]
		}

		if note == "" {
			note = element[8]
		}

		note = strings.Trim(note, " ")
		payee = strings.Trim(payee, " ")
		output = append(output, []string{date.Format(TIME_LAYOUT_ACTUAL), payee, note, strconv.FormatFloat(amount, 'f', 2, 64)})
	}

	return output, nil
}
