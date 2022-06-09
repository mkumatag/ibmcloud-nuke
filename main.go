package main

import (
	"github.com/mkumatag/ibmcloud-nuke/cmd"
	"log"
)

func main() {
	log.Println("Hello, Lets's Nuke it!")
	if err := cmd.NewRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
