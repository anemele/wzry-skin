package main

import (
	"fmt"
	"strings"
)

type Hero struct {
	CName    string `json:"cname"`
	EName    int    `json:"ename"`
	Title    string `json:"title"`
	IDName   string `json:"id_name"`
	SkinName string `json:"skin_name"`
}

func (hero Hero) String() string {
	return fmt.Sprintf(
		"%s_%s, (%d %s), %s",
		hero.CName, hero.Title, hero.EName, hero.IDName, hero.SkinName,
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

func (hero Hero) GetSkins() (skins []Skin) {
	for i, name := range strings.Split(hero.SkinName, "|") {
		skins = append(skins, Skin{Idx: i + 1, Name: name})
	}
	return
}
