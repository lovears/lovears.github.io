package main

import (
	"os"
	"strings"
	"net/http"
)

var env = make(map[string]string)

func init() {
	args := os.Args
	var lastVal string
	for index, arg := range args {
		if index == 0 {
			env["boot"] = arg
		} else if strings.Contains(arg, "=") {
			envMapping := strings.Split(arg, "=")
			env[strings.TrimSpace(strings.Replace(envMapping[0], "--", "", -1))] = strings.TrimSpace(envMapping[1])
		} else {
			if lastVal != "" {
				env[lastVal] = arg
				lastVal = ""
			} else {
				lastVal = arg
			}
		}
	}
}

func main() {

	wd := env["workdir"]
	if wd == "" {
		wd = "."
	}

	http.Handle("/", http.FileServer(http.Dir(wd)))

	http.ListenAndServe(":8001",nil)
}
