package guesttransaction

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/xml"
)

type gtCheckerResponseEnvelope struct {
	XMLName xml.Name
	Body    gtCheckerResponseBody
}

type gtCheckerResponseBody struct {
	XMLName     xml.Name
	GetResponse gtCheckerResponse `xml:"FN_SEL_Optimus_GuestTransactions_CheckerResponse"`
}

type gtCheckerResponse struct {
	XMLName xml.Name `xml:"FN_SEL_Optimus_GuestTransactions_CheckerResponse"`
	Result  int      `xml:"FN_SEL_Optimus_GuestTransactions_CheckerResult"`
}

func header() string {
	return "<x:Envelope xmlns:x=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:tem=\"http://tempuri.org/\">" +
		"<x:Header/>" +
		"<x:Body>" +
		" <tem:FN_SEL_Optimus_GuestTransactions_Checker>" +
		" <tem:JackpotDate>%s</tem:JackpotDate>" +
		" <tem:CardNo>%s</tem:CardNo>" +
		" <tem:TerminalAccount>%s</tem:TerminalAccount>" +
		" <tem:CashierAccount>%s</tem:CashierAccount>" +
		" <tem:ClientKey>BsCbcZTDkN5pCNtB35QnMGZ2SBzRpQ</tem:ClientKey>" +
		"</tem:FN_SEL_Optimus_GuestTransactions_Checker>" +
		"</x:Body>" +
		"</x:Envelope>"
}

func Check(jackpotOn, cardNo, terminal, cashier string) (bool, error) {
	q := header()
	q = fmt.Sprintf(q, jackpotOn, cardNo, terminal, cashier)

	client := http.Client{}
	resp, err := client.Post("http://116.93.80.34/erf_svc/statistics.asmx?WSDL", "text/xml; charset=utf-8", bytes.NewBufferString(q))
	if err != nil {
		return false, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	r := gtCheckerResponseEnvelope{}
	err = xml.Unmarshal(buf, &r)
	if err != nil {
		return false, err
	}

	return r.Body.GetResponse.Result == 0, nil
}
