package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	//lint:ignore U1000 Ignore unused variable
	release = "UNKNOWN"
	//lint:ignore U1000 Ignore unused variable
	buildDate = "UNKNOWN"
	//lint:ignore U1000 Ignore unused variable
	gitHash = "UNKNOWN"
)

//lint:ignore U1000 Ignore unused variable
func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
