package main

import (
	"fmt"
	"labs/Lab1/myadder"
	"labs/Lab1/textproc"
)

func main() {

	fmt.Println(myadder.Add(5, 6))

	filename := "passage"

	textproc.TopWords(filename, 3)

}
