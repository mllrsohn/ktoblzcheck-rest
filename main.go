package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ValidResponse struct {
	Iban          string `json:"iban"`
	IbanFormatted string `json:"iban_formatted"`
	BankName      string `json:"bank_name"`
	BankId        string `json:"bank_id"`
	BankLocation  string `json:"bank_location"`
}

type InValidResponse struct {
	Message string `json:"message"`
	Code    int    `json:"error"`
}

func (e *InValidResponse) Error() string {
	return fmt.Sprintf("Code %s: %s", e.Code, e.Message)
}

func validateIBAN(iban string, country string) (valid bool, err error) {
	iban_check := NewIbanCheck("")
	code := iban_check.Check(iban, country)
	iban_check.Free()

	switch code {
	case 1:
		return false, &InValidResponse{"IBAN is too short to even check", code}
	case 2:
		return false, &InValidResponse{"the 2-character IBAN prefix is unknown", code}
	case 3:
		return false, &InValidResponse{"IBAN has the wrong length", code}
	case 4:
		return false, &InValidResponse{"the country code to check against is unknown", code}
	case 5:
		return false, &InValidResponse{"the IBAN doesn't belong to the country", code}
	case 6:
		return false, &InValidResponse{"Bad IBAN checksum, i.e. the IBAN probably contains a typo", code}
	}

	return true, nil
}

func writeError(w http.ResponseWriter, m error) {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(b)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			writeError(w, &InValidResponse{"Error parsing Request.", 0})
			return
		}
		ibanStr := r.FormValue("iban")
		countryCode := strings.ToUpper(r.FormValue("county"))

		_, err := validateIBAN(ibanStr, countryCode)
		if err != nil {
			writeError(w, err)
			return
		}

		ibanChecker := NewIban(ibanStr)
		Iban := ibanChecker.TransmissionForm()
		IbanFormatted := ibanChecker.Printable()
		ibanChecker.Free()

		iban_check := NewIbanCheck("")
		start, end := iban_check.BicPosition(Iban)

		iban_check.Free()
		blz := Iban[start:end]

		numberCheck := NewDefaultAccountNumberCheck()
		bank, err := numberCheck.FindBank(blz)

		if err != nil {
			writeError(w, err)
			return
		}
		result := &ValidResponse{
			Iban:          Iban,
			IbanFormatted: IbanFormatted,
			BankName:      bank.Name,
			BankId:        bank.BankID,
			BankLocation:  bank.Location,
		}

		b, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		w.Write(b)

	default:
		writeError(w, &InValidResponse{"Sorry, only POSTS methods are supported.", 0})
	}
}

var Port string

func main() {
	flag.StringVar(&Port, "port", "4040", "port to listen on")
	flag.Parse()

	http.HandleFunc("/", handler)
	log.Println("Running at http://localhost:" + Port)
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
