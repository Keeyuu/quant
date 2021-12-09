#![feature(linked_list_remove)]
#![feature(exclusive_range_pattern)]
#![feature(derive_default_enum)]
mod base;
mod fx;
use libc::size_t;
use std::slice;
#[no_mangle]
pub extern "C" fn calc(len: size_t, n: *const f32) {
    println!("calc start");
    let numbers = unsafe {
        assert!(!n.is_null());

        slice::from_raw_parts(n, len as usize)
    };
    //base::MetaData::new_arr(len_, l, h, c, o);
    println!("has get data {:?}", numbers);
}

#[no_mangle]
pub extern "C" fn say_hello() {
    println!("hello");
}
