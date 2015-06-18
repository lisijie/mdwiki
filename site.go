package main

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type site struct {
	title       string   //网站标题
	keywords    string   //关键字
	description string   //描述信息
	url         string   //网址
	staticDir   []string //静态文件目录
	docDir      string   //md文档目录
	workDir     string   //工作目录，web根目录
}

type post struct {
	title   string
	time    time.Time
	url     string
	content string
}

type dir struct {
	name  string
	url   string
	posts map[string]*post
}

func make_post(mdfile string) (*post, error) {
	res, err := ioutil.ReadFile(mdfile)
	if err != nil {
		return nil, err
	}

	p := &post{
		title:   mdfile,
		content: string(blackfriday.MarkdownCommon(res)),
	}

	return p, nil
}

func build(path string) {
	fi, err := os.Open(path)
	checkError(err)

	dirs, err := fi.Readdir(0)
	checkError(err)

	for _, v := range dirs {
		if !v.IsDir() {
			continue
		}
		d := new(dir)
		d.name = v.Name()
		d.posts = make(map[string]*post)

		log.Println(v.Name())
	}
}
