package helpers

import (
	"fmt"
	"github.com/spf13/cobra"
)

func SplitString(s string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func FlagLookup(cmd *cobra.Command, flagName string) (string, error) {
	flag := cmd.Flags().Lookup(flagName)

	if flag != nil && flag.Value != nil {
		return flag.Value.String(), nil
	}

	return "", fmt.Errorf("flag '%s' not found", flagName)
}
