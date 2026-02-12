package cli

import "strings"

func oneLine(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " | ")
	if len(s) > 140 {
		return s[:137] + "..."
	}
	return s
}

func joinArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}
	s := args[0]
	for i := 1; i < len(args); i++ {
		s += " " + args[i]
	}
	return s
}
