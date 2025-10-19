package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Student struct {
	Name   string `validate:"required"`
	Avatar string `validate:"omitempty,url"`
	Grade  string `validate:"oneof=A B C D F"`
}

func main() {
	validator := validator.New()

	student := Student{}

	err := validator.Struct(student)
	fmt.Println("1.", err)

	student = Student{
		Name:  "Peter",
		Grade: "E",
	}

	err = validator.Struct(student)
	fmt.Println("2.", err)

	student = Student{
		Name:   "Peter",
		Avatar: "picture",
		Grade:  "A",
	}

	err = validator.Struct(student)
	fmt.Println("3.", err)

	student = Student{
		Name:   "Peter",
		Avatar: "https://your.image.com/avatar",
		Grade:  "A",
	}

	err = validator.Struct(student)
	fmt.Println("4.", err)
}
