package main

import (
	"log"

	"github.com/jare-abc/XrayR-dev/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
