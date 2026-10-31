// Command embedded prints notices that are compiled into
// the binary, so it runs without its source folder.
package main

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed welcome.txt
var welcome string

//go:embed notices/*.txt
var notices embed.FS

func main() {
	fmt.Print(welcome)
	names, err := fs.Glob(notices, "notices/*.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, name := range names {
		data, _ := notices.ReadFile(name)
		fmt.Printf("%s: %s", name, data)
	}
}
