package main

import "fmt"

type Men interface {
	SayHi()
}

type Human struct {
	name  string
	age   int
	phone string
}

func (h Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

type Student struct {
	Human
	school string
	loan   float32
}

//Human实现Sing方法
func (h Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

func main() {
	nike := Student{
		Human{
			"nike",
			20,
			"13012345678",
		},
		"SBC",
		1.0,
	}

	nike.SayHi()

	var men Men
	men = nike
	men.SayHi()
}
