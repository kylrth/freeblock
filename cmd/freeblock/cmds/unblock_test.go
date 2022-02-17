package cmds_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

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

//nolint:paralleltest // This test modifies package state.
func TestUnblock_fail(t *testing.T) {
	// We want to make sure that Unblock won't unblock when inside a "#freeblock:" time range.

	hostsFile := filepath.Join("testdata", t.Name())

	// Restore the hosts file once the test is over.
	backupFile(t, hostsFile)

	now, err := time.ParseInLocation("15:04", "06:30", time.Local)
	if err != nil {
		t.Fatal(err)
	}

	// Run Unblock and fail.
	err = cmds.Unblock(
		[]string{"google.com", "example.com", "internal.example.com", "github.com"},
		hostsFile, MockNower{now},
	)
	if err == nil {
		t.Fatal("expected error, didn't get one")
	}
	errStr := err.Error()
	want := "it's 06:30 and line 1 of the hosts file disallows" +
		" unblocking google.com from 06:00 to 07:00"
	if diff := cmp.Diff(want, errStr); diff != "" {
		t.Fatal("unexpected error (-want +got):\n" + diff)
	}

	// Check the file.
	checkWantFile(t, hostsFile)
}

//nolint:paralleltest // This test modifies package state.
func TestUnblock_success(t *testing.T) {
	// We want to make sure that Unblock will unblock when outside a "#freeblock:" time range.

	hostsFile := filepath.Join("testdata", t.Name())

	// Restore the hosts file once the test is over.
	backupFile(t, hostsFile)

	now, err := time.ParseInLocation("15:04", "07:30", time.Local)
	if err != nil {
		t.Fatal(err)
	}

	// Run Unblock.
	err = cmds.Unblock(
		[]string{"google.com", "example.com", "internal.example.com", "github.com"},
		hostsFile, MockNower{now},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Check the file.
	checkWantFile(t, hostsFile)
}

// MockNower is a Nower that always returns the same time.Time.
type MockNower struct {
	T time.Time
}

func (m MockNower) Now() time.Time {
	return m.T
}
