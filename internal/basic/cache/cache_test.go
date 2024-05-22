package cache

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	Add("a", "AAA")
	data := Get("a")
	fmt.Println(data)

	data = Get("a")
	fmt.Println(data)

	Add("a", "BBB")
	data = Get("a")
	fmt.Println(data)
}
