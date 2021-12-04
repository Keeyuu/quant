use crate::base::*;
use crate::fx::*;
use crate::*;
use anyhow::Result;
use futures::TryStreamExt;
use sqlx::mysql::MySqlPoolOptions;
use sqlx::{Executor, MySql, Pool, Row};
use std::collections::HashMap;

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
        let sql = format!("SELECT * from {} WHERE code = ?", table);
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
    pub async fn calc(&self, table: &str, code: &str) -> Result<()> {
        let arr = self.query_data_by_code(table, code).await?;
        let (list_fx, last_k) = calc(&arr);
        let mut zous = ZouS::new(code.to_string(), CalcType::from(table), arr);
        zous.calc(&list_fx, last_k);
        zous.to_rsb().save()?;
        Ok(())
    }
}
