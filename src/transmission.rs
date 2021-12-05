use crate::base::*;
use anyhow::Result;
use serde::*;
use std::collections::HashMap;
use std::fs::{File, OpenOptions};
use std::io::Write;
pub const LEVELDAY: &str = "day";
pub const LEVEL15M: &str = "15m";

#[derive(Deserialize, Debug)]
pub struct Code {
    pub code: String,
    pub level: String,
    pub index: i64,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct RsbZouS {
    pub index: i64,
    pub code: String,
    pub souce_data: Vec<MetaData>,
    pub line: HashMap<usize, Vec<Line>>,
    pub zs: HashMap<usize, Vec<Zs>>,
}
impl RsbZouS {
    pub fn save(&self) -> Result<()> {
        let s = serde_json::to_string(&self).unwrap();
        let mut f = OpenOptions::new()
            .write(true)
            .truncate(true)
            .create(true)
            .open("./view/src/views/demo.json")?;
        f.write_all(s.as_bytes())?;
        Ok(())
    }
}

impl ZouS {
    pub fn to_rsb(self, index: i64) -> RsbZouS {
        RsbZouS {
            index: index,
            code: self.code,
            souce_data: self.souce_data,
            line: self.line,
            zs: self.zs,
        }
    }
}

impl CalcType {
    pub fn from(s: &str) -> Self {
        match s {
            LEVELDAY => Self::D,
            LEVEL15M => Self::Min15,
            _ => Self::None,
        }
    }
}
