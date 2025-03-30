package cache

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	Set("a", "A1")
	Set("b", "B1")
	data := Get("a")
	fmt.Println(data)

	data = Get("a")
	fmt.Println(data)

	Set("a", "A2")
	data = Get("a")
	fmt.Println(data)
	ls := GetAll()
	fmt.Println(ls)
}
