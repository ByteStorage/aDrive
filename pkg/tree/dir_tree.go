package tree

import (
	"fmt"
	"log"
	"strings"
)

type DirTree struct {
	Root *DirTreeNode
}

type DirTreeNode struct {
	Name     string
	Children []*DirTreeNode
}

// Insert 插入目录 参数应带/
func (tree *DirTree) Insert(path string) bool {
	// len(strings.Split("/")) == 2 ["", ""]
	// len(strings.Split("/hello/")) == 3 ["", "hello", ""]
	ancestors := strings.Split(path, "/")[1:]

	if len(ancestors) < 1 && ancestors[0] == "" {
		return false
	}

	var parent *DirTreeNode = tree.Root
	var current = parent

	//确保所有的父节点都存在
	for i := 0; i < len(ancestors)-2; i++ {
		parent = tree.Seek(parent, ancestors[i])
		if parent == nil {
			//父节点不存在
			current.Children = append(current.Children, &DirTreeNode{Name: ancestors[i]})
			parent = current
			parent = tree.Seek(parent, ancestors[i])
			current = parent
		} else {
			current = parent
		}
	}

	if parent == nil {
		return false
	}

	dirName := ancestors[len(ancestors)-2]

	parent.Children = append(parent.Children, &DirTreeNode{Name: dirName})

	return true
}

// Seek 从node开始查找name节点
func (tree *DirTree) Seek(node *DirTreeNode, name string) *DirTreeNode {
	if node == nil {
		return nil
	}

	if node.Name == name {
		return node
	}

	if node.Children != nil && len(node.Children) > 0 {
		for i := 0; i < len(node.Children); i++ {
			node := tree.Seek(node.Children[i], name)
			if node != nil {
				return node
			}
		}
	}
	return nil
}

// FindSubDir 查找路径下的子目录 参数应带上/
func (tree *DirTree) FindSubDir(path string) (subDirs []string) {
	ancestors := strings.Split(path, "/")[1:]

	if len(ancestors) < 1 && ancestors[0] == "" {
		return
	}

	var parent *DirTreeNode = tree.Root

	if len(ancestors) == 1 && ancestors[0] == "" {
		for _, child := range parent.Children {
			subDirs = append(subDirs, child.Name)
		}

		return
	}

	// 找到最后一个parent
	for i := 0; i < len(ancestors)-1; i++ {
		for _, c := range parent.Children {
			if ancestors[i] == c.Name {
				parent = c
				continue
			}
		}
	}

	//如果目录不存在，应该报错你要查到的目录不存在
	if parent.Name != ancestors[len(ancestors)-2] {
		log.Println("你要查找的目录不存在")
		return
	}

	for _, child := range parent.Children {
		subDirs = append(subDirs, child.Name)
	}

	return
}

// LookAll 调试用, DFS查看整个目录树的内容
func (tree *DirTree) LookAll() string {
	nodes := make([]string, 0)
	// 初始化队列
	queue := []*DirTreeNode{tree.Root}
	// 当队列中没有元素，那么结束
	for len(queue) > 0 {
		var count = 0
		for i := range queue {
			// 计数+1
			count++
			// 保存值
			nodes = append(nodes, queue[i].Name)
			// 子节点入队
			for j := range queue[i].Children {
				queue = append(queue, queue[i].Children[j])
			}

		}
		// 类似于出队，将遍历过的删掉
		queue = queue[count:]

	}

	return fmt.Sprintf("%s", nodes) // [1 2 3 4 5 6 7 8 9 10]
}

func (tree *DirTree) Rename(node *DirTreeNode, old string, new string) {
	split := strings.Split(old, "/")
	dirTreeNode := tree.Seek(node, split[len(split)-2])
	newSplit := strings.Split(new, "/")
	dirTreeNode.Name = newSplit[len(newSplit)-2]
}

func (tree *DirTree) Delete(node *DirTreeNode, path string) {
	split := strings.Split(path, "/")
	i := split[1 : len(split)-1]
	if len(i) == 1 {
		for k, v := range node.Children {
			if v.Name == i[0] {
				node.Children = removeElementFromSlice(node.Children, k)
			}
		}
		return
	}
	for _, s := range i[:len(i)-1] {
		node = tree.Seek(node, s)
	}
	for k, c := range node.Children {
		if c.Name == i[len(i)-1] {
			node.Children = removeElementFromSlice(node.Children, k)
		}
	}
}

func removeElementFromSlice(elements []*DirTreeNode, index int) []*DirTreeNode {
	return append(elements[:index], elements[index+1:]...)
}
