package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/alyosha/gungan"
)

func main() {
	jarJar, err := gungan.NewJarJar(true)
	if err != nil {
		log.Fatal("failed to create new jar jar", err)
	}

	input := flag.String("text", "", "text to translate")
	flag.Parse()

	if *input == "" {
		log.Fatal("must specify input via the -text flag")
	}

	fmt.Println(jarJar.Spake(*input))
}
