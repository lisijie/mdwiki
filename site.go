package main

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type site struct {
	title        string   //网站标题
	keywords     string   //关键字
	description  string   //描述信息
	url          string   //网址
	staticDir    []string //静态文件目录
	docDir       string   //md文档目录
	categoryList map[string]*category
}

func (s *site) build() {
	basePath := filepath.ToSlash(workPath + "/" + s.docDir)
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		path = filepath.ToSlash(path)
		if path == basePath {
			return nil
		}
		path = path[len(s.docDir):]

		if info.IsDir() {
			cat := newCategory(info.Name(), path)
			s.categoryList[path] = cat
			log.Println("初始化目录: ", path)
			return nil
		}

		if filepath.Ext(path) == ".md" {
			catkey := filepath.ToSlash(filepath.Dir(path))
			p, err := newPost(basePath, path)
			if err != nil {
				return err
			}
			if cat, ok := s.categoryList[catkey]; ok {
				cat.addPost(p)
				log.Println("初始化文章: ", path)
			}
		}

		return nil
	})

	checkError(err)
}

type category struct {
	path  string
	name  string
	posts map[string]*post
}

func (c *category) addPost(p *post) {
	c.posts[p.path] = p
}

type post struct {
	path    string
	title   string
	time    time.Time
	url     string
	content string
}

func newCategory(name string, path string) *category {
	return &category{name: name, path: path, posts: make(map[string]*post)}
}

func newSite() *site {
	return &site{
		categoryList: make(map[string]*category),
	}
}

func newPost(base string, path string) (*post, error) {
	ret, err := ioutil.ReadFile(base + path)
	if err != nil {
		return nil, err
	}

	p := &post{
		path:    path,
		title:   filepath.Base(path)[0 : len(filepath.Base(path))-3],
		content: string(blackfriday.MarkdownCommon(ret)),
	}

	return p, nil
}
