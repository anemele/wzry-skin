import { join } from "@std/path";
import { exists } from "@std/fs/exists";

const API_HEROLIST = "https://pvp.qq.com/web201605/js/herolist.json";
const API_SKINFILE = function (ename: number, skin_id: number): string {
  return `https://game.gtimg.cn/images/yxzj/img201606/heroimg/${ename}/${ename}-bigskin-${skin_id}.jpg`;
};
const LOCAL_SAVE_PATH = "wzry-skin-ts";
// if not exist, create it
if (!await exists(LOCAL_SAVE_PATH)) {
  await Deno.mkdir(LOCAL_SAVE_PATH);
}
await Deno.writeTextFile(join(LOCAL_SAVE_PATH, ".gitignore"), "*");

interface Hero {
  cname: string;
  ename: number;
  title: string;
  skin_name: string;
}

const tasks: Promise<void>[] = [];

async function downloadSkin(url: string, path: string) {
  const resp = await fetch(url);
  const arrayBuffer = await resp.arrayBuffer();
  const uint8Array = new Uint8Array(arrayBuffer);
  await Deno.writeFile(path, uint8Array);
}

async function handleHero(hero: Hero) {
  const heroPath = join(LOCAL_SAVE_PATH, `${hero.cname}_${hero.title}`);
  if (!await exists(heroPath)) {
    await Deno.mkdir(heroPath);
  }

  hero.skin_name.split("|").forEach((skin_name, index) => {
    const skinPath = join(heroPath, `${index + 1}_${skin_name}.jpg`);
    const url = API_SKINFILE(hero.ename, index + 1);
    tasks.push(downloadSkin(url, skinPath));
  });
}

const resp = await fetch(API_HEROLIST);
const herolist: Hero[] = await resp.json();

herolist.forEach((hero) => {
  handleHero(hero);
});

await Promise.all(tasks);
