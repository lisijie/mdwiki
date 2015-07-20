package main

import (
	"flag"
	"fmt"
	"github.com/lisijie/go-conf"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	VERSION = "1.0"
	RELEASE = ""
)

var (
	config   *goconf.Config
	siteInfo *Site
	workPath string
	docPath  string
	confFile = flag.String("conf", "./config.ini", "配置文件路径")
	httpAddr = flag.String("http-addr", ":8080", "HTTP接口端口，默认为:8080")
	showVer  = flag.Bool("version", false, "显示版本信息")
	isDebug  = flag.Bool("debug", true, "是否启用调试模式")
)

func main() {
	flag.Parse()

	// 显示版本信息
	if *showVer {
		fmt.Printf("mdwiki v%s(%s)\n", VERSION, RELEASE)
		return
	}

	// 加载配置文件
	var err error
	config, err = goconf.NewConfig(*confFile)
	checkError(err)

	// 工作目录,为配置文件所在目录
	workPath, err = filepath.Abs(filepath.Dir(*confFile))
	checkError(err)
	debug("工作目录:", workPath)

	//初始化配置信息
	siteInfo = NewSite()
	siteInfo.Url = config.GetString("site_url")
	siteInfo.Name = config.GetString("site_name")
	siteInfo.Keywords = config.GetString("site_keywords")
	siteInfo.Description = config.GetString("site_description")
	siteInfo.PostDir = config.GetString("post_dir", "posts")
	siteInfo.StaticDir = strings.Split(config.GetString("static_dir", "static"), ",")

	siteInfo.build()

	BuildTemplates(config.GetString("theme", "default"))

	go fswatch(filepath.Join(workPath, siteInfo.PostDir), func() {
		debug("重新构建网站数据...")
		siteInfo.rebuild()
	})

	go fswatch(filepath.Join(workPath, ThemeDir), func() {
		debug("重新编译模版...")
		RebuildTemplates(config.GetString("theme", "default"))
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
	if *isDebug {
		log.Println(s...)
	}
}
