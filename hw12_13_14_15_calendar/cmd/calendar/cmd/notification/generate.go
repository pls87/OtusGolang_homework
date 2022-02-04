package notification

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate notifications",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generate!")
	},
}
