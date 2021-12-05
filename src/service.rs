use crate::base::*;
use crate::fx::*;
use crate::*;
use anyhow::Result;
use futures::TryStreamExt;
use sqlx::mysql::MySqlPoolOptions;
use sqlx::{Executor, MySql, Pool, Row};
use std::collections::HashMap;
use std::fs::{File, OpenOptions};
use std::io::Write;
pub struct Service {
    pool: Pool<MySql>,
}

impl Service {
    pub async fn new() -> Result<Self> {
        let pool = MySqlPoolOptions::new()
            .max_connections(5)
            .connect("mysql://quant:Quant.8.cc@localhost/quant")
            .await?;
        Ok(Self { pool: pool })
    }
    async fn query_data_by_code(&self, table: &str, code: &str) -> Result<Vec<MetaData>> {
        let sql = format!(
            "SELECT high,low,open,close,timestamp,volume from {} WHERE code = ?",
            table
        );
        let mut cursor = sqlx::query(&sql).bind(code).fetch(&self.pool);
        let mut arr_data = Vec::new();
        let mut index = 0;
        while let Some(row) = cursor.try_next().await? {
            let data = MetaData {
                index: index,
                high: row.get("high"),
                low: row.get("low"),
                open: row.get("open"),
                close: row.get("close"),
                timestamp: row.get("timestamp"),
                volume: row.get("volume"),
            };
            arr_data.push(data);
            index += 1
        }
        Ok(arr_data)
    }
    pub async fn calc(&self, table: &str, code: &str, index: i64) -> Result<()> {
        let arr = self.query_data_by_code(table, code).await?;
        let (list_fx, last_k) = calc(&arr);
        let mut zous = ZouS::new(code.to_string(), CalcType::from(table), arr);
        zous.calc(&list_fx, last_k);
        zous.to_rsb(index).save()?;
        Ok(())
    }
    pub async fn calc_no_save(&self, table: &str, code: &str) -> Result<ZouS> {
        let arr = self.query_data_by_code(table, code).await?;
        let (list_fx, last_k) = calc(&arr);
        let mut zous = ZouS::new(code.to_string(), CalcType::from(table), arr);
        zous.calc(&list_fx, last_k);
        Ok(zous)
    }
    pub async fn calc_all(&self, table: &str) -> Result<()> {
        let mut codes = Vec::new();
        let mut sum = 0;
        for code in self.get_all_code().await? {
            sum += 1;
            let zs = self.calc_no_save(table, &code).await?;
            if let Some(arr) = zs.status {
                if arr.len() > 0 {
                    println!("yes {} len :{} sum :{}", table, codes.len(), sum);
                    codes.push(zs.code)
                }
            }
        }
        let s = serde_json::to_string(&codes)?;
        let mut f = OpenOptions::new()
            .write(true)
            .truncate(true)
            .create(true)
            .open(format!("Result{}.json", table))?;
        f.write_all(s.as_bytes())?;
        println!("calc all {} fin", table);
        Ok(())
    }
    pub async fn get_all_code(&self) -> Result<Vec<String>> {
        let mut cursor =
            sqlx::query("SELECT code FROM code where end_date>=NOW()").fetch(&self.pool);
        let mut arr_data = Vec::new();
        while let Some(row) = cursor.try_next().await? {
            arr_data.push(row.get("code"));
        }
        Ok(arr_data)
    }
}
