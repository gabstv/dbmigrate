package dbmigrate

// Config file of dbmigrate
type Config struct {
	DefaultDatabase string     `toml:"default_database"`
	Driver          string     `toml:"driver"`
	Migrations      Migrations `toml:"migrations"`
	Databases       []Database `toml:"databases"`
}

// Migrations path and new file config
type Migrations struct {
	Root        string   `toml:"root" mapstructure:"root"`
	DefaultType FileType `toml:"default_type" mapstructure:"default_type"`
}

// Database to apply migrations.
//
// You can set the connection string directly or
// set a file. The file option will read the specified
// file and use it as the Connection String.
type Database struct {
	Name string `toml:"name" mapstructure:"name"`
	CS   string `toml:"cs" mapstructure:"cs"`
	File string `toml:"file" mapstructure:"file"`
}
