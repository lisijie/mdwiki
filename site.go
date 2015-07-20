package main

import (
	"log"
	"os"
	"path/filepath"
)

type Site struct {
	Name          string         //网站名称
	Keywords      string         //关键字
	Description   string         //描述信息
	Url           string         //网址
	StaticDir     []string       //静态文件目录
	PostDir       string         //md文档目录
	PostTable     *PostTable     //文章表
	CategoryTable *CategoryTable //分类表
}

func (s *Site) rebuild() {
	s.CategoryTable = NewCategoryTable()
	s.PostTable = NewPostTable()
	s.build()
}

func (s *Site) build() {
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
			post, err := MakePost(basePath, path)
			if err != nil {
				log.Println("文章", path, "解析失败: ", err)
				return nil
			}

			if _, err := s.CategoryTable.NewCategory(post.Category, filepath.ToSlash(filepath.Dir(path))); err != nil {
				log.Println(err)
				return nil
			}
			s.PostTable.AddPost(post)
		}

		return nil
	})

	// 更新分类统计
	for _, v := range s.CategoryTable.GetAll() {
		s.CategoryTable.UpdateCount(v.Name, s.PostTable.GetCountByCategory(v.Name))
	}

	checkError(err)
}

func NewSite() *Site {
	return &Site{
		PostTable:     NewPostTable(),
		CategoryTable: NewCategoryTable(),
	}
}
