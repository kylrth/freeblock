package cmds_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"

	"github.com/kylrth/freeblock/cmd/freeblock/cmds"
)

//nolint:paralleltest // This test modifies package state.
func TestOpenCmd(t *testing.T) {
	hostsFile := filepath.Join("testdata", "TestUnblockCmd")

	// At the end of the test, make sure the file ends up the same as it started.
	fileShouldNotChange(t, hostsFile)

	// Run OpenCmd in a separate routine.
	osSignals := make(chan os.Signal, 1)
	var g errgroup.Group
	g.Go(func() error {
		return cmds.Open(
			[]string{"google.com", "example.com", "internal.example.com", "github.com"},
			hostsFile,
			osSignals,
		)
	})

	// Wait until the file has been changed, and then check the file.
	time.Sleep(10 * time.Millisecond)
	checkWantFile(t, hostsFile)

	// Kill the routine and check the error.
	osSignals <- os.Interrupt
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}

func fileShouldNotChange(t *testing.T, file string) {
	t.Helper()

	backup, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		final, e := os.ReadFile(file)
		if e != nil {
			t.Fatal(e)
		}
		diff := cmp.Diff(string(backup), string(final))
		if diff != "" {
			t.Error("hosts file ended up changed (-want +got):\n" + diff)
		}
	})
}
