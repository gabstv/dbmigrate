package main

import (
	"github.com/gabstv/dbmigrate/pkg/dbmigrate/mh"
)

// DO NOT REMOVE THE COMMENTS BELOW
// [DBMIGRATE:UUID="00b66979-4426-41dc-8918-398bd588e7dc"]
// [DBMIGRATE:DATE="2019-02-07 22:00:10"]
// [DBMIGRATE:AUTHOR="gabs"]

func main() {
	mh.Run(func(tx mh.Mtx) {
		tx.Exec(`
		CREATE TABLE users (
			Id Int,
			Name Varchar
		  );
		`)
		tx.Exec("INSERT INTO users (Id, Name) VALUES (1, 'Gabs');")
	})
}
