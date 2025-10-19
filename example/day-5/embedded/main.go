package main

import (
	"encoding/json"
	"fmt"
)

type Author struct {
	Name  string
	Email string
}

func (author Author) Signature() string {
	return fmt.Sprintf("%s signed!", author.Name)
}

type Blog struct {
	Author
	ID      int
	Upvotes int32
}

func (blog Blog) Signature() string {
	return fmt.Sprintf("%s unsigned!", blog.Name)
}

func main() {
	blog := Blog{
		ID:      1,
		Upvotes: 99,
		Author: Author{
			Name:  "Peter",
			Email: "peter@email.com",
		},
	}

	fmt.Println("Author name:", blog.Name)
	fmt.Println("Author email:", blog.Email)

	fmt.Println("Signature:", blog.Signature())
	fmt.Println("Embedded Signature:", blog.Author.Signature())

	disp, _ := json.MarshalIndent(blog, "", "  ")
	fmt.Println(string(disp))
}
