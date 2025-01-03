package conv

import (
	"fmt"
	"testing"
)

func ConvertAA[To, From int | int8 | int16 | int32](from From) To {
	to := To(from)
	if From(to) != from {
		panic("conversion out of range")
	}
	return to
}

func TestNewStr_Slice(t *testing.T) {
	list := []int{9}

	err := NewStr().Slice("1,2,3", &list)
	t.Log(err)
	fmt.Printf("结果 %#v\n", list)

	a := ConvertAA[int](5)
	fmt.Println(a)
}
