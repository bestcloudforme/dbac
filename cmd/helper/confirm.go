package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm prints prompt to stderr and returns true only if the user types "y" or "yes".
func Confirm(prompt string) bool {
	fmt.Fprint(os.Stderr, prompt+" [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes"
}
