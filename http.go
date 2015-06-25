package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type httpServer struct {
	listener net.Listener
}

func (s *httpServer) run(addr string) {
	log.Println("启动HTTP服务...")
	listener, err := net.Listen("tcp", addr)
	checkError(err)
	s.listener = listener
	server := &http.Server{
		Handler: s,
	}
	server.Serve(s.listener)
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Form == nil {
		req.ParseForm()
	}

	if req.URL.Path == "/" || req.URL.Path == "/index.html" {
		s.indexHandler(w, req)
		return
	}

	for _, v := range siteInfo.StaticDir {
		if strings.Index(req.URL.Path, "/"+v+"/") == 0 {
			s.staticHandler(w, req)
			return
		}
	}

	s.pageHandler(w, req)
}

// 首页
func (s *httpServer) indexHandler(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	data["title"] = siteInfo.Title
	data["keywords"] = siteInfo.Keywords
	data["description"] = siteInfo.Description
	data["categoryList"] = siteInfo.CategoryList

	RenderTemplate(w, "index.html", data)
}

// 处理页面
func (s *httpServer) pageHandler(w http.ResponseWriter, req *http.Request) {
	post := siteInfo.GetPost(req.URL.Path)
	if post == nil {
		s.errorPage(w, 404)
		return
	}

	data := make(map[string]interface{})
	data["title"] = siteInfo.Title
	data["keywords"] = siteInfo.Keywords
	data["description"] = siteInfo.Description
	data["categoryList"] = siteInfo.CategoryList
	data["post"] = post

	RenderTemplate(w, "page.html", data)
}

// 处理静态文件
func (s *httpServer) staticHandler(w http.ResponseWriter, req *http.Request) {
	fileName := workPath + strings.Replace(req.URL.Path, "..", "", -1)
	http.ServeFile(w, req, fileName)
}

// 显示错误页面
func (s *httpServer) errorPage(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte("<h1>" + strconv.Itoa(code) + " " + http.StatusText(404) + "</h1>"))
	w.Write([]byte("<hr /> <span style=\"font-size:11px\">Powered by mdwiki " + VERSION + "</span>"))
}
