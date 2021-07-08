package hosts

import (
	"bufio"
	"fmt"
	"io"
)

// ReadLines returns the lines in the specified file.
func ReadLines(r io.Reader) ([]Line, error) {
	var lines []Line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, Line(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return lines, fmt.Errorf("read file: %w", err)
	}

	return lines, nil
}

// WriteLines writes the lines to the specified path.
func WriteLines(w io.Writer, lines []Line) error {
	for _, line := range lines {
		_, err := w.Write([]byte(string(line) + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
