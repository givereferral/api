// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

// +build ignore

package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	for _, typ := range []struct {
		Type, Name, Atomic string
		Number             bool
	}{
		{"int32", "Int32", "Int32", true},
		{"int64", "Int64", "Int64", true},
		{"uint32", "Uint32", "Uint32", true},
		{"uint64", "Uint64", "Uint64", true},
		{"float32", "Float32", "Float32", true},
		{"float64", "Float64", "Float64", true},
		{"bool", "Bool", "Bool", false},
		{"string", "String", "String", false},
	} {
		f, err := os.Create(typ.Type + ".go")
		if err != nil {
			log.Fatal(err)
		}

		if err = tmpl.Execute(f, typ); err != nil {
			log.Fatal(err)
		}

		f.Close()
	}
}

var tmpl = template.Must(template.New("impl").Parse(`// Code generated by go run generate-int.go.

// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package maps

import (
	"sync"

	"github.com/tmthrgd/atomics"
)

// {{.Name}} provides a map of atomic {{.Type}}s.
type {{.Name}} struct {
	m sync.Map // map[interface{}]*atomics.{{.Atomic}}
}

// Retrieve returns the atomics.{{.Atomic}} associated with
// the given key or nil if it does not exist in the map.
func (m *{{.Name}}) Retrieve(key interface{}) *atomics.{{.Atomic}} {
	v, ok := m.m.Load(key)
	if !ok {
		return nil
	}

	return v.(*atomics.{{.Atomic}})
}

// Insert inserts the atomics.{{.Atomic}} into the map for
// the given key.
func (m *{{.Name}}) Insert(key interface{}, val *atomics.{{.Atomic}}) {
	m.m.Store(key, val)
}

// Value returns the atomics.{{.Atomic}} associated with the
// given key or atomically inserts a new atomics.{{.Atomic}}
// into the map if an entry did not exist in the map
// for the given key.
func (m *{{.Name}}) Value(key interface{}) *atomics.{{.Atomic}} {
	v, ok := m.m.Load(key)
	if !ok {
		v, _ = m.m.LoadOrStore(key, new(atomics.{{.Atomic}}))
	}

	return v.(*atomics.{{.Atomic}})
}

// Delete removes an atomics.{{.Atomic}} from the map.
func (m *{{.Name}}) Delete(key interface{}) {
	m.m.Delete(key)
}

// Range calls f for each entry in the map. If f
// returns false Range stops iterating over the map.
func (m *{{.Name}}) Range(f func(key interface{}, val *atomics.{{.Atomic}}) bool) {
	m.m.Range(func(key, val interface{}) bool {
		return f(key, val.(*atomics.{{.Atomic}}))
	})
}

// Load is a wrapper for Value(key).Load().
func (m *{{.Name}}) Load(key interface{}) (val {{.Type}}) {
	return m.Value(key).Load()
}

// Store is a wrapper for Value(key).Store(val).
func (m *{{.Name}}) Store(key interface{}, val {{.Type}}) {
	m.Value(key).Store(val)
}

// Swap is a wrapper for Value(key).Swap(new).
func (m *{{.Name}}) Swap(key interface{}, new {{.Type}}) (old {{.Type}}) {
	return m.Value(key).Swap(new)
}

{{- if not (eq .Atomic "String")}}

// CompareAndSwap is a wrapper for
// Value(key).CompareAndSwap(old, new).
func (m *{{.Name}}) CompareAndSwap(key interface{}, old, new {{.Type}}) (swapped bool) {
	return m.Value(key).CompareAndSwap(old, new)
}

{{- end}}

{{- if .Number}}

// Add is a wrapper for Value(key).Add(delta).
func (m *{{.Name}}) Add(key interface{}, delta {{.Type}}) (new {{.Type}}) {
	return m.Value(key).Add(delta)
}

// Increment is a wrapper for Value(key).Increment().
func (m *{{.Name}}) Increment(key interface{}) (new {{.Type}}) {
	return m.Value(key).Increment()
}

// Subtract is a wrapper for Value(key).Subtract(delta).
func (m *{{.Name}}) Subtract(key interface{}, delta {{.Type}}) (new {{.Type}}) {
	return m.Value(key).Subtract(delta)
}

// Decrement is a wrapper for Value(key).Decrement().
func (m *{{.Name}}) Decrement(key interface{}) (new {{.Type}}) {
	return m.Value(key).Decrement()
}

{{- end}}

{{- if eq .Atomic "Bool"}}

// Set is a wrapper for Value(key).Set().
func (m *{{.Name}}) Set(key interface{}) (old {{.Type}}) {
	return m.Value(key).Set()
}

{{- end}}

// Reset is a wrapper for Value(key).Reset().
func (m *{{.Name}}) Reset(key interface{}) (old {{.Type}}) {
	return m.Value(key).Reset()
}
`))
