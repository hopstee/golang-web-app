package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func PromtLine(promt string) string {
	fmt.Print(promt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func PromtPassword(promt string) string {
	fmt.Print(promt)
	bytePw, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePw)
}
