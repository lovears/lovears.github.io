package config

import "strconv"

func load(key string) string {
	return ""
}

func LoadString(key string) string {
	return load(key)
}

func LoadInt(key string) int {
	s := load(key)
	r, e := strconv.Atoi(s)
	if e != nil {
		panic(e.Error())
	}
	return r
}
