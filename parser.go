package main

import (
	"encoding/json"
	"strings"

	"github.com/antchfx/htmlquery"
)

// 解析 herolist 文件，转换成 Hero 列表
// herolist 是字节流形式的 JSON
func ParseJson(data []byte) ([]Hero, error) {
	// 由于键值类型不定，需要用 interface{} 类型
	herolist := new([]map[string]any)

	err := json.Unmarshal(data, &herolist)
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}

	ret := make([]Hero, 0)
	for _, hero := range *herolist {
		ret = append(ret, Hero{
			hero["cname"].(string),
			int(hero["ename"].(float64)),
			hero["title"].(string),
		})
	}

	return ret, nil
}

// 解析每个英雄页面的 HTML ，返回皮肤列表
func ParseHtml(html string) ([]string, error) {
	root, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}

	ul := htmlquery.FindOne(root, "//div[@class=\"pic-pf\"]/ul")
	skins := htmlquery.SelectAttr(ul, "data-imgname")

	return SplitSkin(skins), nil
}
