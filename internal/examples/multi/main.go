package main

import (
	"fmt"

	"github.com/ovsinc/errors"
)

func main() {
	fmt.Printf("%v\n",
		errors.Combine(
			errors.New("one error"),
			errors.New("two error"),
			errors.New("three error"),
		),
	)
}
