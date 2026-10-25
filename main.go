package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

const (
	HeroListURL = "https://pvp.qq.com/web201605/js/herolist.json"
	LocalDir    = "skin-dirs"
	MaxTask     = 1 << 8
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

func getSkinURL(ename int, idx int) string {
	return fmt.Sprintf(
		"https://game.gtimg.cn/images/yxzj/img201606/heroimg/%d/%d-bigskin-%d.jpg",
		ename, ename, idx,
	)
}

var client http.Client

type Task struct {
	Path string
	URL  string
}

var taskChan = make(chan Task, 1)

func downloadSkin(task Task) (err error) {
	skinPath := task.Path
	skinURL := task.URL

	resp, err := client.Get(skinURL)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s %s", resp.Status, skinURL)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = os.WriteFile(skinPath, buf, 0o644)
	if err != nil {
		return
	}

	fmt.Fprintln(os.Stderr, "Downloaded", skinPath)
	return nil
}

func downloadHero(hero Hero) {
	heroDir := path.Join(LocalDir, hero.DirName())
	ensureExists(heroDir)

	skins := hero.GetSkins()
	for _, skin := range skins {
		taskChan <- Task{
			Path: path.Join(heroDir, skin.FileName()),
			URL:  getSkinURL(hero.EName, skin.Idx),
		}
	}
}

func run() (err error) {
	resp, err := client.Get(HeroListURL)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}
	defer resp.Body.Close()

	var heroes []Hero
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&heroes)
	if err != nil {
		return
	}

	ensureExists(LocalDir)

	var wg0 sync.WaitGroup
	for range MaxTask {
		wg0.Go(func() {
			for task := range taskChan {
				err := downloadSkin(task)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		})
	}

	var wg1 sync.WaitGroup
	for _, hero := range heroes {
		wg1.Go(func() {
			downloadHero(hero)
		})
	}

	wg1.Wait()
	close(taskChan)
	wg0.Wait()

	fmt.Fprintln(os.Stderr, "All done.")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
