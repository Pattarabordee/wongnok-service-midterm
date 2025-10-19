package main

import "fmt"

type Animal interface {
	Say() string
}

type Dog struct {
	Name string
}

func (d Dog) Say() string {
	return fmt.Sprintf("%s say woof woof", d.Name)
}

type Cat struct {
	Name string
}

func (c Cat) Say() string {
	return fmt.Sprintf("%s say meow meow", c.Name)
}

func Say(a Animal) {
	fmt.Println(a.Say())
}

func main() {
	Say(Dog{Name: "Lucky"})
	Say(Cat{Name: "Kitty"})
}
