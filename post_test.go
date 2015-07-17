package main

import (
	"fmt"

	"testing"
	"time"
)

func TestPostTable(t *testing.T) {
	table := NewPostTable()

	p1 := &Post{
		FilePath: "post1",
		Title:    "post 1",
		PostTime: time.Unix(100, 0),
	}

	p2 := &Post{
		FilePath: "post2",
		Title:    "post 2",
		PostTime: time.Unix(200, 0),
	}

	p3 := &Post{
		FilePath: "post3",
		Title:    "post 3",
		PostTime: time.Unix(300, 0),
	}

	table.AddPost(p3)
	table.AddPost(p1)
	table.AddPost(p2)

	fmt.Println(table.GetPostByFilePath("post1") == p1)

	for _, v := range table.list.data {
		fmt.Println(v.Title)
	}

	table.RemovePost(table.GetPostByFilePath("post1"))

	for _, v := range table.list.data {
		fmt.Println(v.Title)
	}

	lst := table.GetPostList(21, 2)
	fmt.Println(lst)

}
