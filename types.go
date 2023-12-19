package main

// 由于 herolist 内容不全面，格式不统一，只保留其中的
// cname ename title 三个字段封装成 Hero ，详细内容另外发送请求获取
type Hero struct {
	cname string
	ename int
	title string
}

type Chan1 struct {
	hero  Hero
	path  string
	skins []string
}

type Chan2 struct {
	content []byte
	path    string
}
