package util

import (
	"fmt"
	"os"
	"text/template"
)

const tpltoml = `# Configuration for DBMigrate

default_database = "default"
driver = "sqlite3"

[migrations]
root = "{{.MigrationsPath}}"
default_type = "sql"

[[databases]]
name = "default"
cs = "./default.db"
`

const tpljson = `{
	"migrations": {
		"root": "{{.MigrationsPath}}",
		"default_type": "sql"
	}
}`

type ConfigType string

const (
	CfgTypeTOML ConfigType = "toml"
	CfgTypeJSON ConfigType = "json"
)

type NewConfigDefaults struct {
	MigrationsPath string
}

// NewConfig config file based on templates
func NewConfig(path string, ctype ConfigType, defaults NewConfigDefaults) error {
	var tf string
	switch ctype {
	case CfgTypeTOML:
		tf = tpltoml
	case CfgTypeJSON:
		tf = tpljson
	default:
		return fmt.Errorf("no templates for config type %v", ctype)
	}
	tpl, err := template.New("dbmigrate").Parse(tf)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file \"%v\" already exists", path)
		}
		return err
	}
	defer f.Close()
	return tpl.Execute(f, defaults)
}
