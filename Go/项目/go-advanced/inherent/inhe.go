package main

type A interface {
	f()
}
type B struct {
	A
}
type C struct {
	a A
}

var b B = B{p}
var c C = C{p}

type impl int

var p impl

func (i impl) f() {
}

func main() {
	b.f()
	c.a.f()
}
