//字典包，采用的是单词查找树
package trie

import (
	"log"
)

type Trie struct {
	Root *TrieNode
}

type TrieNode struct {
	Children map[rune]*TrieNode
	End      bool
}

//初始化trie数
func NewTrie() *Trie {
	t := &Trie{}
	t.Root = NewTrieNode()

	return t
}

//初始化一个节点
func NewTrieNode() *TrieNode {
	n := &TrieNode{}
	n.Children = make(map[rune]*TrieNode)
	n.End = false

	return n
}

//新增要过滤的词
func (this *Trie) Add(txt string) {
	if len(txt) < 1 {
		return
	}
	chars := []rune(txt)
	slen := len(chars)
	node := this.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; !exists {
			node.Children[chars[i]] = NewTrieNode()
		}
		node = node.Children[chars[i]]
	}
	node.End = true
}

//屏蔽字搜索替换
func (this *Trie) Replace(txt string) (string, []string) {
	chars := []rune(txt)
	result := []rune(txt)
	find := make([]string, 0, 10)
	slen := len(chars)
	node := this.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; exists {
			node = node.Children[chars[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[chars[j]]; !exists {
					break
				}
				node = node.Children[chars[j]]
				if node.End == true {
					for t := i; t <= j; t++ {
						result[t] = '*'
					}
					find = append(find, string(chars[i:j+1]))
					i = j
					node = this.Root
					break
				}
			}
			node = this.Root
		}
	}

	return string(result), find
}

//查找字符串
func (this *Trie) Find(txt string) bool {
	chars := []rune(txt)
	slen := len(chars)
	node := this.Root
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[chars[i]]; exists {
			node = node.Children[chars[i]]
			//若全部字符都存在匹配，判断最终停留的节点是否为树叶，若是，则返回真，否则返回假。
			if node.End == true && i == slen-1 {
				return true
			}
		}
	}

	return false
}

//删除字符串
//首先查找该字符串，边查询边将经过的节点压栈，若找不到，则返回假；否则依次判断栈顶节点是否为树叶，
//若是则删除该节点，否则返回真。
func (this *Trie) Delete(txt string) bool {
	log.Printf("正在查找: %s ...\n", txt)
	s := []rune(txt)
	slen := len(s)
	if slen == 0 {
		return true
	}

	node := this.Root
	child, ok := node.Children[s[0]]
	if ok && child.delete(txt[1:]) {
		//delete(node.Children, s[0])
		this.Root.Children[s[0]].Children = nil
		this.Root.Children[s[0]].End = false
	}
	return false
}

//删除字符串
func (this *TrieNode) delete(txt string) bool {
	log.Printf("正在查找: %s ...\n", txt)
	s := []rune(txt)
	slen := len(s)
	if slen == 0 {
		return true
	}

	child, ok := this.Children[s[0]]
	if ok {
		if slen > 1 {
			if !child.delete(txt[1:]) {
				return false
			}
		}

		//delete(this.Children, s[0])
		this.Children[s[0]].Children = nil
		this.Children[s[0]].End = false
		return true
	}
	return false
}
