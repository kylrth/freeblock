package cmds

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// UnblockCmd is a command that unblocks domains in the hosts file.
var UnblockCmd = &cobra.Command{
	Use:   "unblock DOMAIN [DOMAIN...]",
	Short: "unblock domains",
	Long: `Unblock domains by commenting out 0.0.0.0 entries from the hosts file for each
domain.

For blocked hosts with a comment that has another IP address, the domain is
reverted back to to that IP address and the comment is deleted.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := Unblock(args, hostsFile, DefaultNower{}); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	UnblockCmd.Flags().StringVar(
		&hostsFile, "hosts-file", defaultHostsFile, "Change the default hosts file.")
}

// Nower is something that can return the current time. Used for mocking during tests.
type Nower interface {
	Now() time.Time
}

// DefaultNower implements Nower with time.Now.
type DefaultNower struct{}

func (DefaultNower) Now() time.Time {
	return time.Now()
}

// Unblock blocks the domains in the hostsFile.
func Unblock(domains []string, hostsFile string, nower Nower) error {
	lines, err := readLines(hostsFile)
	if err != nil {
		return err
	}

	// Modify the lines in place.
	for i, line := range lines {
		hostnames := line.Hostnames()
		if len(hostnames) == 0 {
			continue
		}

		// See if this line refers to one or more of the domains we want to block.
		for _, hostname := range hostnames {
			if !contains(domains, hostname) {
				continue
			}

			// See if there's a time range we need to respect.
			blockStart, blockEnd := line.Timing()
			now := nower.Now()
			if now.Hour() >= blockStart && now.Hour() < blockEnd {
				return &ErrBlockTiming{i + 1, hostname, blockStart, blockEnd, now}
			}

			oldIP := line.GetIP()
			if oldIP != blockedIP {
				// We don't want to comment this one out, because it's already unblocked.
				continue
			}

			// Check for a commented IP address.
			commentIdx := strings.Index(string(line), commentPrefix) //nolint:gocritic // false positive
			if commentIdx == -1 {
				// There's no IP address to revert to, so we'll just comment out the line.
				line.Comment()
			} else {
				commentFields := strings.Fields(string(line)[commentIdx+len(commentPrefix):])
				if len(commentFields) == 1 && net.ParseIP(commentFields[0]) != nil {
					// We'll revert to the commented IP address.
					line.SetIP(commentFields[0])
					line = line[:commentIdx]
				} else {
					// This comment contains something else besides an IP address.
					line.Comment()
				}
			}
			lines[i] = line

			break
		}
	}

	return writeLines(lines, hostsFile)
}

// ErrBlockTiming is returned when the hosts file has specified that this domain is not to be
// unblocked right now.
type ErrBlockTiming struct {
	lineNum    int
	domain     string
	start, end int
	now        time.Time
}

func (e *ErrBlockTiming) Error() string {
	return fmt.Sprintf(
		"it's %02d:%02d and line %d of the hosts file disallows unblocking %s from %02d:00 to %02d:00",
		e.now.Hour(), e.now.Minute(), e.lineNum, e.domain, e.start, e.end,
	)
}
