package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	f, _ := os.CreateTemp(".", "example*.txt")
	fmt.Println(filepath.Base(f.Name()))
	f.WriteString("kek")

	f.Close()

}
