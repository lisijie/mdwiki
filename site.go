package main

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
	baseLen := len(basePath)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		path = filepath.ToSlash(path)
		debug("walk ", path)
		if path == basePath {
			return nil
		}
		path = path[baseLen:]

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			_, err := s.makePost(basePath, path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	checkError(err)
}

func (s *site) makePost(base, path string) (*post, error) {
	source, err := ioutil.ReadFile(filepath.Join(base, path))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?sm)^(---\s*\n)(.*?)^(---\s*\n)`)

	m := re.Find(source)

	debug("m:", m)

	p := &post{
		Url:     path,
		Content: string(blackfriday.MarkdownCommon(source)),
	}

	return p, nil
}

type category struct {
	Name  string
	Posts map[string]*post
}

type post struct {
	Url     string
	Title   string
	Time    time.Time
	Content string
}

func NewCategory(name string, path string) *category {
	return &category{Name: name, Posts: make(map[string]*post)}
}

func NewSite() *site {
	return &site{
		CategoryList: make(map[string]*category),
	}
}
