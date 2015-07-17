package main

import (
	"errors"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

type Post struct {
	FilePath  string    // 文件物理路径
	Layout    string    // 使用的模版名
	Permalink string    // 文章的URL
	Category  string    // 分类名
	Title     string    // 文章标题
	Keywords  string    // 文章关键字
	Date      string    // 发布时间字符串
	PostTime  time.Time // 发布时间对象
	Content   string    //文章内容
}

type PostTable struct {
	posts        map[string]*Post
	linkMap      map[string]*Post
	list         *PostList
	categoryList map[string]*PostList
}

// 添加文章
func (pt *PostTable) AddPost(p *Post) bool {
	if _, ok := pt.posts[p.FilePath]; ok {
		return false
	}

	pt.posts[p.FilePath] = p
	pt.linkMap[p.Permalink] = p
	pt.list.Add(p)

	if v, ok := pt.categoryList[p.Category]; ok {
		v.Add(p)
	} else {
		pt.categoryList[p.Category] = NewPostList()
		pt.categoryList[p.Category].Add(p)
	}

	return true
}

// 删除文章
func (pt *PostTable) RemovePost(p *Post) bool {
	if _, ok := pt.posts[p.FilePath]; !ok {
		return false
	}

	delete(pt.posts, p.FilePath)
	delete(pt.linkMap, p.Permalink)
	pt.list.Remove(p)
	pt.categoryList[p.Category].Remove(p)

	return true
}

// 获取文章总数
func (pt *PostTable) GetCount() int {
	return pt.list.Count()
}

// 根据链接获取文章对象
func (pt *PostTable) GetPostByPermalink(link string) *Post {
	if v, ok := pt.linkMap[link]; ok {
		return v
	}
	return nil
}

// 根据文件路径获取文章对象
func (pt *PostTable) GetPostByFilePath(path string) *Post {
	if v, ok := pt.posts[path]; ok {
		return v
	}
	return nil
}

// 分页获取文章列表
func (pt *PostTable) GetPostList(page, pageSize int) []*Post {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize
	return pt.list.GetList(start, pageSize)
}

// 分页获取某个分类下的文章列表
func (pt *PostTable) GetPostListByCategory(cat string, page, pageSize int) []*Post {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize

	if v, ok := pt.categoryList[cat]; ok {
		return v.GetList(start, pageSize)
	}

	return make([]*Post, 0)
}

// 获取某个分类的文章总数
func (pt *PostTable) GetCountByCategory(cat string) int {
	if v, ok := pt.categoryList[cat]; ok {
		return v.Count()
	}
	return 0
}

type PostList struct {
	data []*Post
}

func (pl *PostList) Add(p *Post) {
	pl.data = append(pl.data, p)
	pl.Sort()
}

func (pl *PostList) Remove(p *Post) {
	for k, v := range pl.data {
		if v.FilePath == p.FilePath {
			data := make([]*Post, 0, len(pl.data)-1)
			data = append(data, pl.data[0:k]...)
			data = append(data, pl.data[k+1:len(pl.data)]...)
			pl.data = data
			break
		}
	}
	pl.Sort()
}

func (pl *PostList) GetList(start, size int) []*Post {
	if start > pl.Len() {
		return make([]*Post, 0)
	}
	end := start + size
	if end > pl.Len() {
		end = pl.Len()
	}
	return pl.data[start:end]
}

func (pl *PostList) Count() int {
	return len(pl.data)
}

func (pl *PostList) Len() int {
	return len(pl.data)
}

func (pl *PostList) Less(i, j int) bool {
	return pl.data[i].PostTime.Unix() > pl.data[j].PostTime.Unix()
}

func (pl *PostList) Swap(i, j int) {
	pl.data[j], pl.data[i] = pl.data[i], pl.data[j]
}

func (pl *PostList) Sort() {
	sort.Sort(pl)
}

func NewPostTable() *PostTable {
	return &PostTable{
		posts:        make(map[string]*Post),
		linkMap:      make(map[string]*Post),
		list:         NewPostList(),
		categoryList: make(map[string]*PostList),
	}
}

func NewPostList() *PostList {
	return &PostList{data: make([]*Post, 0)}
}

// 根据文件路径创建文章对象
func MakePost(base, path string) (*Post, error) {
	absPath := filepath.Join(base, path)
	source, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?sm)^(---\s*\n)(.*?)^(---\s*\n)`)
	matches := re.FindSubmatch(source)
	if len(matches) != 4 {
		return nil, errors.New("找不到文章头信息")
	}

	p := &Post{FilePath: absPath}

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
		if pt, err := time.Parse("2006-01-02", p.Date); err == nil { //只有日期
			p.PostTime = pt
		} else if pt, err := time.Parse("2006-01-02 15:04", p.Date); err == nil {
			p.PostTime = pt
		}
	}
	// 如果头信息解析失败或者没指定，则使用文件的修改时间
	if p.PostTime.IsZero() {
		if fi, err := os.Stat(filepath.Join(base, path)); err == nil {
			p.PostTime = fi.ModTime()
		}
	}

	// 切掉头信息，并解析为HTML
	idx := re.FindIndex(source)
	p.Content = string(blackfriday.MarkdownCommon(source[idx[1]-1:]))

	return p, nil
}
