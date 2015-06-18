package main

import (
	"io/ioutil"
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
	for _, v := range siteInfo.staticDir {
		if strings.Index(req.URL.Path, "/"+v+"/") == 0 {
			s.handleStaticFile(w, req)
		}
	}

}

// 处理静态文件
func (s *httpServer) handleStaticFile(w http.ResponseWriter, req *http.Request) {

	fileName := siteInfo.workDir + strings.Replace(req.URL.Path, "..", "", -1)

	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		s.errorPage(w, 404)
		return
	}

	w.Write(raw)
}

// 显示错误页面
func (s *httpServer) errorPage(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte("<h1>" + strconv.Itoa(code) + " " + http.StatusText(404) + "</h1>"))
	w.Write([]byte("<hr /> <span style=\"font-size:11px\">Powered by mdwiki " + VERSION + "</span>"))
}
