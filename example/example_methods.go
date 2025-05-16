package example

import (
	"log"
	"math/rand"
)

func SimpleExampleMethod() {
	num := rand.Intn(100)

	switch {
	case num%2 == 0:
		log.Print(num)
	default:
		panic("some panic")
	}
}
