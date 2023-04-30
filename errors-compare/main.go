package main

import (
	"errors"
	"fmt"
)

type MyError struct {
	message string
	cause   error
}

type SomeOtherError struct {
	message string
	cause   error
}

func (e *SomeOtherError) Error() string {
	return e.message
}

func (e *MyError) Error() string {
	return e.message
}

func (e *MyError) Unwrap() error {
	return e.cause
}

func main() {
	err := someFunction()
	if errors.Is(err, &SomeOtherError{}) {
		fmt.Println("SomeOtherError occurred")
	} else if errors.As(err, MyError{}) {
		fmt.Println("MyError occurred")
	} else {
		fmt.Println("Unknown error occurred")
	}
}

func someFunction() error {
	return &MyError{}
}
