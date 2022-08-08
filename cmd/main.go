package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
)

func main() {

	createAppFlag := flag.Bool("create-app", false, "Create Kubefirst Go Application")
	languageFlag := flag.String("language", "", "Set application programming language to be created")
	flag.Parse()

	if !*createAppFlag {
		log.Warn().Msg("create app not enable, exiting...")
		return
	}

	if *languageFlag != "go" {
		log.Warn().Msg("Go is the best, no other language is necessary!")
		return
	}

	fmt.Println("hi")
}
