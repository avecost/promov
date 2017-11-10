package raffle

import (
	"github.com/avecost/promov/db"
)

type NBRaffle struct {
	Id         int
	Cardno     string
	Terminal   string
	Provider   string
	Outlet     string
	Game       string
	JackpotAt  string
	Cashier    string
	JackpotAmt float32
}

func GetAllPendingNonBaccaratResultsByDate(db *db.DB, dateTo string) ([]NBRaffle, error) {
	rows, err := db.Query("SELECT r.id AS pid, r.cardno AS pcard, r.terminal_acct AS terminal, "+
		" TRIM(p.name) AS pname, TRIM(r.outlet) AS outlet, TRIM(g.name) AS game, DATE_FORMAT(r.jackpot_at, '%m-%d-%Y') AS hiton, "+
		" r.cashier_name AS cashier, r.jackpot_amt AS hitamt "+
		"FROM raffles as r "+
		" LEFT JOIN providers as p ON r.provider_id = p.id "+
		" LEFT JOIN games as g ON r.game_id = g.id "+
		"WHERE r.validated = ? "+
		" AND r.game_id != 355 "+
		" AND r.game_id != 356 "+
		" AND r.game_id != 357 "+
		" AND r.game_id != 358 "+
		" AND r.game_id != 359 "+
		" AND r.game_id != 360 "+
		" AND r.game_id != 365 "+
		" AND r.game_id != 366 "+
		" AND r.game_id != 367 "+
		" AND r.game_id != 368 "+
		" AND r.game_id != 369 "+
		" AND r.jackpot_at = ?", 0, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nbrs []NBRaffle
	for rows.Next() {
		n := NBRaffle{}
		err = rows.Scan(&n.Id, &n.Cardno, &n.Terminal, &n.Provider, &n.Outlet, &n.Game, &n.JackpotAt, &n.Cashier, &n.JackpotAmt)
		if err != nil {
			return nil, err
		}
		nbrs = append(nbrs, n)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return nbrs, nil
}

// status:
// 0 invalid / 1 pending / 2 valid
func UpdateRaffleStatus(db *db.DB, raffleId, status int) error {
	_, err := db.Exec("UPDATE raffles SET valid = ?, validated = ? WHERE id = ?", status, 1, raffleId)
	if err != nil {
		return err
	}

	updateEntriesStatus(db, raffleId, status)

	return nil
}

// status:
// 0 invalid / 1 pending / 2 valid
func updateEntriesStatus(db *db.DB, raffleId, status int) error {
	_, err := db.Exec("UPDATE entries SET valid = ? WHERE player_id = ?", status, raffleId)
	if err != nil {
		return err
	}

	return nil
}
