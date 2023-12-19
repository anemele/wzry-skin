package main

import (
	"os"
	"regexp"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func MkDir(path string) string {
	if !Exists(path) {
		os.Mkdir(path, os.ModePerm)
		Logger.Info("MKDIR", "path", path)
	}
	return path
}

func SplitSkin(skins string) []string {
	re, _ := regexp.Compile(`(\S+?)[\s(?:&\d+)\|]+`)
	matches := re.FindAllStringSubmatch(skins, -1)
	ret := make([]string, len(matches))
	for i, match := range matches {
		ret[i] = match[1]
	}
	return ret
}
