package main

import (
	"fmt"
	"os"

	calendarcmd "github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/cmd/calendar/cmd"
)

func main() {
	if err := calendarcmd.Execute(); err != nil {
		fmt.Printf("Couldn't run app: %s", err)
		os.Exit(1)
	}
}
