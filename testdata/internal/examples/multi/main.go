package main

import (
	"fmt"

	"gitlab.com/ovsinc/errors"
)

func main() {
	fmt.Printf("%v\n",
		errors.Append(
			errors.New("one error"),
			errors.New("two error"),
			errors.New("three error"),
		),
	)
}
