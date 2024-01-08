package utils

import "log"

func ErrorPanicPrinter(err error, shouldPanic bool) {
	if err != nil {
		log.Print(err)

		if shouldPanic {
			panic(err)
		}
	}
}
