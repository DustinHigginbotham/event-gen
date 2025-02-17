package main

import (
	"fmt"
	"os"

	eventgen "github.com/DustinHigginbotham/event-gen/internal"
)

func main() {
	if err := eventgen.Generate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
