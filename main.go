package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"trie"
)

const (
	PORT     = 8080           //服务监听端口
	FILENAME = "badwords.txt" //敏感词库
)

var T *trie.Trie

//导入过滤词库
func importWords(T *trie.Trie, file string) (err error) {
	rd, err := os.Open(file)
	if err != nil {
		return
	}
	defer rd.Close()
	r := bufio.NewReader(rd)
	for {
		line, isPrefix, e := r.ReadLine()
		if e != nil {
			if e != io.EOF {
				err = e
			}
			break
		}
		if isPrefix {
			continue
		}
		if word := strings.TrimSpace(string(line)); word != "" {
			T.Add(word)
		}
	}
	return
}

//HTTP请求处理器
func mainHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	result, find := T.Replace(content)

	m := make(map[string]interface{})
	m["result"] = result
	m["find"] = find
	if len(find) > 0 {
		m["ret"] = 1
	} else {
		m["ret"] = 0
	}

	bytes, err := json.Marshal(m)
	if err == nil {
		w.Write(bytes)
	} else {
		log.Println(err.Error())
	}
}

func main() {
	var err error

	T = trie.NewTrie()

	err = importWords(T, FILENAME)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		log.Printf("服务正在启动，监听端口: %d ...\n", PORT)
		http.HandleFunc("/", mainHandler)
		err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
		if err != nil {
			log.Fatalln("启动失败: ", err.Error())
		}
	}
}
