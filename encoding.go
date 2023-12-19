package main

/* 与编码相关的函数
HTTP 请求的结果是字节流，错误的编码会导致乱码，无法解析数据。
主要是 utf-8 和 gbk 之间的转换。 */

import "golang.org/x/text/encoding/simplifiedchinese"

/* func ConvertString_from_UTF8_to_GBK(str string) string {

	//utf-8编码的字符串
	ret, err := simplifiedchinese.GBK.NewEncoder().String(str)
	if err != nil {
		return "" //如果转换失败返回空字符串
	}
	return ret

}

func ConvertBytes_from_UTF8_to_GBK(b []byte) string {

	//[]byte格式的字符串
	bytes, err := simplifiedchinese.GBK.NewEncoder().Bytes(b)
	if err != nil {
		return "" //如果转换失败返回空字符串
	}
	return string(bytes)

}

func ConvertString_from_GBK_to_UTF8(str string) string {

	//GBK编码的字符串
	ret, err := simplifiedchinese.GBK.NewDecoder().String(str)
	if err != nil {
		return "" //如果转换失败返回空字符串
	}
	return ret

} */

func ConvertBytes_from_GBK_to_UTF8(b []byte) string {

	//[]byte格式的字符串
	bytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	if err != nil {
		Logger.Error(err.Error())
		return "" //如果转换失败返回空字符串
	}
	return string(bytes)

}
