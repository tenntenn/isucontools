package isucontools

import (
	"log"
)

func Must(err error) {
	if err != nil {
		log.Panic(err)
	}
}
