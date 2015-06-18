package main

import (
	"flag"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"strings"
)

const (
	VERSION = "1.0"
	RELEASE = ""
)

var (
	config   *goconfig.ConfigFile
	siteInfo *site
	dirTree  map[string]*dir
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
	if err != nil {
		fmt.Println(err)
		return
	}

	//初始化配置信息
	siteInfo = new(site)
	siteInfo.url = config.MustValue("site", "url")
	siteInfo.title = config.MustValue("site", "title")
	siteInfo.keywords = config.MustValue("site", "keywords")
	siteInfo.description = config.MustValue("site", "description")
	siteInfo.docDir = config.MustValue("site", "doc_dir", "docs")
	siteInfo.staticDir = strings.Split(config.MustValue("site", "static_dir", "static"), ",")
	siteInfo.workDir = config.MustValue("site", "work_dir", ".")

	dirTree = make(map[string]*dir)
	build(siteInfo.workDir + "/" + siteInfo.docDir)

	//启动HTTP服务
	http := &httpServer{}
	http.run(*httpAddr)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
