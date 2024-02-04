package main

import (
	"fmt"
	"os"
)

func main() {
	command := CreateOwnableERC721Command()
	err := command.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
