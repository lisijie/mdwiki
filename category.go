package main

type Category struct {
	Name      string
	Permalink string
	PostCount int
}

type CategoryTable struct {
	data map[string]*Category
}

func (t *CategoryTable) NewCategory(name, link string) bool {
	if _, ok := t.data[name]; ok {
		return false
	}
	t.data[name] = &Category{
		Name:      name,
		Permalink: link,
		PostCount: 0,
	}
	return true
}

func (t *CategoryTable) UpdateCount(name string, count int) bool {
	if v, ok := t.data[name]; ok {
		v.PostCount = count
		return true
	}
	return false
}

func (t *CategoryTable) GetAll() map[string]*Category {
	return t.data
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
