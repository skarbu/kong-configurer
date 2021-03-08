package logger

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s \n", msg, err)
	}
}

func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s : %s \n", msg, err)
	}
}

func FailOnErrors(errs []error, msg string) {
	if len(errs) != 0 {
		for _, err := range errs {
			log.Print(err)
		}
		log.Fatal(msg)
	}
}
