# DBParam

Uses reflection to make it easier to interact with a database

```go
type S struct {
  Y int `db:"y"`
  Z int `db:"z"`
  Geom orb.Geometry `db:"geom"`
  Unexported float32
}

p := db.Params(S{Y: 10, Z: 15, Geom: orb.Point{1,2}, Unexported: 0.5})

// Generating maps for NamedQuery
p.Map()           // map[string]any{"y": 10, "z": 15, "geom": `{"type":"Point",...}`}
p.Use("y").Map()  // map[string]any{"y": 10}
p.Omit("y").Map() // map[string]any{"z": 15, "geom": `{"type":"Point",...}`}

// Getting columns (you can filter some out if that's required)
p.Omit("geom").Cols() // []string{"y", "z"}

// Getting values
p.Use("y", "z").Vals() // []any{10, 15}

// Zipping cols and vals
p.Omit("geom").FlatVals() // ([]string{"y", "z"}, []any{10, 15})

// Custom logic for mapping things (note that this is absolutely wrong for the Geom column)
p.Mapper(func (x any) any { return 123 }).Vals() // []any{123, 123, 123}
```
