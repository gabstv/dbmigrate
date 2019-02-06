package dbmigrate

import "fmt"

// FileType is to specify a new migration file type
type FileType string

const (
	// TypeSQL = "sql" (migration file is a raw .sql)
	TypeSQL FileType = "sql"
	// TypeGo = "go" (migration file is a .go script)
	TypeGo FileType = "go"
)

// New migration file
func New(name string, ftype FileType) error {
	return fmt.Errorf("TODO: this function")
}
