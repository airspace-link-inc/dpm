package dbm

import (
	"reflect"
	"testing"
)

// Test struct
type s struct {
	X          *float32 `db:"x" json:"foo"`
	Y          int      `db:"y" json:"bar" custom:"y-axis"`
	Z          string   `db:"z" json:"biz"`
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
		TestName     string
		ExpectedCols []string
		ExpectedVals []any
		OmitColumns  []string
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
		TestName     string
		ExpectedCols []string
		ExpectedVals []any
		Use          []string
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

// Test Map(), which when invoked transforms a struct into a map<tag, fieldValue>
func TestMap(test *testing.T) {

	test.Run("Should generate map from struct when Map() invoked", func(t *testing.T) {
		expected := [][]any{
			{"x", nil},
			{"y", 1},
			{"z", "foo"},
		}

		// Only tags specifed in Use() should be returned
		m := Params(s{Y: 1, Z: "foo"}).Map()

		for _, pair := range expected {
			expectedCol := pair[0].(string)
			expectedVal := pair[1]

			// Search for col name in the strct map
			actual, exists := m[expectedCol]

			if !exists {
				t.Errorf("Expected key value pair [%v, %v] not present in output.", expectedCol, expectedVal)
			}

			// Since the result is of type interface{}, null values need to be treated a bit differently
			// in order to avoid a panic
			if expectedVal == nil {

				if !reflect.ValueOf(actual).IsNil() {
					t.Errorf("Expected value %v not equal to expected %v", actual, expectedVal)
				}

				continue
			}

			if actual != expectedVal {
				t.Errorf("Expected value %v not equal to expected %v", actual, expectedVal)
			}
		}
	})

}

// Test Tag() modifies the name of the columns returned to match the provided string.
// But does not modify the values returned
func TestTag(test *testing.T) {

	testCases := []struct {
		TestName     string
		ExpectedCols []string
		ExpectedVals []any
		TagName      string
	}{
		{
			"Should use json tag value as col names",
			[]string{"foo", "bar", "biz"},
			[]any{nil, 1, "foo"},
			"json",
		},
		{
			"Should yield no cols nor vals when provided invalid tag name",
			[]string{},
			[]any{},
			"M",
		},
		{
			"should yield no cols nor vals when provided empty string as tag name",
			[]string{},
			[]any{},
			"",
		},
		{
			"should support custom tag values set by the user",
			[]string{"y-axis"},
			[]any{1},
			"custom",
		},
	}

	for _, tc := range testCases {
		test.Run(tc.TestName, func(t *testing.T) {

			// Col names should be associated with the value provided to Tag()
			cols, vals := Params(s{Y: 1, Z: "foo"}).Tag(tc.TagName).FlatVals()

			// Validate output matches expected
			validateColumns(test, tc.ExpectedCols, cols)
			validateValues(test, tc.ExpectedVals, vals)
		})
	}
}
