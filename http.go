package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type httpServer struct {
	listener net.Listener
	rsp      http.ResponseWriter
	req      *http.Request
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

	s.rsp = w
	s.req = req

	if req.URL.Path == "/" {
		s.handleDefaultPage()
		return
	}

	for _, v := range siteInfo.staticDir {
		if strings.Index(req.URL.Path, "/"+v+"/") == 0 {
			s.handleStaticFile()
			return
		}
	}

}

func (s *httpServer) handleDefaultPage() {
	log.Println(siteInfo)

}

// 处理静态文件
func (s *httpServer) handleStaticFile() {

	fileName := workPath + strings.Replace(s.req.URL.Path, "..", "", -1)

	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		s.errorPage(404)
		return
	}

	s.rsp.Write(raw)
}

// 显示错误页面
func (s *httpServer) errorPage(code int) {
	s.rsp.WriteHeader(code)
	s.rsp.Write([]byte("<h1>" + strconv.Itoa(code) + " " + http.StatusText(404) + "</h1>"))
	s.rsp.Write([]byte("<hr /> <span style=\"font-size:11px\">Powered by mdwiki " + VERSION + "</span>"))
}
