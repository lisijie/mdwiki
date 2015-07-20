package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ThemeDir    = "themes"
	TplVarLeft  = "{{"
	TplVarRight = "}}"
	T           *template.Template
)

var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

func init() {
	InitTemplate()
}

func RebuildTemplates(theme string) {
	InitTemplate()
	BuildTemplates(theme)
}

func InitTemplate() {
	T = template.New("__top__")
	funcMap := make(template.FuncMap)

	funcMap["str2html"] = func(raw string) template.HTML {
		return template.HTML(raw)
	}
	funcMap["GetPostListByCategory"] = func(cat string, page, size int) []*Post {
		return siteInfo.PostTable.GetPostListByCategory(cat, page, size)
	}
	funcMap["GetPostList"] = func(page, size int) []*Post {
		return siteInfo.PostTable.GetPostList(page, size)
	}
	funcMap["GetPost"] = func(link string) *Post {
		return siteInfo.PostTable.GetPostByPermalink(link)
	}
	funcMap["ShowPageBar"] = func(page, total int, url string) template.HTML {
		pageSize := config.GetInt("page_size", 10)
		pager := NewPager(page, total, pageSize, url, true)
		return template.HTML(pager.ToString())
	}
	funcMap["Date"] = func(t time.Time, format string) string {
		replacer := strings.NewReplacer(datePatterns...)
		format = replacer.Replace(format)
		return t.Format(format)
	}
	funcMap["GetCategoryList"] = func() []*Category {
		return siteInfo.CategoryTable.GetAll()
	}

	T.Funcs(funcMap)
}

func BuildTemplates(theme string) {
	root := filepath.Join(workPath, ThemeDir, theme)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".html" {
			source, err := ioutil.ReadFile(path)
			checkError(err)
			_, err = T.New(info.Name()).Parse(string(source))
			checkError(err)
		}
		return nil
	})
}

func RenderTemplate(w io.Writer, tpl string, data interface{}) error {
	return T.ExecuteTemplate(w, tpl, data)
}
