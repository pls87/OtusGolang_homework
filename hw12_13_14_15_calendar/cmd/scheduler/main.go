package main

import (
	"fmt"
	"os"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/cmd/scheduler/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Couldn't run app: %s", err)
		os.Exit(1)
	}
}
