package main

import (
	"fmt"

	"gitlab.com/ovsinc/errors"
)

func main() {
	fmt.Printf("%v\n", errors.New("simple error"))
}
