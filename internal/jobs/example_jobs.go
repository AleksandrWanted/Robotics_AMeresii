package jobs

import (
	"log"
	"math/rand"
)

func ExampleJob() {
	num := rand.Intn(100)

	switch {
	case num%2 == 0:
		log.Print(num)
	default:
		panic("some panic")
	}
}
