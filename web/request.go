package web

import (
	"consts"
	"errors"
)

// 请求 herolist.json ，解析转换返回 Hero 列表
func GetData() ([]Hero, error) {
	bytes, err := GetBytes(consts.ApiDataUrl)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	if bytes == nil {
		return nil, errors.New("herolist.json is empty")
	}

	return ParseJson(bytes)
}

// 请求英雄页面，返回 HTML 字符串
func GetPage(ename int) (string, error) {
	heroPageUrl := consts.GetPageUrl(ename)
	pageBytes, err := GetBytes(heroPageUrl)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	html := ConvertBytes_from_GBK_to_UTF8(pageBytes)

	return html, nil
}
