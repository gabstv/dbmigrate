package dbmigrate

import (
	"os"
	"path/filepath"
	"strings"
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
	npath, err := New("migration1", TypeGo, filepath.Join(gr, "test/tmp/pkg_dbmigrate_new"))
	if err != nil {
		t.Fatal(err.Error())
	}
	if !strings.HasSuffix(npath, "migration1/main.go") {
		t.Fatal(npath)
	}
	t.Log(npath)
}
