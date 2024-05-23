package auth

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func init() {
	os.Chdir(path.Dir("../../../")) // 设置运行根目录
}

func TestParseJwtToken(t *testing.T) {
	token, err := GenJwtToken(2, "AAA")
	if err != nil {
		t.Fatal(err)
	}
	//token := ctx.Request.Header.Get("Authorization")

	u, err := ParseJwtToken(token)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(u)
}
