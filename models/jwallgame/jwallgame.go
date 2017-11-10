package jwallgame

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/xml"
)

type jwAllGamesCheckerResponseEnvelope struct {
	XMLName		xml.Name
	Body  		jwAllGamesCheckerResponseBody
}

type jwAllGamesCheckerResponseBody struct {
	XMLName     xml.Name
	GetResponse jwAllGamesCheckerResponse `xml:"FN_SEL_Celebra_RPT_JackpotWinnings_All_Games_CheckerResponse"`
}

type jwAllGamesCheckerResponse struct {
	XMLName 	xml.Name	`xml:"FN_SEL_Celebra_RPT_JackpotWinnings_All_Games_CheckerResponse"`
	Result		int			`xml:"FN_SEL_Celebra_RPT_JackpotWinnings_All_Games_CheckerResult"`
}

func header() string {
	return "<x:Envelope xmlns:x=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:tem=\"http://tempuri.org/\">" +
		"<x:Header/>" +
		"<x:Body>" +
		"<tem:FN_SEL_Celebra_RPT_JackpotWinnings_All_Games_Checker>" +
		" <tem:JackpotDate>%s</tem:JackpotDate>" +
		" <tem:Provider>%s</tem:Provider>" +
		" <tem:Terminal>%s</tem:Terminal>" +
		" <tem:Outlet>%s</tem:Outlet>" +
		" <tem:ProgressiveName>%s</tem:ProgressiveName>" +
		" <tem:Payout>%g</tem:Payout>" +
		" <tem:ClientKey>BsCbcZTDkN5pCNtB35QnMGZ2SBzRpQ</tem:ClientKey>" +
		"</tem:FN_SEL_Celebra_RPT_JackpotWinnings_All_Games_Checker>" +
		"</x:Body>" +
		"</x:Envelope>"
}

func Check(jackpotOn, provider, terminal, outlet, gameName string, jackpot float32) (bool, error) {
	q := header()
	q = fmt.Sprintf(q, jackpotOn, provider, terminal, outlet, gameName, jackpot)

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

	r := jwAllGamesCheckerResponseEnvelope{}
	err = xml.Unmarshal(buf, &r)
	if err != nil {
		return false, err
	}

	return r.Body.GetResponse.Result == 0, nil
}