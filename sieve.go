package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ==== 检测关键词 =====
type Sieve struct {
	mu   sync.RWMutex
	trie *node
}

func New() *Sieve {
	s := &Sieve{
		trie: &node{},
	}
	return s
}

// 简单添加关键词
func (s *Sieve) Add(words []string) (fail []string) {
	return s.add(words, 0)
}

// 从文本添加关键词，打标签并设定是否自动替换为*
func (s *Sieve) AddByFile(filename string, tag uint8, autoReplace bool) (fails []string, err error) {
	const delim = '\n'
	words := make([]string, 0, 2048)

	var reader io.Reader
	// 远程文件
	if strings.HasPrefix(filename, "http") {
		resp, err := http.Get(filename)
		if err != nil {
			return nil, err
		}
		reader = resp.Body
		defer resp.Body.Close()
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		reader = f
		defer f.Close()
	}

	br := bufio.NewReader(reader)
	for {
		b, err := br.ReadBytes(delim)
		words = append(words, string(b))
		if err == io.EOF {
			break
		}
	}

	fails = s.add(words, tag)

	return
}

// 移除关键词
func (s *Sieve) Remove(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		s.trie.RemoveWord(w)
	}
}

// 返回文本中第一个关键词及其标签
func (s *Sieve) Search(text string) (string, uint8) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws := []rune(text)
	start, end, tag, _, wildcards := s.trie.Search(ws)

	// wildcards存放模式中通配符位置，下列函数将文本中相应位置替换为*
	if len(wildcards) > 0 {
		for _, a := range wildcards {
			ws[a] = trans('*')
		}
	}

	return string(ws[start:end]), tag
}

// 添加关键词，打标签并设定是否强制替换
func (s *Sieve) add(words []string, tag uint8) (fail []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		if !s.trie.AddWord(w, tag) {
			fail = append(fail, w)
		}
	}

	return
}
