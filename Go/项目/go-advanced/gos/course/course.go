package main

func f1() func() int {
	x := 0
	return func() int {
		x += 1
		return x
	}
}

func main() {
	f := f1()
	f()
}
