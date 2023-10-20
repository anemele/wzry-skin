package wzry

import (
	"fmt"
	"log"
	"path"
	"sync"
)

// size 是随便设置的
const size = 256

var chan1 = make(chan Chan1, size)
var chan2 = make(chan Chan2, size)

func getHeroPage(wg *sync.WaitGroup, hero Hero) {
	defer wg.Done()

	// 创建该英雄的存储目录
	heroTN := hero.cname + "_" + hero.title
	heroSavePath := path.Join(SaveRoot, heroTN)
	MkDir(heroSavePath)
	// 请求英雄页面
	html, err := GetPage(hero.ename)
	if err != nil {
		log.Println(err)
		return
	}

	// 解析英雄页面，获取皮肤列表
	skins, err := ParseHtml(html)
	if err != nil {
		log.Println(err)
		return
	}

	chan1 <- Chan1{hero, heroSavePath, skins}
}

func getSkinBytes(wg *sync.WaitGroup, url, path string) {
	defer wg.Done()

	bytes, err := GetBytes(url)
	if err == nil {
		chan2 <- Chan2{bytes, path}
	} else {
		log.Println(err)
	}
}

func saveSkin(wg *sync.WaitGroup, bytes []byte, path string) {
	defer wg.Done()

	ok, err := WriteBytes(bytes, path)
	if ok {
		LogInfo("SAVED", path)
	} else {
		log.Println(err)
	}
}

func Run() error {
	// 获取 []Hero
	heros, err := GetData()
	if err != nil {
		return err
	}

	// 获取本地统计信息 statistics.txt
	statistics, err := GetStat()
	if err != nil {
		return err
	}

	// 创建同步锁
	var wg sync.WaitGroup
	// 遍历英雄列表
	wg.Add(len(heros))
	for _, hero := range heros {
		go getHeroPage(&wg, hero)
	}

	count := 0
	for i := 0; i < len(heros); i++ {
		ch1 := <-chan1
		hero := ch1.hero
		heroSavePath := ch1.path
		skins := ch1.skins

		// 截取皮肤列表，更新统计信息
		lenSkin := len(skins)
		lenStat := statistics[hero.cname]
		// 比对皮肤列表长度和统计信息
		if lenStat < lenSkin {
			// 统计信息记录小于皮肤列表长度，说明有更新
			// 则截取更新部分送入下载，并更新统计信息记录
			for i, skin := range skins[lenStat:] {
				skinFileName := fmt.Sprintf("%d_%s.jpg", i+lenStat+1, skin)
				skinSavePath := path.Join(heroSavePath, skinFileName)
				if Exists(skinSavePath) {
					LogInfo("EXISTS", skinSavePath)
					continue
				}
				skinImageUrl := GetImageUrl(hero.ename, i+lenStat+1, SkinSize["b"])
				wg.Add(1)
				go getSkinBytes(&wg, skinImageUrl, skinSavePath)
				count++
			}
			statistics[hero.cname] = lenSkin
		} else if lenStat > lenSkin {
			// 统计信息记录大于皮肤列表长度，说明记录存在错误
			// （也可能是请求的皮肤数据错误或解析错误，可能性较小）
			// 仅更正统计信息记录
			statistics[hero.cname] = lenSkin
		}
		// 如果二者相等，说明没有更新，也没有错误，无需操作
	}

	wg.Add(count)
	for i := 0; i < count; i++ {
		ch2 := <-chan2
		go saveSkin(&wg, ch2.content, ch2.path)
	}

	wg.Wait()
	SetStat(statistics)
	log.Println("DONE")

	return nil
}
