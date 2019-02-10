package dbmigrate

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gabstv/dbmigrate/pkg/dbmigrate/mh"
	"github.com/jmoiron/sqlx"
)

var amtagsql = regexp.MustCompile(`-- \[DBMIGRATE:(.+?)\]`)
var amtaggo = regexp.MustCompile(`// \[DBMIGRATE:(.+?)\]`)

var validExts = map[string]bool{
	".sql": true,
	".go":  true,
}

var tagExps = map[FileType]*regexp.Regexp{
	TypeGo:  amtaggo,
	TypeSQL: amtagsql,
}

type MigFile struct {
	Name string
	Type FileType
	// tags tier 1
	UUID   string
	T      time.Time
	Unix   int64
	Author string
	// tags tier 2
	IsNew bool
	Error string
}

var existsQuery = map[string]string{
	"sqlite3": "SELECT 1 FROM sqlite_master WHERE type='table' AND name='db_migrations';",
	"mysql":   "SELECT 1 FROM db_migrations LIMIT 1;",
}

func MigrationTableExists() (bool, error) {
	db, err := mh.EnvConnect()
	if err != nil {
		return false, err
	}
	defer db.Close()
	qq := existsQuery[db.DriverName()]
	if qq == "" {
		return false, fmt.Errorf("no query for this driver: %v", db.DriverName())
	}
	var nn int
	err = db.QueryRowx(qq).Scan(&nn)
	if err != nil && err.Error() == sql.ErrNoRows.Error() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

var createQuery = map[string]string{
	"sqlite3": "CREATE TABLE `db_migrations` ( `migration_id` TEXT NOT NULL UNIQUE, `author_name` TEXT NOT NULL, `created_at` TIMESTAMP NOT NULL, `applied_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP );",
	"mysql":   "CREATE TABLE `db_migrations` ( `migration_id` char(36) NOT NULL DEFAULT '', `author_name` varchar(150) NOT NULL DEFAULT '', `created_at` datetime NOT NULL, `applied_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`migration_id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
}

func CreateMigrationTable() error {
	db, err := mh.EnvConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	qq := createQuery[db.DriverName()]
	if qq == "" {
		return fmt.Errorf("no query for this driver: %v", db.DriverName())
	}
	_, err = db.Exec(qq)
	return err
}

func ListMigrations(rootp string) (newlist, oldlist []*MigFile, err error) {
	onurdered := make([]*MigFile, 0)
	filepath.Walk(rootp, func(path string, info os.FileInfo, err error) error {
		if validExts[filepath.Ext(path)] {
			var t FileType
			switch filepath.Ext(path) {
			case ".sql":
				t = TypeSQL
			case ".go":
				t = TypeGo
			}
			onurdered = append(onurdered, &MigFile{
				Name: path,
				Type: t,
			})
		}
		return nil
	})
	db, err := mh.EnvConnect()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()
	for _, v := range onurdered {
		tagFile(v, db)
	}
	//TODO: put on correct slice based on IsNew
	//TODO: sort by v.Unix ASC
	return
}

func tagFile(mf *MigFile, db *sqlx.DB) error {
	re := tagExps[mf.Type]
	f, err := os.Open(mf.Name)
	if err != nil {
		mf.Error = "open error: " + err.Error()
		return err
	}
	//
	bbff := new(bytes.Buffer)
	io.Copy(bbff, f)
	bbff.WriteRune('\n')
	f.Close()
	//
	br := bufio.NewReader(bbff)
	for ls, err := br.ReadString('\n'); err == nil; ls, err = br.ReadString('\n') {
		if re.MatchString(ls) {
			matches := re.FindStringSubmatch(ls)
			if len(matches) == 2 {
				keyval := strings.Split(matches[1], "=")
				switch strings.ToUpper(keyval[0]) {
				case "UUID":
					if len(keyval) != 2 {
						mf.Error = "invalid UUID"
						return fmt.Errorf("invalid UUID")
					}
					mf.UUID = unquoteIf(keyval[1])
				case "DATE":
					if len(keyval) != 2 {
						mf.Error = "invalid DATE"
						return fmt.Errorf("invalid DATE")
					}
					var terr error
					mf.T, terr = time.Parse("2006-02-01 15:04:05", unquoteIf(keyval[1]))
					if err != nil {
						mf.Error = "invalid date: " + terr.Error()
						return fmt.Errorf("invalid DATE: %v", err.Error())
					}
					mf.Unix = mf.T.Unix()
				case "AUTHOR":
					if len(keyval) != 2 {
						mf.Error = "invalid AUTHOR"
						return fmt.Errorf("invalid AUTHOR")
					}
					mf.Author = unquoteIf(keyval[1])
				}
			}
		}
	}
	n := 0
	if err := db.QueryRowx("SELECT COUNT(*) FROM db_migrations WHERE migration_id=?", mf.UUID).Scan(&n); err != nil {
		mf.Error = err.Error()
		return err
	}
	if n == 0 {
		mf.IsNew = true
	} else {
		mf.IsNew = false
	}
	return nil
}

func unquoteIf(v string) string {
	if v[0] == '"' && v[len(v)-1] == '"' {
		return v[1 : len(v)-1]
	}
	return v
}
