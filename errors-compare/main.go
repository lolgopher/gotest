package main

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	UsecaseErrIndureUDPTransportFail = errors.New("indure udp transport fail")
	DtoErrIndureUDPTransportFail     = errors.New("indure udp transport fail")
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
	myErr := &MyError{}
	if errors.Is(err, &SomeOtherError{}) {
		fmt.Println("SomeOtherError occurred")
	} else if errors.As(err, &myErr) {
		fmt.Println("MyError occurred")
	} else {
		fmt.Println("Unknown error occurred")
	}
	fmt.Println()

	err = UsecaseErrIndureUDPTransportFail
	fmt.Println(err)

	if errors.Is(err, UsecaseErrIndureUDPTransportFail) {
		fmt.Println("UsecaseErrIndureUDPTransportFail errors.Is match")
	}
	if errors.As(err, &UsecaseErrIndureUDPTransportFail) {
		fmt.Println("UsecaseErrIndureUDPTransportFail errors.As match")
	}
	fmt.Println()

	err = DtoErrIndureUDPTransportFail
	fmt.Println(err)

	if errors.Is(err, UsecaseErrIndureUDPTransportFail) {
		fmt.Println("DtoErrIndureUDPTransportFail errors.Is match")
	}
	if errors.As(err, &UsecaseErrIndureUDPTransportFail) {
		fmt.Println("DtoErrIndureUDPTransportFail errors.As match")
	}
	fmt.Println()

	err = errors.Wrap(DtoErrIndureUDPTransportFail, "fail to send rtp packet")
	fmt.Println(err)

	if errors.Is(err, UsecaseErrIndureUDPTransportFail) {
		fmt.Println("Wrapped DtoErrIndureUDPTransportFail errors.Is match")
	}
	if errors.As(err, &UsecaseErrIndureUDPTransportFail) {
		fmt.Println("Wrapped DtoErrIndureUDPTransportFail errors.As match")
	}
	fmt.Println()

	err = errors.Wrap(DtoErrIndureUDPTransportFail, "write: operation not permitted")
	err = errors.Wrap(err, "fail to send rtp packet")
	fmt.Println(err)

	if errors.Is(err, UsecaseErrIndureUDPTransportFail) {
		fmt.Println("Two Wrapped DtoErrIndureUDPTransportFail errors.Is match")
	}
	if errors.As(err, &UsecaseErrIndureUDPTransportFail) {
		fmt.Println("Two Wrapped DtoErrIndureUDPTransportFail errors.As match")
	}
	fmt.Println()
}

func someFunction() error {
	return &MyError{}
}
