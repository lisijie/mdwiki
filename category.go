package main

import (
	"errors"
	"sort"
)

type Category struct {
	Name      string
	Permalink string
	PostCount int
}

type CategoryTable struct {
	data map[string]*Category
}

func (t *CategoryTable) NewCategory(name, link string) (bool, error) {
	if _, ok := t.data[name]; ok {
		return false, nil
	}
	for _, v := range t.data {
		if v.Permalink == link {
			return false, errors.New("分类 " + name + " 与 " + v.Name + " 的URL重复: " + link)
		}
	}
	t.data[name] = &Category{
		Name:      name,
		Permalink: link,
		PostCount: 0,
	}
	return true, nil
}

func (t *CategoryTable) UpdateCount(name string, count int) bool {
	if v, ok := t.data[name]; ok {
		v.PostCount = count
		return true
	}
	return false
}

// 获取所有分类，返回排好序的列表
func (t *CategoryTable) GetAll() []*Category {
	tmp := make(sort.StringSlice, 0, len(t.data))
	for k, _ := range t.data {
		tmp = append(tmp, k)
	}
	tmp.Sort()
	ret := make([]*Category, 0, len(t.data))
	for _, v := range tmp {
		ret = append(ret, t.data[v])
	}
	return ret
}

func (t *CategoryTable) GetByName(name string) *Category {
	if v, ok := t.data[name]; ok {
		return v
	}
	return nil
}

func (t *CategoryTable) GetByPermalink(link string) *Category {
	for _, v := range t.data {
		if v.Permalink == link {
			return v
		}
	}
	return nil
}

func NewCategoryTable() *CategoryTable {
	return &CategoryTable{
		data: make(map[string]*Category),
	}
}
