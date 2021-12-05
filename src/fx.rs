use std::collections::LinkedList;

use crate::base::*;
use base::*;

pub fn calc(arr: &Vec<MetaData>) -> (LinkedList<MetaFX>, PureK) {
    if arr.len() == 0 {
        return (LinkedList::new(), PureK::default());
    }
    let a5 = MetaData::get_a5_array(&arr);
    let arr_k = contain(arr, &a5);
    let last_k = arr_k.last().unwrap().clone();
    let list_fx = fx(arr_k);
    (remove_invalid_fx(list_fx), last_k)
}
pub fn view(list: &Vec<MetaData>) -> Vec<MetaDataFX_> {
    let (list_, _) = calc(list);
    conversion(list_)
}
fn conversion(list_fx: LinkedList<MetaFX>) -> Vec<MetaDataFX_> {
    //*转化为输出结构
    let mut arr_fx = Vec::new();
    for i in list_fx {
        arr_fx.push(i.conversion())
    }
    arr_fx
}
fn contain(arr: &Vec<MetaData>, a5: &Vec<f64>) -> Vec<PureK> {
    let mut arr_k = Vec::with_capacity(arr.len() / 2);
    let mut i = 1;
    let mut this = arr[0].conversion();
    while i < arr.len() {
        this = do_contain(a5, &mut arr_k, this, arr[i].conversion(), i);
        i += 1;
    }
    arr_k.push(this);
    arr_k
}

fn do_contain<'a>(
    a5: &Vec<f64>,
    arr_k: &'a mut Vec<PureK>,
    mut last: PureK,
    this: PureK,
    i: usize,
) -> PureK {
    if last.item.contain(&this.item) {
        //last.index_range.1 = i as i64;
        if last.item.close >= a5[i - 1] {
            last.item.close = bigger(last.item.close, this.item.close);
            last.item.open = bigger(last.item.open, this.item.open);
            last.item.low = bigger(last.item.low, this.item.low);
        } else {
            last.item.close = smaller(last.item.close, this.item.close);
            last.item.open = smaller(last.item.open, this.item.open);
            last.item.high = smaller(last.item.high, this.item.high);
        }
        return last;
    } else {
        arr_k.push(last);
        return this;
    }
}
fn fx(arr_k: Vec<PureK>) -> LinkedList<MetaFX> {
    let mut list_fx = LinkedList::new();
    let mut i = 1;
    while i < arr_k.len() - 1 {
        if let Some(fx) = arr_k[i].check_di(&arr_k[i - 1], &arr_k[i + 1]) {
            list_fx.push_back(fx)
        }
        if let Some(fx) = arr_k[i].check_din(&arr_k[i - 1], &arr_k[i + 1]) {
            list_fx.push_back(fx)
        }
        i += 1;
    }
    list_fx
}
fn check_reasonable_fx(list_fx: &LinkedList<MetaFX>) -> bool {
    if list_fx.len() <= 1 {
        return true;
    }
    let mut iter_ = list_fx.iter();
    let mut this = iter_.next().unwrap();
    loop {
        if let Some(fx) = iter_.next() {
            if !this.valid(fx) {
                return false;
            } else if let Some(fx) = iter_.next() {
                this = fx
            } else {
                return true;
            }
        } else {
            return true;
        }
    }
}
fn remove_invalid_fx(mut list_fx: LinkedList<MetaFX>) -> LinkedList<MetaFX> {
    let mut ok = false;
    while !ok || !check_reasonable_fx(&mut list_fx) {
        ok = do_remove_invalid_fx(&mut list_fx);
    }
    list_fx
}
fn do_remove_invalid_fx(list_fx: &mut LinkedList<MetaFX>) -> bool {
    let mut iter_ = list_fx.iter();
    let len_ = iter_.len();
    let mut this;
    if let Some(fx) = iter_.next() {
        this = fx.clone()
    } else {
        return true;
    }
    loop {
        if let Some(next) = iter_.next() {
            let this_len = iter_.len() + 1;
            if this.fx_type == next.fx_type {
                if this.valid(next) {
                    list_fx.remove(len_ - this_len);
                    return false;
                }
                list_fx.remove(len_ - this_len - 1);
                return false;
            } else {
                if this.valid(next) {
                    if this.valid_time(next) {
                        this = next.clone();
                        continue;
                    } else {
                        list_fx.remove(len_ - this_len);
                        return false;
                    }
                }
                list_fx.remove(len_ - this_len - 1);
                return false;
            }
        } else {
            return true;
        }
    }
}
//*-------------------------分割--------------------------------------------------------
//

//*-------------------------分割--------------------------------------------------------
pub mod base {
    pub fn bigger<T: std::cmp::PartialOrd>(a: T, b: T) -> T {
        if a >= b {
            return a;
        }
        b
    }
    pub fn smaller<T: std::cmp::PartialOrd>(a: T, b: T) -> T {
        if a <= b {
            return a;
        }
        b
    }
}
//*-------------------------分割--------------------------------------------------------
#[cfg(test)]
mod tests {
    use super::{base::*, *};
    use crate::base::*;
    use rand::Rng;
    use std::collections::LinkedList;
    use std::fs::{self, File, OpenOptions};
    use std::io::prelude::*;
    fn build_test_data(size: i64, is_show: bool) -> Vec<MetaData> {
        let mut rand_number = rand::thread_rng();
        let mut org = Vec::with_capacity(size as usize);
        for i in 0..size {
            let mut item = MetaData::default();
            item.index = i;
            item.close = rand_number.gen_range(0.1..9.9);
            item.high = rand_number.gen_range(0.1..9.9);
            item.low = rand_number.gen_range(0.1..9.9);
            org.push(item)
        }
        if is_show {
            println!("{:?}\nbuild_test_data len: {}", org, org.len())
        }
        org
    }
    #[test]
    fn test_a5() {
        let tmp = MetaData::get_a5_array(&build_test_data(10, true));
        println!("{:?}", tmp)
    }
    #[test]
    fn test_max() {
        assert_eq!(2.1, bigger(1.1, 2.1),);
    }
    #[test]
    fn test_list_len() {
        let mut list = LinkedList::new();
        for i in 0..9 {
            list.push_back(i)
        }
        let mut iter_ = list.iter();
        let len_ = iter_.len();
        println!("{}", len_ - iter_.len());
        iter_.next();
        println!("{}", len_ - iter_.len());
        iter_.next();
        println!("{:?},{}", list, list.len());
        list.front().unwrap();
        println!("{:?},{}", list, list.len());
    }
    #[test]
    fn test_vec_len() {
        let mut arr = Vec::new();
        for i in 0..2 {
            arr.push(i);
        }
        arr.clear();
        let mut arr = arr.iter();
        println!("{}", arr.len());
        arr.next();
        println!("{}", arr.len());
        arr.next();
    }
    #[test]
    //fn test_calc_fx() {
    //    calc(build_test_data(100, false));
    //}
    #[test]
    fn test_calc_fx_new() {}
    #[test]
    fn test_range_slice() {
        let a = vec![0, 1, 2, 3, 4];
        println!("{:?}", &a[0..2 + 1]);
    }
}
