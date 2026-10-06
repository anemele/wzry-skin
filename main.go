package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	ApiUrl   = "https://pvp.qq.com/web201605/js/herolist.json"
	LocalDir = "skin-dirs"
)

// mkdir if not exists?
// right?
func ensureExists(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(path, 0o755)
		}
	}
}

func getSkinUrl(ename int, idx int) string {
	return fmt.Sprintf(
		"https://game.gtimg.cn/images/yxzj/img201606/heroimg/%d/%d-bigskin-%d.jpg",
		ename, ename, idx,
	)
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

func main() {
	ensureExists(LocalDir)

	c1 := colly.NewCollector(
		colly.Async(true),
	)
	c1.Limit(&colly.LimitRule{
		Parallelism: 16,
	})
	c2 := c1.Clone()
	c2.SetRequestTimeout(120 * time.Second)

	defer func() {
		c1.Wait()
		c2.Wait()
		fmt.Println("all done")
	}()

	c1.OnResponse(func(r *colly.Response) {
		var heroes []Hero
		err := json.Unmarshal(r.Body, &heroes)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, hero := range heroes {
			heroDir := path.Join(LocalDir, hero.DirName())
			ensureExists(heroDir)

			skins := hero.GetSkins()
			for _, skin := range skins {
				skinPath := path.Join(heroDir, skin.FileName())
				skinUrl := getSkinUrl(hero.EName, skin.Idx)
				ctx := colly.NewContext()
				ctx.Put("path", skinPath)
				c2.Request("GET", skinUrl, nil, ctx, nil)
			}
		}
	})

	c2.OnResponse(func(r *colly.Response) {
		skinPath := r.Ctx.Get("path")
		err := os.WriteFile(skinPath, r.Body, 0o644)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("done", skinPath)
	})

	c2.OnError(func(r *colly.Response, err error) {
		fmt.Println(err, r.Request.URL)
	})

	err := c1.Visit(ApiUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
}
