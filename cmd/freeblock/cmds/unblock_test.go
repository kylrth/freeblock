package cmds_test

import (
	"path/filepath"
	"testing"

	"github.com/kylrth/freeblock/cmd/freeblock/cmds"
)

//nolint:paralleltest // This test modifies package state.
func TestUnblockCmd(t *testing.T) {
	hostsFile := filepath.Join("testdata", t.Name())

	// Restore the hosts file once the test is over.
	backupFile(t, hostsFile)

	// Run UnblockCmd.
	cmds.UnblockCmd.SetArgs([]string{
		"--hosts-file", hostsFile,
		"google.com", "example.com", "internal.example.com", "github.com",
	})
	if err := cmds.UnblockCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Check the file.
	checkWantFile(t, hostsFile)
}
