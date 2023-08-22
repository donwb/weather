package main

import "fmt"

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println("--------------------")
		fmt.Println("PANIC: ", msg)
		panic(err)
	}
}

func celsiusToFahrenheit(celsius float64) int {
	calc := (celsius * 9 / 5) + 32

	return int(calc)

}
