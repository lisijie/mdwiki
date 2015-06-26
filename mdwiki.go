package main

import (
	"flag"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	VERSION = "1.0"
	RELEASE = ""
	DEBUG   = true
)

var (
	config   *goconfig.ConfigFile
	siteInfo *site
	workPath string
	docPath  string
	confFile = flag.String("conf", "./config.ini", "配置文件路径")
	httpAddr = flag.String("http-addr", ":8080", "HTTP接口端口，默认为:8080")
	showVer  = flag.Bool("version", false, "显示版本信息")
)

func main() {
	flag.Parse()

	//显示版本信息
	if *showVer {
		fmt.Printf("mdwiki v%s(%s)\n", VERSION, RELEASE)
		return
	}

	//加载配置文件
	config, err := goconfig.LoadConfigFile(*confFile)
	checkError(err)

	//工作目录,web根目录
	workPath = config.MustValue("site", "work_dir", ".")
	workPath, err = filepath.Abs(workPath)
	checkError(err)

	//初始化配置信息
	siteInfo = NewSite()
	siteInfo.Url = config.MustValue("site", "url")
	siteInfo.Title = config.MustValue("site", "title")
	siteInfo.Keywords = config.MustValue("site", "keywords")
	siteInfo.Description = config.MustValue("site", "description")
	siteInfo.PostDir = config.MustValue("site", "post_dir", "posts")
	siteInfo.StaticDir = strings.Split(config.MustValue("site", "static_dir", "static"), ",")

	siteInfo.build()

	BuildTemplates(config.MustValue("site", "theme", "default"))

	go fswatch(filepath.Join(workPath, siteInfo.PostDir), func() {
		debug("重新构建网站数据...")
		siteInfo.rebuild()
	})

	go fswatch(filepath.Join(workPath, ThemeDir), func() {
		debug("重新编译模版...")
		RebuildTemplates(config.MustValue("site", "theme", "default"))
	})

	//启动HTTP服务
	http := &httpServer{}
	http.run(*httpAddr)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func debug(s ...interface{}) {
	log.Println(s...)
}