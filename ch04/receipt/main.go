package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	receiptID := uuid.NewString()
	fmt.Println("Receipt", receiptID)
}
