import asyncio
from dataclasses import dataclass
from itertools import chain, starmap
from pathlib import Path

import aiofiles
import aiohttp
import orjson
from fake_useragent import FakeUserAgent
from mashumaro.mixins.orjson import DataClassORJSONMixin

local_save_path = Path("wzry-skin-dirs")
local_save_path.mkdir(exist_ok=True)

API_BASEURL = "https://pvp.qq.com/"
API_HEROLIST = "https://pvp.qq.com/web201605/js/herolist.json"
API_SKINFILE = "https://game.gtimg.cn/images/yxzj/img201606/heroimg/{ename}/{ename}-bigskin-{skin_id}.jpg"


@dataclass
class Hero(DataClassORJSONMixin):
    cname: str
    ename: int
    title: str
    id_name: str
    skin_name: str

    @property
    def to_str(self):
        return f"{self.cname}_{self.title}, ({self.ename} {self.id_name}), {self.skin_name}"

    @property
    def dirname(self):
        return f"{self.cname}_{self.title}"

    def get_skins(self):
        names = self.skin_name.split("|")
        it = enumerate(names, 1)
        return starmap(Skin, it)


@dataclass
class Skin:
    skin_id: int
    name: str

    @property
    def filename(self):
        return f"{self.skin_id}_{self.name}.jpg"


async def download_and_write(client: aiohttp.ClientSession, url: str, path: Path):
    async with client.get(url) as resp:
        # TODO: response status?
        content = await resp.read()

    async with aiofiles.open(path, "wb") as fp:
        await fp.write(content)


def handle_hero(client: aiohttp.ClientSession, hero: Hero):
    hero_path = local_save_path.joinpath(hero.dirname)
    hero_path.mkdir(exist_ok=True)

    tasks = []
    for skin in hero.get_skins():
        skin_path = hero_path.joinpath(skin.filename)

        # comment the following to re-download
        if skin_path.exists():
            continue

        skin_url = API_SKINFILE.format(ename=hero.ename, skin_id=skin.skin_id)
        task = download_and_write(client, skin_url, skin_path)
        tasks.append(task)

    return tasks


async def main():
    async with aiohttp.ClientSession(
        # NOTE: don't set this True, or it will crash
        # raise_for_status=True,
        base_url=API_BASEURL,
        headers={"User-Agent": str(FakeUserAgent.random)},
    ) as client:
        async with client.get(API_HEROLIST) as first_resp:
            content = await first_resp.read()
        data = orjson.loads(content)
        herolist = map(Hero.from_dict, data)
        tasks = (handle_hero(client, hero) for hero in herolist)
        tasks = chain.from_iterable(tasks)

        await asyncio.gather(*tasks)


if __name__ == "__main__":
    asyncio.run(main())
