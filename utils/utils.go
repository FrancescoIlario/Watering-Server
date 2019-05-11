package utils

import "log"

func PanicIf(err error) {
	if err != nil {
		log.Panicf("Panic! Error occurred: %s\n", err.Error())
	}
}
