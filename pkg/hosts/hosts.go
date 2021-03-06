// Package hosts provides tools for handling the hosts file, usually found at /etc/hosts.
package hosts

import (
	"net"
	"strconv"
	"strings"
	"unicode"
)

// Line is a single line in a hosts file. It may contain a host line, or a comment.
type Line string

// IsHostLine returns whether this Line contains a (possibly commented-out) host line.
func (l Line) IsHostLine() bool {
	s := strings.TrimSpace(string(l))

	if s == "" {
		return false
	}

	if s[0] != '#' {
		// We'll assume a non-empty, uncommented line is a host line.
		return true
	}

	// This is a comment line, but it might be a commented host line.
	f := strings.Fields(strings.TrimSpace(s[1:]))
	if len(f) < 2 {
		return false
	}

	return isIPAddress(f[0])
}

func isIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

// GetIP returns the IP of a host line. The empty string is returned if l is not a host line.
func (l Line) GetIP() string {
	if !l.IsHostLine() {
		return ""
	}

	s := strings.TrimSpace(string(l))

	if s[0] == '#' {
		s = s[1:]
	}

	return strings.Fields(strings.TrimSpace(s))[0]
}

// SetIP sets the IP of a host line to the provided string. Nothing happens if l is not a host line.
func (l *Line) SetIP(ip string) {
	if !l.IsHostLine() {
		return
	}

	s := string(*l)

	// Find the first non-space character.
	startOfIP := strings.IndexFunc(s, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	if s[startOfIP] == '#' {
		startOfIP++
	}
	endOfIP := startOfIP + strings.IndexFunc(s[startOfIP:], unicode.IsSpace)

	*l = Line(s[:startOfIP] + ip + s[endOfIP:])
}

// Hostnames returns the domain name and aliases referred to in the Line. If the Line is not a host
// line, the slice will be empty.
func (l Line) Hostnames() []string {
	s := strings.TrimSpace(string(l))

	if s == "" {
		return nil
	}
	if s[0] == '#' {
		s = s[1:]
	}

	f := strings.Fields(strings.TrimSpace(s))
	if len(f) < 1 || !isIPAddress(f[0]) {
		return nil
	}

	for i := 1; i < len(f); i++ {
		if f[i][0] == '#' {
			return f[1:i]
		}

		idx := strings.Index(f[i], "#")
		if idx != -1 {
			// Comment starts here.
			return append(f[1:i], f[i][:idx])
		}
	}

	return f[1:]
}

// Comment adds '#' to the beginning of the line, if it's not already there.
func (l *Line) Comment() {
	s := strings.TrimSpace(string(*l))
	if s != "" && s[0] == '#' {
		return
	}

	numSpaces := strings.Index(string(*l), s) //nolint:gocritic // false positive

	*l = Line(string(*l)[:numSpaces] + "#" + s)
}

// Uncomment removes '#' from the beginning of the line if it's there.
func (l *Line) Uncomment() {
	s := strings.TrimSpace(string(*l))
	if s == "" || s[0] != '#' {
		return
	}

	numSpaces := strings.Index(string(*l), s) //nolint:gocritic // false positive

	*l = Line(string(*l)[:numSpaces] + strings.TrimLeftFunc(s[1:], unicode.IsSpace))
}

// Timing returns the times when the line shouldn't be unblocked. If that isn't specified or if this
// isn't a host line at all, zeros are returned.
func (l *Line) Timing() (start, end int) {
	s := strings.TrimSpace(string(*l))

	if s == "" {
		return 0, 0
	}
	if s[0] == '#' {
		s = s[1:]
	}

	f := strings.Fields(strings.TrimSpace(s))
	if len(f) < 2 || !isIPAddress(f[0]) {
		return 0, 0
	}

	for i := 1; i < len(f); i++ {
		idx := strings.Index(f[i], "#")
		if idx == -1 {
			continue
		}

		s = f[i][idx:]
		const prefix = "#freeblock:"
		if len(s) < len(prefix) || s[:len(prefix)] != prefix {
			return 0, 0
		}
		s = s[len(prefix):]

		hours := strings.Split(s, "-")
		if len(hours) != 2 {
			return 0, 0
		}

		start, err := strconv.Atoi(hours[0])
		if err != nil {
			return 0, 0
		}
		end, err := strconv.Atoi(hours[1])
		if err != nil {
			return 0, 0
		}

		return start, end
	}

	return 0, 0
}
