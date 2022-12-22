package main

import (
	"fmt"

	"github.com/ovsinc/errors"
)

func main() {
	fmt.Printf("%v\n", errors.New("simple error"))
}
