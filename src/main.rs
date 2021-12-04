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
        .route("/result", get(get_result))
        .route("/all_code", get(get_all_code));

    axum::Server::bind(&"0.0.0.0:8000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn get_result(code: Query<Code>) -> Json<Value> {
    if let Ok(ser) = Service::new().await {
        if code.level == LEVELDAY || code.level == LEVEL15M {
            if let Ok(_) = ser.calc(&code.level, &code.code).await {
                return Json(json!({ "code":0,"msg":"ok","data": "yes" }));
            }
        }
    }
    Json(json!({"code":-1,"msg":"fuck you err"}))
}

async fn get_all_code() {}
