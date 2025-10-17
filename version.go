package main

import "fmt"

var version = "dev"

func getVersion() string {
	return version
}

func printVersion() {
	fmt.Printf("cogmit version %s\n", version)
}
