package dbmigrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gabstv/dbmigrate/internal/pkg/util"
)

func TestNew(t *testing.T) {
	gr, _ := util.GitRoot()
	if err := os.MkdirAll(filepath.Join(gr, "tmp/t_new_test"), 0644); err != nil {
		t.Fatal("could not create test dir", err.Error())
	}
	npath, err := New("migration1", TypeGo, filepath.Join(gr, "tmp/t_new_test"))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(npath)
}
