package ioutil

import (
	"fmt"
	"os"

	"github.com/pHo9UBenaA/cmdbook/internal/constant"
	"github.com/pHo9UBenaA/cmdbook/internal/domain"
	"golang.org/x/term"
)

const (
	AnsiReset = "\033[0m"
	AnsiCyan  = "\033[1;36m"
	AnsiGreen = "\033[1;32m"
	AnsiRed   = "\033[1;31m"
)

func PrintInteractiveList(entries []domain.CommandEntry, pageSize, offset int) int {
	width := getTerminalWidth()
	cmdWidth := width - constant.MaxShortLen - 4
	printed := 0

	for i := offset; i < len(entries) && printed < pageSize; i++ {
		entry := entries[i]
		if entry.Short == "" {
			fmt.Printf("%s%s%s\n", AnsiCyan, entry.Prefix, AnsiReset)
		} else {
			short := entry.Short + ":"
			cmd := truncateString(entry.Command, cmdWidth)
			fmt.Printf("  %s%-*s%s %-*s\n",
				AnsiGreen, constant.MaxShortLen, short, AnsiReset,
				cmdWidth, cmd)
		}
		printed++
	}
	return printed
}

func getTerminalWidth() int {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if width < 40 {
		return 80
	}
	return width
}

func truncateString(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
