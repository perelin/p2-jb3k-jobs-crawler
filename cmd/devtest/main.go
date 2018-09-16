package main

import (
	"fmt"
	"log"
	"os"
)

func printDir() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
func main() {
	printDir()
	//db.GetLastEntryDate()
}
