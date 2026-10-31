package main

import "fmt"

type MyError struct{ msg string }

func (e *MyError) Error() string { return e.msg }

func check(ok bool) error {
	var e *MyError
	if !ok {
		e = &MyError{"check failed"}
	}
	return e // wrong: a nil *MyError is not a nil error
}

func checkFixed(ok bool) error {
	if !ok {
		return &MyError{"check failed"}
	}
	return nil
}

func main() {
	fmt.Println(check(true) != nil)      // true
	fmt.Println(checkFixed(true) != nil) // false
}
