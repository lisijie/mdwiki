package main

import (
	"errors"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type category struct {
	Name      string
	Permalink string
	Posts     map[string]*post
	PostCount int
}

type post struct {
	Layout    string
	Permalink string
	Category  string
	Title     string
	Keywords  string
	Date      string
	PostTime  time.Time
	Content   string
}

type site struct {
	Title        string   //网站标题
	Keywords     string   //关键字
	Description  string   //描述信息
	Url          string   //网址
	StaticDir    []string //静态文件目录
	PostDir      string   //md文档目录
	CategoryList map[string]*category
	Posts        map[string]*post
}

func (s *site) getPost(uri string) *post {
	if post, ok := s.Posts[uri]; ok {
		return post
	}
	return nil
}

func (s *site) getCategory(uri string) *category {
	if cat, ok := s.CategoryList[uri]; ok {
		return cat
	}
	return nil
}

func (s *site) getCategoryByName(name string) *category {
	for _, v := range s.CategoryList {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (s *site) rebuild() {
	s.CategoryList = make(map[string]*category)
	s.Posts = make(map[string]*post)
	s.build()
}

func (s *site) build() {
	basePath := filepath.ToSlash(workPath + "/" + s.PostDir)
	baseLen := len(basePath)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		defer func() {
			if err := recover(); err != nil {
				debug("panic ", err)
				return
			}
		}()

		path = filepath.ToSlash(path)
		if path == basePath {
			return nil
		}
		path = path[baseLen:]

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			post, err := s.makePost(basePath, path)
			if err != nil {
				log.Println("文章", path, "解析失败: ", err)
				return nil
			}

			s.Posts[post.Permalink] = post
			if cat := s.getCategoryByName(post.Category); cat != nil {
				cat.addPost(post)
			} else {
				cat := NewCategory(post.Category, filepath.ToSlash(filepath.Dir(path)))
				cat.addPost(post)
				s.CategoryList[cat.Permalink] = cat
			}
		}

		return nil
	})

	checkError(err)
}

// 创建文章对象
func (s *site) makePost(base, path string) (*post, error) {
	source, err := ioutil.ReadFile(filepath.Join(base, path))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?sm)^(---\s*\n)(.*?)^(---\s*\n)`)
	matches := re.FindSubmatch(source)
	if len(matches) != 4 {
		return nil, errors.New("找不到文章头信息")
	}

	p := &post{}

	if err := yaml.Unmarshal(matches[2], p); err != nil {
		return nil, err
	}

	if p.Title == "" {
		return nil, errors.New("缺少文章头信息: title")
	}
	if p.Category == "" {
		return nil, errors.New("缺少文章头信息: category")
	}
	if p.Permalink == "" {
		p.Permalink = path
	}

	// 文章发布时间，如果头信息有指定日期，则使用头信息的日期进行解析
	if p.Date != "" {
		if pt, err := time.Parse("2006-01-02", p.Date); err != nil { //只有日期
			p.PostTime = pt
		} else if pt, err := time.Parse("2006-01-02 15:04", p.Date); err != nil {
			p.PostTime = pt
		}
	}
	// 如果头信息解析失败或者没指定，则使用文件的修改时间
	if p.PostTime.IsZero() {
		if fi, err := os.Stat(filepath.Join(base, path)); err != nil {
			p.PostTime = fi.ModTime()
		}
	}

	// 切掉头信息，并解析为HTML
	idx := re.FindIndex(source)
	p.Content = string(blackfriday.MarkdownCommon(source[idx[1]-1:]))

	return p, nil
}

// 分类增加文章
func (c *category) addPost(p *post) {
	c.Posts[p.Permalink] = p
	c.PostCount = len(c.Posts)
}

func NewCategory(name string, link string) *category {
	return &category{
		Name:      name,
		Permalink: link,
		Posts:     make(map[string]*post),
	}
}

func NewSite() *site {
	return &site{
		CategoryList: make(map[string]*category),
		Posts:        make(map[string]*post),
	}
}
