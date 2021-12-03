//use std::collections::HashMap;
use sqlx::{
    mysql::{MySql, MySqlPoolOptions},
    types::chrono::{NaiveDate, NaiveDateTime},
    Pool,
};

use std::time::SystemTime;
//use sqlx::Row;
//mod base;
use anyhow::Result;
//use futures::TryStreamExt;
//mod fx;

pub struct Mysql {
    pool: Pool<MySql>,
}
#[derive(sqlx::FromRow, Debug)]
pub struct Code {
    pub code: String,
    pub name: String,
    pub type_: String,
}

impl Mysql {
    pub async fn new() -> Result<Self> {
        let pool = MySqlPoolOptions::new()
            .connect("mysql://quant:Quant.8.cc@localhost/quant")
            .await?;
        Ok(Self { pool: pool })
    }
    pub async fn get_all_code(&self) -> Result<Vec<Code>> {
        let stream = sqlx::query_as::<_, Code>("SELECT * FROM code where end_date>=NOW()")
            .fetch_all(&self.pool)
            .await?;
        Ok(stream)
    }
    pub async fn get_happy_code(&self) -> Result<Vec<Code>> {
        let stream = sqlx::query_as::<_, Code>("SELECT * FROM code where end_date>=NOW()")
            .fetch_all(&self.pool)
            .await?;
        Ok(stream)
    }
}

//    let mut youxi: Vec<(String, Vec<(usize, base::Status)>)> = Vec::new();
//    let mut cursor = sqlx::query("SELECT DISTINCT(code) FROM day").fetch(&pool);
//    while let Some(row) = cursor.try_next().await? {
//        let str: String = row.get("code");
//        let mut cursor = sqlx::query("SELECT * from day WHERE code = ?")
//            .bind(str.clone())
//            .fetch(&pool);
//        let mut arr_data = Vec::new();
//        let mut index = 0;
//        while let Some(row) = cursor.try_next().await? {
//            arr_data.push(base::MetaData {
//                index: index,
//                high: row.get("high"),
//                low: row.get("low"),
//                open: row.get("open"),
//                close: row.get("close"),
//            });
//            index += 1
//        }

//        let (list_fx, last_k) = fx::calc(arr_data);
//        let mut zous = base::ZouS {
//            name: str.clone(),
//            code: str.clone(),
//            calc_at: 1,
//            calc_type: base::CalcType::D,
//            souce_data: Vec::new(),
//            status: None,
//            line: HashMap::new(),
//            zs: HashMap::new(),
//        };
//        zous.calc(&list_fx, last_k);
//        if let Some(data) = zous.status {
//            youxi.push((str, data));
//        }
//        println!("you xi len {}", youxi.len());
//    }

//    println!("\n\n\n\n {:?}", youxi);
