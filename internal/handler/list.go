package handler

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"golang.org/x/term"

	"github.com/pHo9UBenaA/cmdbook/internal/config"
	"github.com/pHo9UBenaA/cmdbook/internal/domain"
	"github.com/pHo9UBenaA/cmdbook/pkg/ioutil"
)

func ListCommands(configPath string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	grouped := domain.GroupCommands(cfg.Commands)
	entries := domain.PrepareInteractiveEntries(grouped)

	if len(entries) == 0 {
		fmt.Println("No commands saved")
		return nil
	}

	height, _, _ := term.GetSize(int(os.Stdout.Fd()))
	pageSize := calculatePageSize(height)

	if err := keyboard.Open(); err != nil {
		return fmt.Errorf("failed to initialize keyboard input: %w", err)
	}
	defer keyboard.Close()

	return runInteractiveViewer(entries, pageSize)
}

func calculatePageSize(height int) int {
	pageSize := height - 2
	if pageSize < 1 {
		return 10
	}
	return pageSize
}

func runInteractiveViewer(entries []domain.CommandEntry, pageSize int) error {
	offset := 0
	for {
		printInteractiveView(entries, pageSize, offset)

		char, key, err := keyboard.GetKey()
		if err != nil {
			return fmt.Errorf("failed to get key input: %w", err)
		}

		switch {
		case shouldScrollUp(key, offset):
			offset--
		case shouldScrollDown(key, offset, len(entries), pageSize):
			offset++
		case shouldExit(char, key):
			return nil
		}
	}
}

func printInteractiveView(entries []domain.CommandEntry, pageSize, offset int) {
	fmt.Fprint(os.Stdout, "\033[2J\033[H") // clear display
	printed := ioutil.PrintInteractiveList(entries, pageSize, offset)
	printFooter(offset, printed, len(entries))
}

func printFooter(offset, printed, total int) {
	end := offset + printed
	if end > total {
		end = total
	}
	footer := fmt.Sprintf("\nCommands %d-%d of %d (▲/▼ scroll, q quit)",
		offset+1, end, total)
	fmt.Fprintln(os.Stdout, footer)
}

func shouldScrollUp(key keyboard.Key, offset int) bool {
	return key == keyboard.KeyArrowUp && offset > 0
}

func shouldScrollDown(key keyboard.Key, offset, total, pageSize int) bool {
	return key == keyboard.KeyArrowDown && offset < total-pageSize
}

func shouldExit(char rune, key keyboard.Key) bool {
	return char == 'q' || key == keyboard.KeyEsc
}
