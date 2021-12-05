#![feature(linked_list_remove)]
#![feature(exclusive_range_pattern)]
#![feature(derive_default_enum)]
use axum::{extract::Query, routing::get, Json, Router};
mod base;
mod fx;
mod service;
mod transmission;
use serde_json::{json, Value};
use service::Service;
use transmission::*;
#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", get(|| async { "Hello, World!" }))
        .route("/api/result", get(get_result))
        .route("/api/all_code", get(get_all_code))
        .route("/api/calc_all", get(calc_all));

    axum::Server::bind(&"0.0.0.0:8000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn get_result(code: Query<Code>) -> Json<Value> {
    println!("get_result get request {:?}", code);
    if let Ok(ser) = Service::new().await {
        if code.level == LEVELDAY || code.level == LEVEL15M {
            if let Ok(_) = ser.calc(&code.level, &code.code, code.index).await {
                return Json(json!({ "code":0,"msg":"ok","data": "yes" }));
            }
        }
    }
    Json(json!({"code":-1,"msg":"fuck you err"}))
}

async fn get_all_code() {}
async fn calc_all() -> &'static str {
    println!("calc_all get request");
    if let Ok(ser) = Service::new().await {
        //ser.calc_all(LEVELDAY).await;
        ser.calc_all(LEVEL15M).await;
        return "ok";
    }
    "err 1"
}
