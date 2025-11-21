use std::{
    fs,
    path::{Path, PathBuf},
};

use anyhow::Result;
use futures::future::join_all;
use reqwest::Client;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
struct Hero {
    cname: String,
    ename: u16,
    title: String,
    skin_name: String,
}

const LOCAL_SAVE_PATH: &str = "wzry-skin-rs";

const API_BASEURL: &str = "https://pvp.qq.com/";
const API_HEROLIST: &str = "https://pvp.qq.com/web201605/js/herolist.json";
// const API_SKINFILE: &str =
//     "https://game.gtimg.cn/images/yxzj/img201606/heroimg/{ename}/{ename}-bigskin-{skin_id}.jpg";
#[inline]
fn get_skin_url(ename: u16, skin_id: usize) -> String {
    format!(
        "https://game.gtimg.cn/images/yxzj/img201606/heroimg/{ename}/{ename}-bigskin-{skin_id}.jpg"
    )
}

async fn download_skin(client: &Client, url: String, path: PathBuf) -> Result<()> {
    let data = client.get(url).send().await?.bytes().await?;
    fs::write(path, data)?;
    Ok(())
}

async fn handle_hero(
    client: &Client,
    hero: Hero,
    root_path: &Path,
) -> Result<Vec<impl Future<Output = Result<()>>>> {
    let hero_path = root_path.join(format!("{}_{}", hero.cname, hero.title));
    if !hero_path.exists() {
        fs::create_dir(&hero_path)?;
    }

    let tasks = hero
        .skin_name
        .split('|')
        .enumerate()
        .map(|(idx, skin_name)| {
            let idx = idx + 1;
            let skin_path = hero_path.join(format!("{idx}_{skin_name}.jpg"));
            let url = get_skin_url(hero.ename, idx);
            download_skin(client, url, skin_path.to_path_buf())
        })
        .collect();

    Ok(tasks)
}

#[tokio::main]
async fn main() -> Result<()> {
    let root_path = Path::new(LOCAL_SAVE_PATH);
    if !root_path.exists() {
        fs::create_dir(root_path)?;
    }
    fs::write(root_path.join(".gitignore"), b"*")?;

    let client = Client::new();
    if !client.get(API_BASEURL).send().await?.status().is_success() {
        anyhow::bail!("failed to get base url");
    }

    let herolist = client
        .get(API_HEROLIST)
        .send()
        .await?
        .json::<Vec<Hero>>()
        .await?;

    let mut tasks = vec![];

    for hero in herolist {
        let task = handle_hero(&client, hero, root_path).await?;
        tasks.extend(task);
    }
    join_all(tasks).await;

    Ok(())
}
