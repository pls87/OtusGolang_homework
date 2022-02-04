package notification

import (
	"fmt"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send notifications",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Send!")
	},
}
