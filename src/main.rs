#![feature(linked_list_remove)]
#![feature(exclusive_range_pattern)]
#![feature(derive_default_enum)]
#[macro_use]
extern crate rocket;
use rocket::serde::json::{json, Json, Value};
mod base;
use std::collections::HashMap;
#[get("/?<code>")]
fn get_zous(code: Option<String>) -> Option<Json<base::ZouSOut>> {
    if let Some(code) = code {
        return Some(Json(base::ZouSOut {
            code: code,
            souce_data: Vec::new(),
            status: Some(vec![(1, "buy".to_string())]),
            line: HashMap::new(),
            zs: HashMap::new(),
        }));
    }
    None
}
#[get("/all_code")]
fn get_all_code() -> Value {
    json!({ "a": vec!["1", "2", "3"] })
}

#[launch]
fn rocket() -> _ {
    rocket::build()
        .mount("/", routes![get_zous])
        .mount("/", routes![get_all_code])
        .register("/", catchers![not_found])
}

//use std::collections::HashMap;

//use sqlx::mysql::MySqlPoolOptions;
//use sqlx::Row;
//mod base;
//use futures::TryStreamExt;
//mod fx;

//let pool = MySqlPoolOptions::new()
//        .max_connections(5)
//        .connect("mysql://quant:Quant.8.cc@localhost/quant")
//        .await?;
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

#[catch(404)]
fn not_found() -> Value {
    json!({
        "status": "error",
        "reason": "Resource was not found."
    })
}
