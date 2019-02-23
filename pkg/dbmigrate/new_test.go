package dbmigrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gabstv/dbmigrate/internal/pkg/util"
)

func TestNew(t *testing.T) {
	gr, _ := util.GitRoot()
	if gr == "" {
		t.Fail()
	}
	if err := os.MkdirAll(filepath.Join(gr, "test/tmp/pkg_dbmigrate_new"), 0744); err != nil {
		t.Fatal("could not create test dir", err.Error())
	}
	npath, err := New("migration1", TypeSQL, filepath.Join(gr, "test/tmp/pkg_dbmigrate_new"))
	if err != nil {
		t.Fatal(err.Error())
	}
	if npath != "d" {
		t.Fatal(npath)
	}
	t.Log(npath)
}
