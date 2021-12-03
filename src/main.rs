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
        let mut cursor = sqlx::query("SELECT * from 15m WHERE code = ?")
            .bind(str.clone())
            .fetch(&pool);
        let mut arr_data = Vec::new();
        let mut arr_data_ = Vec::new();
        let mut index = 0;
        while let Some(row) = cursor.try_next().await? {
            let org = base::MetaData_ {
                index: index,
                high: row.get("high"),
                low: row.get("low"),
                open: row.get("open"),
                close: row.get("close"),
                times: row.get("date"),
            };
            arr_data_.push(org.clone());
            let data = base::MetaData {
                index: org.index,
                high: org.high,
                low: org.low,
                open: org.open,
                close: org.close,
                timestamp: org.times.timestamp_millis(),
            };
            arr_data.push(data);
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
        if list_fx.len() == 0 {
            continue;
        }
        zous.calc(&list_fx, last_k);
        //let mut list = Vec::new();
        if let Some(data) = zous.status {
            let a = base::ZouS_ {
                name: zous.name,
                code: zous.code,
                calc_at: zous.calc_at,
                calc_type: zous.calc_type,
                souce_data: arr_data_,
                status: None,
                line: zous.line,
                zs: zous.zs,
            };
            println!("{:?}", a);
            return Ok(());
        }
    }

    Ok(())
}
