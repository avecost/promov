package promov

import (
	"fmt"

	"github.com/avecost/promov/db"
	"github.com/avecost/promov/models/raffle"
	"github.com/avecost/promov/models/jwallgame"
	"github.com/avecost/promov/models/guesttransaction"
	"strings"
	"bufio"
	"os"
	"unicode"
)

type Valid struct {
	AppDb *db.DB
}

func Init(conn string) (*Valid, error) {
	appDb, err := db.Open(conn)
	if err != nil {
		return nil, err
	}

	return &Valid{AppDb: appDb}, nil
}

func (v *Valid) Run(dToValidate string) {

	r, err := raffle.GetAllPendingNonBaccaratResultsByDate(v.AppDb, dToValidate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Processing Raffle Entries :", len(r))
	// prompt to continue
	fmt.Print("Continue [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Println(err)
	}

	char = unicode.ToLower(rune(char))
	if char == 'y' {
		v.process(r)
	}

}

func (v *Valid) process(r []raffle.NBRaffle) {

	for _, rafRow := range r {

		p := strings.Split(rafRow.Provider, "-")

		c1, err := jwallgame.Check(rafRow.JackpotAt, strings.TrimSpace(p[1]), rafRow.Terminal, rafRow.Outlet, rafRow.Game, rafRow.JackpotAmt)
		if err != nil {
			fmt.Println(err)
			return
		}

		c2, err := guesttransaction.Check(rafRow.JackpotAt, rafRow.Cardno, rafRow.Terminal, rafRow.Cashier)
		if err != nil {
			fmt.Println(err)
			return
		}

		var status = 0
		if c1 && c2 {
			status = 2
		} else {
			status = 0
		}

		fmt.Println(rafRow.Id, c1, c2, c1 && c2)

		raffle.UpdateRaffleStatus(v.AppDb, rafRow.Id, status)
	}

}
