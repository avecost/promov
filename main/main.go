package main

import (
	"fmt"
	"os"

	"github.com/avecost/promov"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Invalid parameter")
		return
	}

	connStr := os.Args[1]
	dateTo := os.Args[2]
	v, err := promov.Init(connStr)
	if err != nil {
		fmt.Println("Error: DB Connection")
		return
	}
	// make sure we close the DB
	defer v.AppDb.Close()

	// run the validation
	v.Run(dateTo)
}
