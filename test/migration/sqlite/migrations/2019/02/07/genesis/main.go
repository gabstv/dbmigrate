package main

import (
	"github.com/gabstv/dbmigrate/pkg/dbmigrate/mh"
)

func main() {
	mh.Run(func(tx mh.ExitTx) {
		tx.Exec(`
		CREATE TABLE users (
			Id Int,
			Name Varchar
		  );
		`)
		tx.Exec("INSERT INTO users (Id, Name) VALUES (1, 'Gabs');")
	})
}
