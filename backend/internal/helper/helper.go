package helper

import "log"

func ErrorHandler(err error) {
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
