package dbm

import (
	"reflect"
	"testing"
)

// Validates expected and resultant columns match
// Fails the test with error message if cols do not match
func validateColumns(test *testing.T, expected, result []string) {
	numCols := len(result)
	numExpected := len(expected)

	// Verify # of cols returned matches expected
	if numCols != numExpected {
		test.Errorf("%v columns were returned, expected %v", numCols, numExpected)
	}

	if numCols == 0 {
		// If we expected no columns we can break here
		return
	}

	var expectedCol string

	// Verify values are returned in in order
	for i, col := range result {
		expectedCol = expected[i]

		if col != expectedCol {
			test.Errorf("expected column name %v, recieved %v", col,expectedCol)
		}
	}
}

// Validates expected and resultant values match
// Fails the test with error message if vals do not match
func validateValues(test *testing.T, expected, result []any) {
	num := len(result)
	numExpected := len(expected)

	if num != numExpected {
		test.Errorf("%v values were returned, expected %v", num, numExpected)
	}

	if num == 0 {
		// If we expected no values we can break here
		return
	}

	for i := 0; i < num; i++ {
		current := result[i]
		exp := expected[i]

		// Unwrap the interface{} to get the underlying value
		val := reflect.ValueOf(current)

		switch val.Kind() {
		case reflect.Invalid:
			// Result is invalid, this should never happen
			test.Fatalf("Invalid result processed in result %v", current)
		case reflect.Pointer:
			// First check if pointer is null
			// If so we need to reset current to null so it doesnt compare reflect.Value
			if val.IsNil() {
				current = nil
			}

			fallthrough

		default: 
			if exp != current {
				test.Fatalf("Expected value %v, does not match recieved %v", exp, current)
			}
		}
	}
}