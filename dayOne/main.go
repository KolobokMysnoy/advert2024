package main

import (
	"fmt"
	"os"
)

func main() {
	currentWorkDirectory, _ := os.Getwd()
	fmt.Println(currentWorkDirectory)
	f, err := os.Open(currentWorkDirectory + "/../dayOne/.env")
	if err != nil {
		fmt.Print(err)
		return
	}
	b := make([]byte, 100)
	f.Read(b)

	fmt.Print(b)
}
