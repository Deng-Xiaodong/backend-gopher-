package reflect

import "testing"

func TestPrintAny(t *testing.T) {
	println(Any(10))
	println(Any(10.23))
	println(Any("love"))
	println(Any(true))
	println(Any([]int{}))
	println(Any(map[int]int{}))

}
func TestEmpty(t *testing.T) {
	Empty()
}
