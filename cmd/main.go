// cmd/main.go
package main

import (
	"fmt"
	"forum/internal/utils"
)

func main() {
	// Get the current working directory
	filepath, err := utils.GetFolderPath("..","testiiii")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(filepath)
}
