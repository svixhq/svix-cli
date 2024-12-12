use anyhow::Result;
use figment::{
    providers::{Env, Format, Toml},
    Figment,
};
use serde::{Deserialize, Serialize};
use std::io::Write;
use std::os::unix::fs::OpenOptionsExt;
use std::path::PathBuf;

#[derive(Debug, Default, Deserialize, Serialize)]
pub struct Config {
    pub auth_token: Option<String>,
    pub server_url: Option<String>,
}

impl Config {
    pub fn load() -> Result<Config> {
        let cfg_file = get_folder()?.join(FILE_NAME);
        let config: Config = Figment::new()
            .merge(Toml::file(cfg_file))
            .merge(Env::prefixed("SVIX_"))
            .extract()?;
        Ok(config)
    }
}

const FILE_NAME: &str = "config.toml";
const FILE_MODE: u32 = 0o600;

fn get_folder() -> Result<PathBuf> {
    let config_path = if cfg!(windows) {
        std::env::var("APPDATA")
    } else {
        std::env::var("XDG_CONFIG_HOME")
    };

    let pb = match config_path {
        Ok(path) => PathBuf::from(path),
        Err(_e) => {
            // FIXME: home_dir() can give incorrect results on Windows. Docs recommend "use a crate instead"
            #[allow(deprecated)]
            std::env::home_dir().ok_or_else(|| anyhow::anyhow!("unable to find config path"))?
        }
    };
    Ok(pb.join(".config").join("svix"))
}

fn write(settings: Config) -> Result<()> {
    let cfg_path = get_folder()?;
    let mut fh = std::fs::OpenOptions::new()
        .create(true)
        .truncate(true)
        .write(true)
        .mode(FILE_MODE)
        .open(cfg_path)?;

    let source = &toml::to_string_pretty(&settings)?;
    fh.write(source.as_bytes())?;
    Ok(())
}
