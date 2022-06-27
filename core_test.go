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