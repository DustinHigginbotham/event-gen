package main

import (
	"fmt"
	"os"

	eventgen "github.com/DustinHigginbotham/event-gen/pkg"
)

func main() {
	if err := eventgen.Generate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
