package main

import (
	"fmt"

	"github.com/KendallBritton/CloudNativeApp/Labs/Lab1/myadder"
	"github.com/KendallBritton/CloudNativeApp/Labs/Lab1/textproc"
)

func main() {

	fmt.Println(myadder.Add(5, 6))

	filename := "passage"

	textproc.TopWords(filename, 3)

}
