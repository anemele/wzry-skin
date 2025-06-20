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

const ApiUrl = "https://pvp.qq.com/web201605/js/herolist.json"
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

func (hero Hero) String() string {
	return fmt.Sprintf(
		"%s_%s, (%d %s), %s",
		hero.CName, hero.Title, hero.EName, hero.IdName, hero.SkinName,
	)
}

func (hero Hero) DirName() string {
	return fmt.Sprintf("%s_%s", hero.CName, hero.Title)
}

type Skin struct {
	Idx  int
	Name string
}

func (skin Skin) FileName() string {
	return fmt.Sprintf("%d_%s.jpg", skin.Idx, skin.Name)
}

func (hero Hero) GetSkins() []Skin {
	var skins []Skin
	for i, name := range strings.Split(hero.SkinName, "|") {
		skins = append(skins, Skin{Idx: i + 1, Name: name})
	}
	return skins
}

var client http.Client

func downloadSkin(hero Hero, skin Skin, heroDir string, wg *sync.WaitGroup) error {
	defer wg.Done()

	skinPath := path.Join(heroDir, skin.FileName())
	if exists(skinPath) {
		return fmt.Errorf("exists: %s", skinPath)
	}

	skinUrl := getSkinUrl(hero.EName, skin.Idx)

	resp, err := client.Get(skinUrl)
	if err != nil {
		return fmt.Errorf("%s: %s", err, skinUrl)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s %s", resp.Status, skinUrl)
	}
	defer resp.Body.Close()

	fp, err := os.OpenFile(skinPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fp, resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Downloaded", skinPath)

	return nil
}

func downloadHero(hero Hero, wg *sync.WaitGroup) {
	defer wg.Done()

	heroDir := path.Join(LocalDir, hero.DirName())
	ensurePath(heroDir)

	skins := hero.GetSkins()
	wg.Add(len(skins))
	for _, skin := range skins {
		go func() {
			if err := downloadSkin(hero, skin, heroDir, wg); err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func run() error {
	resp, err := client.Get(ApiUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}
	defer resp.Body.Close()

	var heroes []Hero
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&heroes)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(len(heroes))
	for _, hero := range heroes {
		go downloadHero(hero, &wg)
	}

	wg.Wait()
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
