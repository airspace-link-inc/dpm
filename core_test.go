package dbm

import (
	"testing"
)

// Test struct
type s struct {
	X *float32  `db:"x"`
	Y int `db:"y"`
	Z string `db:"z"`
	Unexported float32
}


// Test retrieval of "column names", the tag associated with a struct field
func TestCols(t *testing.T) {
	cols := Params(s{Y: 1, Z: "foo"}).Cols()

	validateColumns(
		t, 
		[]string{"x", "y", "z"}, 
		cols,
	)
}

// Test omission based on input col name
func TestOmit(test *testing.T) {
	testCases := []struct {
		TestName string
		ExpectedCols []string
		ExpectedVals []any
		OmitColumns []string
	}{
		{
			"Successfully omit valid tag",
			[]string{"y", "z"},
			[]any{1, "foo"},
			[]string{"x"},
		},
		{
			"If invalid tag is given nothing should be omitted",
			[]string{"x", "y", "z"},
			[]any{nil, 1, "foo"},
			[]string{"X"}, // Uppercase when tag is lowercase
		},
		{
			"Omit all fields in struct should return empty slices",
			[]string{},
			[]any{},
			[]string{"x", "y", "z"},
		},
		{
			"Nothing provided to omission doesn't error nor modify result",
			[]string{"x", "y", "z"},
			[]any{nil, 1, "foo"},
			[]string{},
		},
	}

	for _, tc := range testCases {
		test.Run(tc.TestName, func(t *testing.T) {

			// Omit in the response the struct field with the (?) tag
			cols, vals := 
				Params(s{Y: 1, Z: "foo"}).
					Omit(tc.OmitColumns...).
					FlatVals()
			

			// Validate output matches expected 
			validateColumns(test, tc.ExpectedCols, cols)
			validateValues(test, tc.ExpectedVals, vals)
		})
	}	
}

// Test Use() restricts output to only fields specified in its input params
func TestUse(test *testing.T) {
	testCases := []struct {
		TestName string
		ExpectedCols []string
		ExpectedVals []any
		Use []string
	}{
		{
			"Should only include values of tags passed to the Use() modifier",
			[]string{"y", "z"},
			[]any{1, "foo"},
			[]string{"y", "z"},
		},
		{
			"Should yield no cols nor vals when provided invalid tag",
			[]string{},
			[]any{},
			[]string{"M"},
		},
		{
			"should yield no cols nor vals when provided empty slice",
			[]string{},
			[]any{},
			[]string{},
		},
	}

	for _, tc := range testCases {
		test.Run(tc.TestName, func(t *testing.T) {

			// Only tags specifed in Use() should be returned
			cols, vals := Params(s{Y: 1, Z: "foo"}).Use(tc.Use...).FlatVals()
			
			// Validate output matches expected 
			validateColumns(test, tc.ExpectedCols, cols)
			validateValues(test, tc.ExpectedVals, vals)
		})
	}
}
