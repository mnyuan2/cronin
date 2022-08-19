package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestElasticConf(t *testing.T) {
	fmt.Printf("%#v\n", ElasticConf())
}

func TestDbConf(t *testing.T) {
	err := os.Chdir("/")
	fmt.Printf("%s\n", err)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Printf("err:%s,	dir:%#v\n", err, dir)

	file, err := exec.LookPath(os.Args[0])
	fmt.Printf("err:%s,	dir:%#v\n", err, file)

	path, err := os.Executable()
	fmt.Printf("err:%s,	dir:%#v\n", err, path)

	c := MainConf()
	fmt.Printf("%#v\n", c)
}

func TestHttpConfs_GetConf(t *testing.T) {
	os.Chdir("/")
	conf := Http()

	fmt.Printf("%#v\n", conf)
}