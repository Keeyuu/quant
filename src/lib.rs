#![feature(linked_list_remove)]
#![feature(exclusive_range_pattern)]
#![feature(derive_default_enum)]
mod base;
mod fx;
use libc::size_t;
use std::slice;
#[no_mangle]
pub extern "C" fn calc(len: size_t, l: *const f32, h: *const f32, c: *const f32, o: *const f32) {
    println!("calc start");
    let l = unsafe {
        assert!(!l.is_null());

        slice::from_raw_parts(l, len as usize)
    };
    let h = unsafe {
        assert!(!h.is_null());

        slice::from_raw_parts(h, len as usize)
    };
    let c = unsafe {
        assert!(!c.is_null());

        slice::from_raw_parts(c, len as usize)
    };
    let o = unsafe {
        assert!(!o.is_null());

        slice::from_raw_parts(o, len as usize)
    };
    let numbers = base::MetaData::new_arr(len as i64, l, h, c, o);
    println!("has get data {:?}", numbers);
}

#[no_mangle]
pub extern "C" fn say_hello() {
    println!("hello");
}
