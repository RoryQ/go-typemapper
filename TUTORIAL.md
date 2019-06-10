# `typemapper` Tutorial

The initial files for this tutorial can be found in the [examples](examples/how-it-works) directory.

Given a package main with the following types:

**main.go**
```go
package main

import (
	"fmt"
)

type Foo struct {
    FieldOne string
    FieldTwo int
}

type Bar struct {
    FieldOne string
    Field2 int
}

func main() {
	src := Foo{
		FieldOne: "one",
		FieldTwo: 2,
	}
	dst := MapFooToBar(src)
	fmt.Printf("Src: %#v\nDst: %#v\n", src, dst)
}
```

To start using `typemapper`, start by creating a file to hold your declarations for mapping functions. This file must have the `typemapper` build flag and import the `typemapper` package:

**map.go**
```go
// +build typemapper

package main

import (
    "github.com/paultyng/go-typemapper"
)

func MapFooToBar(src Foo) Bar {
    dst := Bar{}
    typemapper.CreateMap(src, dst)
    return dst
}
```

Running `typemapper` will generate two files:

**map.generated.go**
```go
// Code generated by "typemapper "; DO NOT EDIT.

// +build !typemapper

package main

func MapFooToBar(src Foo) Bar {
	dst := Bar{}
	dst.FieldOne = src.FieldOne
	// no match for "Field2"
	return dst
}
```

**map.generated_test.go**
```go
// Code generated by "typemapper "; DO NOT EDIT.

// +build !typemapper

package main

import "testing"

func TestMapFooToBar(t *testing.T) {
	t.Fatal("no mapping for: [Field2]")
}
```

You can see in the first file, a comment is left for a destination field that was unable to be mapped. In the generated test file, a failing test was written which also indicates which fields were unmapped on the destination struct.

To handle this unmapped field, you have two different approaches you can take. You could ignore the field:

```go
func MapFooToBar(src Foo) Bar {
    dst := Bar{}
	typemapper.CreateMap(src, dst)
	typemapper.IgnoreFields(dst.Field2)
    return dst
}
```

Which would remove the comment and make the test pass. There is a field on the `src` type though that could also be manually mapped:

```go
func MapFooToBar(src Foo) Bar {
    dst := Bar{}
	typemapper.CreateMap(src, dst)
	typemapper.MapField(src.FieldTwo, dst.Field2)
    return dst
}
```

This results in a mapping function generated that looks like this:

```go
func MapFooToBar(src Foo) Bar {
	dst := Bar{}
	dst.FieldOne = src.FieldOne
	dst.Field2 = src.FieldTwo
	return dst
}
```

As well as making the unit test now passing.