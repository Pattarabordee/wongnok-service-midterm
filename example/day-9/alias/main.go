package main

import "fmt"

type Base struct {
	ID uint
}

func (base Base) Greeting() {
	fmt.Println("Hello")
}

type Teacher struct {
	Base
	Name string
}

type Student Base

func main() {
	t := Teacher{
		Name: "Peter",
	}

	t.Greeting()

	s := Student{ID: 1}

	fmt.Println(s.ID)
}
