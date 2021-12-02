#![feature(linked_list_remove)]
#![feature(exclusive_range_pattern)]
#![feature(derive_default_enum)]
use std::collections::HashMap;

use sqlx::mysql::MySqlPoolOptions;
use sqlx::Row;
mod base;
use futures::TryStreamExt;
mod fx;

#[tokio::main]
async fn main() -> Result<(), sqlx::Error> {
    let pool = MySqlPoolOptions::new()
        .max_connections(5)
        .connect("mysql://quant:Quant.8.cc@localhost/quant")
        .await?;
    let mut youxi: Vec<(String, Vec<(usize, base::Status)>)> = Vec::new();
    let mut cursor = sqlx::query("SELECT DISTINCT(code) FROM day").fetch(&pool);
    while let Some(row) = cursor.try_next().await? {
        let str: String = row.get("code");
        let mut cursor = sqlx::query("SELECT * from day WHERE code = ?")
            .bind(str.clone())
            .fetch(&pool);
        let mut arr_data = Vec::new();
        let mut index = 0;
        while let Some(row) = cursor.try_next().await? {
            arr_data.push(base::MetaData {
                index: index,
                high: row.get("high"),
                low: row.get("low"),
                open: row.get("open"),
                close: row.get("close"),
            });
            index += 1
        }

        let (list_fx, last_k) = fx::calc(arr_data);
        let mut zous = base::ZouS {
            name: str.clone(),
            code: str.clone(),
            calc_at: 1,
            calc_type: base::CalcType::D,
            souce_data: Vec::new(),
            status: None,
            line: HashMap::new(),
            zs: HashMap::new(),
        };
        zous.calc(&list_fx, last_k);
        if let Some(data) = zous.status {
            youxi.push((str, data));
        }
        println!("you xi len {}", youxi.len());
    }

    println!("\n\n\n\n {:?}", youxi);
    Ok(())
}
