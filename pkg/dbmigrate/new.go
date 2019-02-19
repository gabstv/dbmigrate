package dbmigrate

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/google/uuid"
)

// FileType is to specify a new migration file type
type FileType string

const (
	// TypeSQL = "sql" (migration file is a raw .sql)
	TypeSQL FileType = "sql"
	// TypeGo = "go" (migration file is a .go script)
	TypeGo FileType = "go"
)

var templateGo = `package main

import (
	"github.com/gabstv/dbmigrate/pkg/dbmigrate/mh"
)

// DO NOT REMOVE THE COMMENTS BELOW
// [DBMIGRATE:UUID="{{.UUID}}"]
// [DBMIGRATE:DATE="{{.DATE}}"]
// [DBMIGRATE:AUTHOR="{{.AUTHOR}}"]

func main() {
	mh.Run(func(tx mh.Mtx) {
		// migration code goes here
		tx.Exec("SELECT 1;")
	})
}
`

var templateSQL = `-- DO NOT REMOVE THE COMMENTS BELOW
-- [DBMIGRATE:UUID="{{.UUID}}"]
-- [DBMIGRATE:DATE="{{.DATE}}"]
-- [DBMIGRATE:AUTHOR="{{.AUTHOR}}"]

-- migration code goes here
`

// New migration file
func New(name string, ftype FileType, migrationsRoot string) (string, error) {
	id := uuid.New()
	var tf string
	switch ftype {
	case TypeGo:
		tf = templateGo
	case TypeSQL:
		tf = templateSQL
	}
	tpl, err := template.New("new_migration").Parse(tf)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	tpl.Execute(buf, map[string]interface{}{
		"UUID": id.String(),
		"DATE": time.Now().Format("2006-01-02 15:04:05"),
	})
	return id.String(), fmt.Errorf("TODO: this function")
}
