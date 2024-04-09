package utils

import (
	"log"
	"runtime"
)

// Function to log all errors. Set second parameter
// to true to terminate program on error.
func LogError(err error, isFatal ...bool) error {
	if err != nil {
		_, file, _, _ := runtime.Caller(1)
		log.Println("Error from", file)

		if len(isFatal) > 0 && isFatal[0] {
			log.Fatalln(err)
		}

		log.Println(err)
	}

	return err
}
