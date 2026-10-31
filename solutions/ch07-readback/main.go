package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch07/catalog/book"
)

func main() {
	data, err := os.ReadFile("cs_101.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var b book.Book
	if err := json.Unmarshal(data, &b); err != nil {
		fmt.Println(err)
		return
	}
	b.Display()
}
