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
    let mut cursor = sqlx::query("SELECT DISTINCT(code) FROM day").fetch(&pool);
    while let Some(row) = cursor.try_next().await? {
        let str: String = row.get("code");
        let mut cursor = sqlx::query("SELECT * from 15m WHERE code = ?")
            .bind(str.clone())
            .fetch(&pool);
        let mut arr_data = Vec::new();
        let mut index = 0;
        while let Some(row) = cursor.try_next().await? {
            let data = base::MetaData {
                index: index,
                high: row.get("high"),
                low: row.get("low"),
                open: row.get("open"),
                close: row.get("close"),
                timestamp: row.get("timestamp"),
            };
            arr_data.push(data);
            index += 1
        }

        let (list_fx, last_k) = fx::calc(&arr_data);
        let mut zous = base::ZouS {
            name: str.clone(),
            code: str.clone(),
            calc_at: 1,
            calc_type: base::CalcType::D,
            souce_data: arr_data,
            status: None,
            line: HashMap::new(),
            zs: HashMap::new(),
        };
        if list_fx.len() == 0 {
            continue;
        }
        zous.calc(&list_fx, last_k);
        let mut open = false;
        //let mut list = Vec::new();
        if let Some(_) = zous.status {
            open = true
        } else {
            continue;
        }
        if open {
            zous.save();
            break;
        }
    }

    Ok(())
}
