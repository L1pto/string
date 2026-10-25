package main

import (
	"fmt"
	"myproject/toolkit"
)

func printResult(name string, src, result []int) {
	fmt.Printf("%-10s src=%v result=%v\n", name, src, result)
}

func main() {

	toolkit.RunNumbersDemo()
	toolkit.RunTextDemo("Привет, мир! Привет, Go!")
	toolkit.RunReportDemo()

}
