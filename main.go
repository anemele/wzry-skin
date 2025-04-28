package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

const ApiUrl = "http://pvp.qq.com/web201605/js/herolist.json"
const LocalDir = "wzry-skin-dirs"

func getSkinUrl(ename int, idx int) string {
	return fmt.Sprintf(
		"https://game.gtimg.cn/images/yxzj/img201606/heroimg/%d/%d-bigskin-%d.jpg",
		ename, ename, idx,
	)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func ensurePath(path string) {
	if !exists(path) {
		os.Mkdir(path, 0755)
	}
}

func init() {
	ensurePath(LocalDir)
}

type Hero struct {
	CName    string `json:"cname"`
	EName    int    `json:"ename"`
	Title    string `json:"title"`
	IdName   string `json:"id_name"`
	SkinName string `json:"skin_name"`
}

func (h Hero) String() string {
	return fmt.Sprintf(
		"%s_%s, (%d %s), %s",
		h.CName, h.Title, h.EName, h.IdName, h.SkinName,
	)
}

type Skin struct {
	Idx  int
	Name string
}

func (h Hero) GetSkins() []Skin {
	skins := make([]Skin, 0)
	for i, skin := range strings.Split(h.SkinName, "|") {
		skins = append(skins, Skin{Idx: i + 1, Name: skin})
	}
	return skins
}

var client http.Client

func downloadSkin(hero Hero, skin Skin, heroDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	skinUrl := getSkinUrl(hero.EName, skin.Idx)
	skinPath := path.Join(heroDir, fmt.Sprintf("%d_%s.jpg", skin.Idx, skin.Name))
	if !exists(skinPath) {
		resp, err := client.Get(skinUrl)
		if err != nil {
			fmt.Println(err, skinUrl)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error:", resp.Status, skinUrl)
			return
		}
		fp, err := os.OpenFile(skinPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fp.Close()
		_, err = io.Copy(fp, resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Downloaded", skinPath)
	}
}

func downloadHero(hero Hero, wg *sync.WaitGroup) {
	defer wg.Done()
	heroDir := path.Join(LocalDir, fmt.Sprintf("%s_%s", hero.CName, hero.Title))
	ensurePath(heroDir)

	skins := hero.GetSkins()
	wg.Add(len(skins))
	for _, skin := range skins {
		go downloadSkin(hero, skin, heroDir, wg)
	}
}

func main() {
	resp, err := client.Get(ApiUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: ", resp.Status)
		return
	}

	var heroes []Hero
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&heroes)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(len(heroes))
	for _, hero := range heroes {
		go downloadHero(hero, &wg)
	}

	wg.Wait()
}
