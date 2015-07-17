# mdwiki


mdwiki 是一个Go语言开发的wiki引擎，可用于构建wiki系统或个人博客。

## 特点

- 使用简单
- 不需要数据库，基于文档
- 使用markdown编写文档
- 极速访问

## 安装

	go get -u github.com/lisijie/mdwiki

## 使用说明

### 目录说明：

- static 静态资源目录
- themes 主题模版目录，每个主题一个目录
- posts 文档目录，每篇文档一个markdown文件，建议分2级目录存储

### 创建文档：

在posts目录下新建一个markdown文档，在文件头插入以下信息

    ---
    layout: 文章使用的模版文件
    title: 文章标题
    category: 文章类别
    keywords: 文章关键字
    ---

### 模板函数：

##### str2html

将输出字符串转成html，示例：

	{{str2html .post.Content}}

##### GetPostListByCategory

根据分类名获取某个分类文章列表，示例：

	{{$postList := GetPostListByCategory "foo" 1 10}}
	<ul>
	{{range $kk, $p := $postList}}
		<li>{{$p.Title}}</li>
	{{end}}
	</ul>

##### GetPostList

获取最新的文章列表，示例：

	{{$postList := GetPostList 1 10}}
	<ul>
	{{range $kk, $p := $postList}}
		<li>{{$p.Title}}</li>
	{{end}}
	</ul>

##### ShowPageBar

显示分页栏，生成的是bootstrap的分页格式，示例：

	{{ShowPageBar 1 100 "/foo"}}
