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
		"timed":     {" 1.1.1.1 google.com  #freeblock:08-17", true},
		"timed_com": {"#1.1.1.1 google.com  #freeblock:08-17", true},
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
		"timed":     {" 1.1.1.1 google.com  #freeblock:08-17", "1.1.1.1"},
		"timed_com": {"#1.1.1.1 google.com  #freeblock:08-17", "1.1.1.1"},
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
		"timed": {
			" 1.1.1.1 google.com  #freeblock:08-17",
			"1.2.1.2", " 1.2.1.2 google.com  #freeblock:08-17",
		},
		"timed_com": {
			"#1.1.1.1 google.com  #freeblock:08-17",
			"1.4.2.1", "#1.4.2.1 google.com  #freeblock:08-17",
		},
		"bad": {"#1.1.1.1", "1.2.3.4", "#1.1.1.1"},
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
		"timed":     {" 1.1.1.1 google.com  #freeblock:08-17", []string{"google.com"}},
		"timed_com": {"#1.1.1.1 google.com#freeblock:08-17", []string{"google.com"}},
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
		"timed":     {" 1.1.1.1 google.com  #freeblock:08-17", " #1.1.1.1 google.com  #freeblock:08-17"},
		"timed_com": {"#1.1.1.1 google.com  #freeblock:08-17", "#1.1.1.1 google.com  #freeblock:08-17"},
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
		"timed":     {" 1.1.1.1 google.com  #freeblock:08-17", " 1.1.1.1 google.com  #freeblock:08-17"},
		"timed_com": {"#1.1.1.1 google.com  #freeblock:08-17", "1.1.1.1 google.com  #freeblock:08-17"},
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

func TestLine_Timing(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in hosts.Line

		wantA int
		wantB int
	}{
		"empty":          {"", 0, 0},
		"comment":        {" # this is a comment", 0, 0},
		"normal":         {" 0.0.0.0   google.com", 0, 0},
		"aliases":        {" 1.1.1.1 google.com google", 0, 0},
		"commented":      {"#2.2.2.2 twitter.com twitter", 0, 0},
		"timed":          {" 1.1.1.1 google.com  #freeblock:08-17", 8, 17},
		"timed_com":      {"#1.1.1.1 google.com #freeblock:02-24", 2, 24},
		"always_blocked": {"#1.1.1.1 google.com#freeblock:00-24", 0, 24},
		"bad_format":     {"1.1.1.1 google.com #freeblock: 00-24", 0, 0},
		"bad_format2":    {"1.1.1.1 google.com #freelock:00-24", 0, 0},
		"bad_format3":    {"1.1.1.1 google.com #freeblock:00-24-05", 0, 0},
		"bad_format4":    {"1.1.1.1 google.com #freeblock:ab-20", 0, 0},
		"bad_format5":    {"1.1.1.1 google.com #freeblock:05-dd", 0, 0},
		"bad":            {"#1.1.1.1", 0, 0},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotA, gotB := tc.in.Timing()

			diff := cmp.Diff(tc.wantA, gotA)
			if diff != "" {
				t.Error("unexpected first output (-want +got):\n" + diff)
			}
			diff = cmp.Diff(tc.wantB, gotB)
			if diff != "" {
				t.Error("unexpected second output (-want +got):\n" + diff)
			}
		})
	}
}
