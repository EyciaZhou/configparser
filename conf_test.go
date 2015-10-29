package config

import (
	"testing"
	"encoding/json"
	"fmt"
)

type Config struct {
	Addr string `default:"fffff"`
	Port string
	Path string `default:"path: /root/aaa/aaa/path-/pp.pp"`
	Num int
	noPub int
	List1 []string
	List2 []string
	List3 []int
	List4 []int
}

var (
	a = `{"list4": [11111], "list1": [], "list3": [1, 2, -100000, 1], "list2": ["aaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "adfadf"], "num": 11, "addr": "127.0.0.1", "port": ":80"}`
)

var (
	config Config
)

func TestLoadConfFromJson(t *testing.T) {
	var m map[string]interface{}

	//json.Unmarshal(([]byte)(a), &config)
	//config.noPub = 1431331231

	json.Unmarshal(([]byte)(a), &m);

	LoadConfFromJson(&config, m)

	fmt.Printf("%v", config)
}