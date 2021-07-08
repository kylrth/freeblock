package hosts_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/kylrth/freeblock/pkg/hosts"
)

func TestLine_IsHostLine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   hosts.Line
		want bool
	}{
		"empty":     {"", false},
		"comment":   {" # this is a comment", false},
		"normal":    {" 0.0.0.0   google.com", true},
		"aliases":   {" 1.1.1.1 google.com google", true},
		"commented": {"#2.2.2.2 twitter.com twitter", true},
		"bad":       {"#1.1.1.1", false},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tc.in.IsHostLine()

			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Error("unexpected output (-want +got):\n" + diff)
			}
		})
	}
}

func TestLine_GetIP(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   hosts.Line
		want string
	}{
		"empty":     {"", ""},
		"comment":   {" # this is a comment", ""},
		"normal":    {" 0.0.0.0   google.com", "0.0.0.0"},
		"aliases":   {" 1.1.1.1 google.com google", "1.1.1.1"},
		"commented": {"#2.2.2.2 twitter.com twitter", "2.2.2.2"},
		"bad":       {"#1.1.1.1", ""},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tc.in.GetIP()

			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Error("unexpected output (-want +got):\n" + diff)
			}
		})
	}
}

func TestLine_SetIP(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   hosts.Line
		ip   string
		want hosts.Line
	}{
		"empty":     {"", "8.5.3.1", ""},
		"comment":   {" # this is a comment", "8.5.3.1", " # this is a comment"},
		"normal":    {" 0.0.0.0   google.com", "8.5.3.1", " 8.5.3.1   google.com"},
		"aliases":   {" 1.1.1.1 google.com google", "1.2.3.4", " 1.2.3.4 google.com google"},
		"commented": {"#2.2.2.2 twitter.com twitter", "1.2.3.4", "#1.2.3.4 twitter.com twitter"},
		"bad":       {"#1.1.1.1", "1.2.3.4", "#1.1.1.1"},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.in.SetIP(tc.ip)

			diff := cmp.Diff(tc.want, tc.in)
			if diff != "" {
				t.Error("unexpected output (-want +got):\n" + diff)
			}
		})
	}
}

func TestLine_Hostnames(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   hosts.Line
		want []string
	}{
		"empty":     {"", nil},
		"comment":   {" # this is a comment", nil},
		"normal":    {" 0.0.0.0   google.com", []string{"google.com"}},
		"aliases":   {" 1.1.1.1 google.com google", []string{"google.com", "google"}},
		"commented": {"#2.2.2.2 twitter.com twitter", []string{"twitter.com", "twitter"}},
		"bad":       {"#1.1.1.1", nil},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tc.in.Hostnames()

			diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty())
			if diff != "" {
				t.Error("unexpected output (-want +got):\n" + diff)
			}
		})
	}
}

func TestLine_Comment(t *testing.T) {
	t.Parallel()

	tests := map[string]modifyInPlaceTest{
		"empty":     {"", "#"},
		"comment":   {" # this is a comment", " # this is a comment"},
		"normal":    {" 0.0.0.0   google.com", " #0.0.0.0   google.com"},
		"aliases":   {" 1.1.1.1 google.com google", " #1.1.1.1 google.com google"},
		"commented": {"#2.2.2.2 twitter.com twitter", "#2.2.2.2 twitter.com twitter"},
		"bad":       {"#1.1.1.1", "#1.1.1.1"},
	}
	runModifyInPlaceTest(t, tests, func(l *hosts.Line) { l.Comment() })
}

func TestLine_Uncomment(t *testing.T) {
	t.Parallel()

	tests := map[string]modifyInPlaceTest{
		"empty":     {"", ""},
		"comment":   {" # this is a comment", " this is a comment"},
		"normal":    {" 0.0.0.0   google.com", " 0.0.0.0   google.com"},
		"aliases":   {" # 1.1.1.1 google.com google", " 1.1.1.1 google.com google"},
		"commented": {"#2.2.2.2 twitter.com twitter", "2.2.2.2 twitter.com twitter"},
		"bad":       {"#1.1.1.1", "1.1.1.1"},
	}
	runModifyInPlaceTest(t, tests, func(l *hosts.Line) { l.Uncomment() })
}

type modifyInPlaceTest struct {
	in   hosts.Line
	want hosts.Line
}

//nolint:thelper // We want to see line numbers here.
func runModifyInPlaceTest(t *testing.T, tests map[string]modifyInPlaceTest, f func(*hosts.Line)) {
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f(&tc.in)

			diff := cmp.Diff(tc.want, tc.in)
			if diff != "" {
				t.Error("unexpected output (-want +got):\n" + diff)
			}
		})
	}
}
