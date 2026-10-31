// Package member holds a library member's card.
package member

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Member is a person allowed to borrow books.
type Member struct {
	CardNo string `json:"card_no"`
	Name   string `json:"name"`
}

// New checks that the card number and name are present.
func New(cardNo, name string) (Member, error) {
	if cardNo == "" || name == "" {
		return Member{}, errors.New("card and name needed")
	}
	return Member{CardNo: cardNo, Name: name}, nil
}

// Display prints the member card.
func (m Member) Display() {
	fmt.Printf("Card %s: %s\n", m.CardNo, m.Name)
}

// Save writes the member as JSON.
func (m Member) Save() error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	name := "member_" + strings.ToLower(m.CardNo) + ".json"
	return os.WriteFile(name, data, 0644)
}
