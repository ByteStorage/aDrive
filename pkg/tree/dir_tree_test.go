package tree

import (
	"fmt"
	"testing"
)

// setup 初始化
//          /
//  2       3       4
//5 6 7    8 9     10
func setup() *DirTree {
	root := &DirTreeNode{
		Name:     "/",
		Children: []*DirTreeNode{},
	}

	return &DirTree{root}
}

func TestLookAll(t *testing.T) {
	tree := setup()
	t.Log(tree.LookAll())
}

// TestInsert 结果应为
//          /
//  2       3       4
//5 6 7    8 9     10
//11
func TestInsert(t *testing.T) {
	tree := setup()
	//插入目录+文件
	tree.Insert("/first.txt/")
	tree.Insert("/tds/hello.txt/")
	tree.Insert("/tds/hdfs.txt/")
	tree.Insert("/tds/hello/hello.txt/")

	t.Log(tree.FindSubDir("/tds/"))
	t.Log(tree.FindSubDir("/"))
	t.Log(tree.FindSubDir("/hello/"))
	t.Log(tree.FindSubDir("/tds/hello/"))

	tree.Rename(tree.Root, "/tds/hello/", "/tds/test/")
	t.Log(tree.FindSubDir("/tds/"))
	t.Log(tree.FindSubDir("/tds/hello/"))
	t.Log(tree.FindSubDir("/tds/test/"))
}

func TestDirTree_Delete(t *testing.T) {
	tree := setup()
	//插入目录+文件
	tree.Insert("/first.txt/")
	tree.Insert("/tds/hello.txt/")
	tree.Insert("/tds/hdfs.txt/")
	tree.Insert("/tds/hello/hello.txt/")
	tree.Delete(tree.Root, "/tds/hello.txt/")
	dir := tree.FindSubDir("/tds/")
	fmt.Println(dir)
}

// TestFindSubDir 结果应为 [5, 6, 7]
func TestFindSubDir(t *testing.T) {
	tree := setup()
	t.Log(tree.FindSubDir("/tds/"))
	t.Log("以上结果应为 [5, 6, 7]")
}
