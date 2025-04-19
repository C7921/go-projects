package main

import (
	"errors" // Used for creating custom error messages
	"fmt"
	"math" // Required for floating-point Modulo operation
	"os"
	"strconv"
)

const (
	expectedArgCount = 4 // Program name + 3 user arguments
	num1Index        = 1
	operatorIndex    = 2
	num2Index        = 3
)

func main() {
	// Argument Validation
	if len(os.Args) != expectedArgCount {
		fmt.Fprintf(os.Stderr, "Usage: %s <number1> <operator> <number2>\n", os.Args[0]) // os.Args[0] is the program name
		fmt.Fprintln(os.Stderr, "Example: go run calculator_v2.go 5.5 + 3.1")
		fmt.Fprintln(os.Stderr, "Supported operators: +, -, x, /,%")
		os.Exit(1)
	}

	// Input Parsing
	num1, err := parseArgAsFloat(os.Args[num1Index], "first number")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	operator := os.Args[operatorIndex]

	num2, err := parseArgAsFloat(os.Args[num2Index], "second number")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Calculation
	result, err := calculate(num1, operator, num2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output
	fmt.Printf("%g %s %g = %g\n", num1, operator, num2, result)
}

// parseArgAsFloat converts a string argument to a float64.
// Returns the parsed float and an error if the conversion fails.
func parseArgAsFloat(arg string, argName string) (float64, error) {
	val, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s provided ('%s'): %w", argName, arg, err)
	}
	return val, nil
}

// calculate performs the arithmetic operation based on the operator.
// Returns the result and an error if the operation is invalid.
// Supported operators: +, -, x, /, % (modulo)
func calculate(num1 float64, operator string, num2 float64) (float64, error) {
	var result float64

	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "x":
		result = num1 * num2
	case "/":
		// Check for division by zero *before* performing the operation.
		if num2 == 0 {
			return 0, errors.New("division by zero")
		}
		result = num1 / num2
	case "%":
		if num2 == 0 {
			return 0, errors.New("division by zero (for modulo)")
		}
		result = math.Mod(num1, num2)
	default:
		return 0, fmt.Errorf("unsupported operator '%s'", operator)
	}

	// Return the calculated result and nil error to indicate success.
	return result, nil
}
