package main

import "log"

// logError logs a non-fatal error. Unlike the old checkError, it never
// panics, so a transient network or parse failure no longer takes down
// the whole server.
func logError(err error, msg string) {
	if err != nil {
		log.Printf("ERROR: %s: %v", msg, err)
	}
}

func celsiusToFahrenheit(celsius float64) int {
	return int((celsius * 9 / 5) + 32)
}
