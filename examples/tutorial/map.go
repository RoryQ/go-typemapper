//go:build typemapper
// +build typemapper

package main

func MapFooToBar(src Foo) Bar {
	dst := Bar{}
	typemapper.CreateMap(src, dst)
	return dst
}
