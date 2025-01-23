package ioutil

import (
	"fmt"
	"os"

	"github.com/pHo9UBenaA/cmdbook/internal/domain"
	"golang.org/x/term"
)

const (
	AnsiReset   = "\033[0m"
	AnsiCyan    = "\033[1;36m"
	AnsiGreen   = "\033[1;32m"
	AnsiRed     = "\033[1;31m"
	MaxShortLen = 20
)

func GetTerminalWidth() int {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if width < 40 {
		return 80
	}
	return width
}

func TruncateString(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func PrintInteractiveList(entries []domain.CommandEntry, pageSize, offset int) int {
	width := GetTerminalWidth()
	cmdWidth := width - MaxShortLen - 4
	printed := 0

	for i := offset; i < len(entries) && printed < pageSize; i++ {
		entry := entries[i]
		if entry.Short == "" {
			fmt.Printf("%s%s%s\n", AnsiCyan, entry.Prefix, AnsiReset)
		} else {
			short := entry.Short + ":"
			cmd := TruncateString(entry.Command, cmdWidth)
			fmt.Printf("  %s%-*s%s %-*s\n",
				AnsiGreen, MaxShortLen, short, AnsiReset,
				cmdWidth, cmd)
		}
		printed++
	}
	return printed
}
