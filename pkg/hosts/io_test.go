package hosts_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kylrth/freeblock/pkg/hosts"
)

func TestReadWriteLines(t *testing.T) {
	t.Parallel()

	b, err := os.ReadFile(filepath.Join("testdata", "sample_hosts"))
	if err != nil {
		t.Fatal(err)
	}

	lines, err := hosts.ReadLines(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	wantLines := []hosts.Line{
		"0.0.0.0 google.com",
		"",
		"# Host addresses",
		"127.0.0.1  localhost",
		"127.0.1.1  devicename",
		"::1        localhost ip6-localhost ip6-loopback",
		"ff02::1    ip6-allnodes",
		"ff02::2    ip6-allrouters",
	}

	diff := cmp.Diff(wantLines, lines)
	if diff != "" {
		t.Error("unexpected output lines (-want +got):\n" + diff)
	}

	var out bytes.Buffer

	err = hosts.WriteLines(&out, lines)
	if err != nil {
		t.Fatal(err)
	}

	diff = cmp.Diff(string(b), out.String())
	if diff != "" {
		t.Error("unexpected output bytes (-want +got):\n" + diff)
	}
}
