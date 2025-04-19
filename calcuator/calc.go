package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	// check arguments
	if len(os.Args) != 4 {
		fmt.Println("Usage: calc <num1> <operator> <num2>")
		fmt.Println("Supported operators: +, -, *, /, %")
		os.Exit(1)
	}

	// Parse first num
	num1, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Error parsing first number: %v\n", err)
		os.Exit(1)
	}

	// Get Operator
	operation := os.Args[2]

	// Parse second num
	num2, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Printf("Error parsing second number: %v\n", err)
		os.Exit(1)
	}

	// Perform operation
	var result float64
	switch operation {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			fmt.Println("Error: Division by zero")
			os.Exit(1)
		}
		result = num1 / num2
	case "%":
		if num2 == 0 {
			fmt.Println("Error: Division by zero")
			os.Exit(1)
		}
		result = float64(int(num1) % int(num2))
	default:
		fmt.Printf("Error: Unsupported operator %s\n", operation)
		fmt.Println("Supported operators: +, -, *, /, %")
		os.Exit(1)
	}

	// Print Result, rounding to 3 decimal places.
	fmt.Printf("Result: %5.f %s %5.f = %5.f\n", num1, operation, num2, result)

}
