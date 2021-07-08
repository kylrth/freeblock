package cmds_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kylrth/freeblock/cmd/freeblock/cmds"
)

//nolint:paralleltest // This test modifies package state.
func TestBlockCmd(t *testing.T) {
	hostsFile := filepath.Join("testdata", t.Name())

	// Restore the hosts file once the test is over.
	backupFile(t, hostsFile)

	// Run BlockCmd.
	cmds.BlockCmd.SetArgs([]string{
		"--hosts-file", hostsFile,
		"google.com", "example.com", "internal.example.com",
	})
	if err := cmds.BlockCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Check the file.
	checkWantFile(t, hostsFile)
}

func backupFile(t *testing.T, file string) {
	t.Helper()

	backup, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		e := os.WriteFile(file, backup, 0o600)
		if e != nil {
			t.Fatal(e)
		}
	})
}

func checkWantFile(t *testing.T, file string) {
	t.Helper()

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(file + ".want")
	if err != nil {
		t.Fatal(err)
	}
	diff := cmp.Diff(string(want), string(got))
	if diff != "" {
		t.Error("unexpected output (-want +got):\n" + diff)
	}
}
