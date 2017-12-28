package main

import (
	"fmt"

	"github.com/juju/errors"
)

func main() {
	fmt.Println(errors.ErrorStack(Foo()))
}

func Foo() error {
	if err := Bar(); err != nil {
		return errors.Annotate(err, "foo")
	}

	return nil
}

func Bar() error {
	err := errors.New("bar error")
	return errors.Annotate(err, "bar")
}
