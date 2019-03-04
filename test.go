package main

import (
	"fmt"

	"./lib"
)

func main() {
	urls, _ := lib.GetSankakuPictures("1", true)
	fmt.Println(urls)
}
