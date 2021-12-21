package handlers

import "log"

func PanicOnError(err error) {
	if err == nil {
		return
	}

	log.Panicln(err)
}

func IgnoreError(_ error) {
	// Method created to encapsulate non-critical errors that are (currently) being ignored.
	// If this changes, I'll implement something here.
}