package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	//create file
	f, _ := os.CreateTemp(".", "example*.txt")
	fmt.Println(filepath.Base(f.Name()))
	f.WriteString("kek")
	fmt.Println(filepath.Abs(f.Name()))
	absName, _ := filepath.Abs(f.Name())

	f2, err := os.OpenFile(absName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	} //doesnt depend on the f
	f.Close()
	n, err := f2.WriteString("appended")
	fmt.Println(n, err)
	//f.Close()
	f2.Close() //have to close this one as well along side with f

	//delete file without closing first *file(not possible

	err = os.Remove(absName)
	fmt.Println(err)

}
