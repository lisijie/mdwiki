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
	Title        string   //网站标题
	Keywords     string   //关键字
	Description  string   //描述信息
	Url          string   //网址
	StaticDir    []string //静态文件目录
	DocDir       string   //md文档目录
	CategoryList map[string]*category
}

func (s *site) GetPost(path string) *post {
	ck := filepath.ToSlash(filepath.Dir(path))
	if cat, ok := s.CategoryList[ck]; ok {
		if p, ok := cat.Posts[path]; ok {
			return p
		}
	}
	return nil
}

func (s *site) build() {
	basePath := filepath.ToSlash(workPath + "/" + s.DocDir)
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		path = filepath.ToSlash(path)
		if path == basePath {
			return nil
		}
		path = path[len(s.DocDir):]

		if info.IsDir() {
			cat := NewCategory(info.Name(), path)
			s.CategoryList[path] = cat
			log.Println("初始化目录: ", path)
			return nil
		}

		if filepath.Ext(path) == ".md" {
			catkey := filepath.ToSlash(filepath.Dir(path))
			p, err := NewPost(basePath, path)
			if err != nil {
				return err
			}
			if cat, ok := s.CategoryList[catkey]; ok {
				cat.addPost(p)
				log.Println("初始化文章: ", path)
			}
		}

		return nil
	})

	checkError(err)
}

type category struct {
	Path  string
	Name  string
	Posts map[string]*post
}

func (c *category) addPost(p *post) {
	c.Posts[p.Path] = p
}

type post struct {
	Path    string
	Title   string
	Time    time.Time
	Url     string
	Content string
}

func NewCategory(name string, path string) *category {
	return &category{Name: name, Path: path, Posts: make(map[string]*post)}
}

func NewSite() *site {
	return &site{
		CategoryList: make(map[string]*category),
	}
}

func NewPost(base string, path string) (*post, error) {
	ret, err := ioutil.ReadFile(base + path)
	if err != nil {
		return nil, err
	}

	p := &post{
		Path:    path,
		Title:   filepath.Base(path)[0 : len(filepath.Base(path))-3],
		Content: string(blackfriday.MarkdownCommon(ret)),
	}

	return p, nil
}
