package dbm

import (
	"github.com/goccy/go-reflect"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
)

// DBParam uses reflection to make it easier to interact with a database:
//
//   type S struct {
//     Y int `db:"y"`
//     Z int `db:"z"`
//     Geom orb.Geometry `db:"geom"`
//     Unexported float32
//   }
//
//   p := db.Params(S{Y: 10, Z: 15, Geom: orb.Point{1,2}, Unexported: 0.5})
//
//   // Generating maps for NamedQuery
//   p.Map()           // map[string]any{"y": 10, "z": 15, "geom": `{"type":"Point",...}`}
//   p.Use("y").Map()  // map[string]any{"y": 10}
//   p.Omit("y").Map() // map[string]any{"z": 15, "geom": `{"type":"Point",...}`}
//
//   // Getting columns (you can filter some out if that's required)
//   p.Omit("geom").Cols() // []string{"y", "z"}
//
//   // Getting values
//   p.Use("y", "z").Vals() // []any{10, 15}
//
//   // Zipping cols and vals
//   p.Omit("geom").FlatVals() // ([]string{"y", "z"}, []any{10, 15})
//
//   // Custom logic for mapping things (note that this is absolutely wrong for the Geom column)
//   p.Mapper(func (x any) any { return 123 }).Vals() // []any{123, 123, 123}
type DBParam struct {
	inner       any
	colFilters  []func(s string) bool
	valFilters  []func(s any) bool
	filterZero  bool
	mapper      func(x any) any
	extraCols   []string
	extraValues []any
	tag         string
}

// Default mapper is the mapper to use for turning a struct's value into
// something close to driver.Value
var DefaultMapper = func(x any) any {
	return x
}

// Params creates a new DBParam pointer
func Params(x any) *DBParam {
	d := DBParam{inner: x, colFilters: []func(string) bool{}, mapper: defaultMapper, tag: "db", valFilters: []func(s any) bool{}}
	return &d
}

// Mapper alters the mapping function used by this params struct. This allows you to change the
// logic used for turning the struct into something insertable, if you require that control
func (d *DBParam) Mapper(fn func(any) any) *DBParam {
	d.mapper = fn
	return d
}

func (d *DBParam) Tag(s string) *DBParam {
	d.tag = s
	return d
}

// Use will only include the struct fields specified when exported
func (d *DBParam) Use(fields ...string) *DBParam {
	d.colFilters = append(d.colFilters, func(s string) bool {
		for _, v := range fields {
			if v == s {
				return true
			}
		}
		return false
	})

	return d
}

// Omit will exclude the following fields when exporting
func (d *DBParam) Omit(fields ...string) *DBParam {
	d.colFilters = append(d.colFilters, func(s string) bool {
		for _, v := range fields {
			if v == s {
				return false
			}
		}
		return true
	})
	return d
}

// FilterZero adds a filter function that filters out all zero values
func (d *DBParam) FilterZero() *DBParam {
	d.filterZero = true
	return d
}

func (d *DBParam) Filter(filters ...func(s any) bool) *DBParam {
	d.valFilters = append(d.valFilters, filters...)
	return d
}

func (d *DBParam) AddKV(col string, value any) *DBParam {
	d.extraCols = append(d.extraCols, col)
	d.extraValues = append(d.extraValues, value)
	return d
}

// Cols will return the columns that satisfy all the filters specified
func (d *DBParam) Cols() []string {
	cols, _ := d.FlatVals()
	return cols
}

// Vals will return the values that satisfy all the filters specified
func (d *DBParam) Vals() []any {
	_, vals := d.FlatVals()
	return vals
}

// FlatVals essentially returns (Cols(), Vals())
func (d *DBParam) FlatVals() ([]string, []any) {
	t := d.getType()
	v := d.getVal()

	n := len(d.extraCols)
	tags := make([]string, 0, t.NumField()+n)
	vals := make([]any, 0, t.NumField()+n)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(d.tag)
		if tag == "" {
			continue
		}

		for _, filter := range d.colFilters {
			if !filter(tag) {
				continue
			}
		}

		reflectVal := v.Field(i)
		if d.filterZero && reflectVal.IsZero() {
			continue
		}

		value := reflectVal.Interface()
		for _, filter := range d.valFilters {
			if !filter(value) {
				continue
			}
		}

		tags = append(tags, tag)
		vals = append(vals, d.mapper(value))
	}

	for i := 0; i < n; i++ {
		tags = append(tags, d.extraCols[i])
		vals = append(vals, d.extraValues[i])
	}

	return tags, vals
}

// Map returns a map that points the struct tag to its respective value
func (d *DBParam) Map() map[string]any {
	tags, vals := d.FlatVals()
	n := len(tags)
	m := make(map[string]any, n)
	for i := 0; i < n; i++ {
		m[tags[i]] = vals[i]
	}

	return m
}

func (d *DBParam) getType() reflect.Type {
	t := reflect.TypeOf(d.inner)

	// If it's an interface or a pointer, unwrap it.
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		t = t.Elem()
	}

	return t
}

func (d *DBParam) getVal() reflect.Value {
	v := reflect.ValueOf(d.inner)

	// If it's an interface or a pointer, unwrap it.
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		v = v.Elem()
	}

	return v
}
