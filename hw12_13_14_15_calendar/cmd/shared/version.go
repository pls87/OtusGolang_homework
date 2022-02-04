package shared

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Release   string
	BuildDate string
	GitHash   string
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number of api app",
	Run: func(cmd *cobra.Command, args []string) {
		if err := json.NewEncoder(os.Stdout).Encode(struct {
			Release   string
			BuildDate string
			GitHash   string
		}{
			Release:   Release,
			BuildDate: BuildDate,
			GitHash:   GitHash,
		}); err != nil {
			fmt.Printf("error while decode version info: %v\n", err)
		}
	},
}
