package main

import (
	"net/http"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"strings"
	"crucians-blog/logger"
	"regexp"
)

func main() {
	/*workDir := config.LoadString(constant.CMD_KEY_WD)
	port := config.LoadInt(constant.CMD_KEY_PORT)
	if len(workDir) == 0 {
		workDir = "."
	}
	if port == 0 {
		port = 7240
	}*/
	//http.Handle("/static/", http.FileServer(http.Dir(workDir)))
	//http.HandleFunc("/", HttpRequest)
	//http.ListenAndServe(":"+strconv.Itoa(port), nil)
	md := `sjkdj
---
layout:     post
title:      "【TODO】2019下半年"
subtitle:   "2019下半年TODO,希望不会只是计划"
date:       2019-06-18
author:     "wangnem"
header-img: "img/bg-todo.png"
tags:
    - TODO
    - 计划
    - 2019
---
---
dhdjh
---
放在这里，想到什么都放到里面；希望到年底的时候，列表中的大部分都已经完成；这算是对自己的鞭策吧。
希望这篇文章不要变成今日对自己的安慰，将来对自己的愧疚！

### LIST

- [] 看完《机器学习实战》，至少写一篇文章`
	mdBytes,_:=ioutil.ReadAll(strings.NewReader(md))
	out := blackfriday.Run(mdBytes,blackfriday.WithExtensions(blackfriday.FencedCode))
	logger.Info("%s",string(out))
	reg :=regexp.MustCompile(`-{3}[\s|\S]*?-{3}`)
	find := reg.FindString(md)
	b,e := regexp.MatchString("^--[\\s|\\S]*--\\s",md)
	if e != nil {
		logger.Panic(e.Error())
	}

	logger.Info("%t",b)
	logger.Info("%q",find)

}

func HttpRequest(rw http.ResponseWriter, r *http.Request) {

}
