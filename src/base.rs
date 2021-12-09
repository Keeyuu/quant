use std::collections::{HashMap, LinkedList};
pub const DIN: &str = "Din";
pub const DI: &str = "Di";
pub const INIT: &str = "Init";
pub const UP: &str = "Up";
pub const DOWN: &str = "Down";
// TODO 调整至最佳比例
const WEIGHT_S: f64 = 0.8;
const WEIGHT_O: f64 = 0.05;
const WEIGHT_C: f64 = 0.15;

#[derive(Clone, Copy, Default, Debug)]
pub struct MetaData {
    pub index: i64,
    pub high: f64,
    pub low: f64,
    pub open: f64,
    pub close: f64,
}

#[derive(Clone, Debug)]
pub struct MetaDataFX_ {
    pub item: PureK,
    pub fx_type: String,
}
#[derive(Clone, Copy, Debug, Default)]
pub struct PureK {
    pub index_range: (i64, i64),
    pub item: MetaData,
}
//*-------------------------对内--------------------------------------------------------
#[derive(Copy, Clone, Debug)]
pub struct MetaFX {
    pub item: PureK,
    pub fx_type: &'static str,
}
#[derive(Clone, Copy, Debug, PartialEq)]
pub enum Direction {
    Up,
    Down,
}

#[derive(Clone, Debug)]
pub struct Line {
    pub level: u8,
    pub index_range: (i64, i64),
    pub line_range: (f64, f64),
    pub difference: f64,
    pub direction: Direction,
    pub status: Status,
}

#[derive(Clone, Debug)]
pub struct Zs {
    pub level: usize,
    pub index_range: (i64, i64),
    pub center_range: (f64, f64),
    pub grow_center_range: (f64, f64),
    pub wave_range: (f64, f64),
    pub status: Status,
    pub lines: Vec<Line>,
}
#[derive(Clone, Debug)]
pub struct ZouS {
    pub code: String,
    pub calc_type: CalcType,
    pub souce_data: Vec<MetaData>,
    pub status: Option<Vec<(usize, Status)>>,
    pub line: HashMap<usize, Vec<Line>>,
    pub zs: HashMap<usize, Vec<Zs>>,
}

#[derive(Clone, Copy, Debug, PartialEq)]
pub enum Status {
    BuyTree,
    LowBuyTree,
    BuyTreeNear,
    LowBuyTreeNear,
    SellTree,
    Init,
    Grow,
    Down,
    QSUP,
    QSDOWN,
}
#[derive(Clone, Copy, Debug)]
pub enum CalcType {
    D,
    Min15,
    None,
}

//*-------------------------对内--------------------------------------------------------
//-------------------------方法-------------------------------
impl MetaData {
    // TODO 调整至最佳包含关系
    pub fn contain(&self, next: &Self) -> bool {
        self.high > next.high && self.low < next.low
    }
    pub fn conversion(self) -> PureK {
        PureK {
            item: self,
            index_range: (self.index, self.index),
        }
    }
    pub fn get_a5_array(arr: &Vec<MetaData>) -> Vec<f64> {
        return Self::get_an_array(arr, 5);
    }
    pub fn get_a20_array(arr: &Vec<MetaData>) -> Vec<f64> {
        return Self::get_an_array(arr, 20);
    }
    fn get_an_array(arr: &Vec<MetaData>, n: i64) -> Vec<f64> {
        let mut an = Vec::with_capacity(arr.len());
        for i in arr {
            if i.index <= n - 1 {
                an.push(
                    (arr[..=(i.index) as usize]
                        .iter()
                        .fold(0.0, |x, y| x + y.close))
                        / (i.index + 1) as f64,
                );
            } else {
                an.push(
                    (arr[(i.index - n + 1) as usize..=(i.index) as usize]
                        .iter()
                        .fold(0.0, |x, y| x + y.close))
                        / (n) as f64,
                );
            }
        }
        an
    }
    pub fn new_arr(len_: i64, l: &[f32], h: &[f32], c: &[f32], o: &[f32]) -> Vec<Self> {
        let mut arr = Vec::new();
        for _ in 0..len_ {
            let index = arr.len();
            arr.push(Self {
                index: index as i64,
                low: l[index] as f64,
                high: h[index] as f64,
                close: c[index] as f64,
                open: o[index] as f64,
            })
        }
        println!("new arr len :{}\n{:?}", arr.len(), arr);
        arr
    }
}
//-------------------------方法-------------------------------
impl PureK {
    pub fn check_di(self, last: &Self, next: &Self) -> Option<MetaFX> {
        if (self.item.low <= last.item.low && self.item.low <= next.item.low)
        //&& (self.calc_value_low() < last.calc_value_low()
        //&& self.calc_value_low() < next.calc_value_low())
        {
            return Some(MetaFX {
                item: self,
                fx_type: DI,
            });
        }
        None
    }
    pub fn check_din(self, last: &Self, next: &Self) -> Option<MetaFX> {
        if (self.item.high >= last.item.high && self.item.high >= next.item.high)
        //&& (self.calc_value_high() > last.calc_value_high()
        //&& self.calc_value_high() > next.calc_value_high())
        {
            return Some(MetaFX {
                item: self,
                fx_type: DIN,
            });
        }
        None
    }
    fn calc_value_low(&self) -> f64 {
        return self.item.close * WEIGHT_C + self.item.open * WEIGHT_O + self.item.low * WEIGHT_S;
    }
    fn calc_value_high(&self) -> f64 {
        return self.item.close * WEIGHT_C + self.item.open * WEIGHT_O + self.item.high * WEIGHT_S;
    }
    fn calc_value(&self) -> f64 {
        return self.item.close * WEIGHT_C
            + self.item.open * WEIGHT_O
            + (self.item.low + self.item.high) * WEIGHT_S / 2.;
    }
}
//-------------------------方法-------------------------------
impl MetaFX {
    pub fn valid(&self, next: &Self) -> bool {
        match self.fx_type {
            DIN => {
                return self.item.calc_value_high() >= next.item.calc_value_high();
            }
            DI => {
                return self.item.calc_value_low() <= next.item.calc_value_low();
            }
            _ => {
                println!("fx type err: {:?}", self.clone());
                return false;
            }
        }
    }
    pub fn valid_time(&self, next: &Self) -> bool {
        return next.item.index_range.0 >= self.item.index_range.1 + 3;
    }
    pub fn get_abs_from(&self, from: f64) -> f64 {
        return (self.item.item.close.abs() - from.abs()).abs();
    }
    pub fn conversion(self) -> MetaDataFX_ {
        MetaDataFX_ {
            item: self.item,
            fx_type: String::from(self.fx_type),
        }
    }
}
//-------------------------方法-------------------------------

//-------------------------方法-------------------------------
impl Line {
    fn calc_trait_value(arr_line: &[Line], i: usize, len: usize) -> f64 {
        arr_line[((1 + i) - len)..i + 1]
            .iter()
            .fold(0.0, |x, y| x + y.difference)
    }
    fn check_trait_2_2(arr_line: &[Line], i: usize) -> bool {
        Self::calc_trait_value(arr_line, i, 2) * arr_line[i - 2].difference >= 0.
    }
    fn check_trait_x_y(arr_line: &[Line], i: usize, x: usize, y: usize) -> bool {
        Self::calc_trait_value(arr_line, i, x) * arr_line[i - y].difference >= 0.
    }
    fn calc_intersection_maxmin(arr_line: &[Line]) -> ((f64, f64), (f64, f64), (i64, i64)) {
        let mut intersection = (f64::MIN, f64::MAX);
        let mut maxmin = (f64::MAX, f64::MIN);
        let mut index_range = (i64::MAX, i64::MIN);
        let mut index = 0;
        for i in arr_line {
            maxmin.0 = i.line_range.0.min(maxmin.0);
            maxmin.1 = i.line_range.1.max(maxmin.1);
            if index != 0 && index != arr_line.len() - 1 {
                index_range.1 = i.index_range.1.max(index_range.1);
            }
            if index < 3 {
                //只计算前三笔
                index_range.0 = i.index_range.0.min(index_range.0);
                intersection.0 = i.line_range.0.max(intersection.0);
                intersection.1 = i.line_range.1.min(intersection.1);
            }
            index += 1
        }
        (intersection, maxmin, index_range)
    }
}

//-------------------------方法-------------------------------
fn clone_from<T: Clone + Sized>(arr: &[T]) -> Vec<T> {
    let mut v = Vec::new();
    for i in arr {
        v.push(i.clone())
    }
    v
}
fn clone_from_<T: Clone + Sized>(arr: &[&T]) -> Vec<T> {
    let mut v = Vec::new();
    for i in arr {
        v.push(i.clone().clone())
    }
    v
}

impl ZouS {
    pub fn new(code: String, calc_type: CalcType, arr: Vec<MetaData>) -> Self {
        Self {
            code: code,
            calc_type: calc_type,
            souce_data: arr,
            status: None,
            line: HashMap::new(),
            zs: HashMap::new(),
        }
    }
    pub fn calc(&mut self, list_fx: &LinkedList<MetaFX>, last_k: PureK) {
        if list_fx.len() == 0 {
            return;
        }
        self.calc_level_zore(list_fx, last_k);
        self.calc_level_one();
        self.calc_level_any();
        self.calc_zs();
        self.analyze()
    }
    fn calc_level_zore(&mut self, list_fx: &LinkedList<MetaFX>, last_k: PureK) {
        let mut level_zore = Vec::new();
        let mut last = list_fx.front().unwrap();
        let mut i = 0;
        for this in list_fx {
            if i > 0 {
                match last.fx_type {
                    DIN => level_zore.push(Line {
                        level: 0,
                        index_range: (last.item.index_range.0, this.item.index_range.1),
                        line_range: (this.item.item.low, last.item.item.high),
                        difference: this.item.item.low - last.item.item.high,
                        direction: Direction::Down,
                        status: Status::Down,
                    }),
                    DI => level_zore.push(Line {
                        level: 0,
                        index_range: (last.item.index_range.0, this.item.index_range.1),
                        line_range: (last.item.item.low, this.item.item.high),
                        difference: this.item.item.high - last.item.item.low,
                        direction: Direction::Up,
                        status: Status::Down,
                    }),
                    _ => {
                        panic!("calc_level_zore fx_type err :{:?}", last)
                    }
                }
            }
            last = this;
            i += 1
        }
        if last.fx_type == DIN {
            level_zore.push(Line {
                level: 0,
                index_range: (last.item.index_range.0, last_k.index_range.1),
                line_range: (last_k.item.low, last.item.item.high),
                difference: last_k.item.low - last.item.item.high,
                direction: Direction::Down,
                status: Status::Grow,
            })
        } else {
            level_zore.push(Line {
                level: 0,
                index_range: (last.item.index_range.0, last_k.index_range.1),
                line_range: (last.item.item.low, last_k.item.high),
                difference: last_k.item.high - last.item.item.low,
                direction: Direction::Up,
                status: Status::Grow,
            })
        };
        self.line.insert(0, level_zore);
    }
    fn level_one_init(
        zore: &Vec<Line>,
        arr_up: &mut Vec<Line>,
        arr_down: &mut Vec<Line>,
        //arr_zs: &mut Vec<Zs>,
        i: usize,
        item: &Line,
        status: &mut Status,
    ) {
        if Line::check_trait_2_2(zore, i) {
            if item.direction == Direction::Up {
                *arr_up = clone_from(&zore[(i - 2)..i + 1]);
                *status = Status::QSUP;
            } else {
                *arr_down = clone_from(&zore[(i - 2)..i + 1]);
                *status = Status::QSDOWN;
            }
            //* 暂时不处理第一条线段准确性
            //} else if Line::check_trait_3(zore, i) {
            //    let (intersection, maxmin) = Line::calc_intersection_maxmin(&zore[(i - 2)..i + 1]);
            //    arr_zs.push(Zs {
            //        level: 1,
            //        index_range: (zore[i - 2].index_range.0, zore[i].index_range.1),
            //        center_range: intersection,
            //        wave_range: (zore[i - 1].line_range.0, zore[i - 1].line_range.1),
            //        status: Status::Grow,
            //        lines: clone_from(&zore[(i - 2)..i + 1]),
            //    });
            //    *status = Status::QSForming;
        }
    }
    fn calc_level_one(&mut self) {
        if let Some(zore) = self.line.get(&0) {
            let mut i = 0;
            let mut zs: HashMap<usize, Vec<Zs>> = HashMap::new();
            let mut level_one = Vec::new();
            let mut arr_up = Vec::new();
            let mut arr_down = Vec::new();
            let mut cache = Vec::new();
            let mut status = Status::Init;
            for item in zore {
                match i {
                    0..=1 => {}
                    2 => {
                        Self::level_one_init(zore, &mut arr_up, &mut arr_down, i, item, &mut status)
                    }
                    _ => {
                        if status == Status::Init {
                            Self::level_one_init(
                                zore,
                                &mut arr_up,
                                &mut arr_down,
                                i,
                                item,
                                &mut status,
                            )
                        } else {
                            match cache.len() {
                                0 if status == Status::QSDOWN || status == Status::QSUP => {
                                    cache.push(i);
                                }
                                _ => {
                                    if cache.len() % 2 == 0 {
                                        if Line::check_trait_x_y(zore, i, cache.len(), cache.len())
                                        {
                                            if status == Status::QSUP {
                                                level_one.push(
                                                    self.build_new_line_one_up(
                                                        &arr_up,
                                                        Status::Down,
                                                    ),
                                                );
                                                arr_up.clear();
                                                status = Status::QSDOWN;
                                                for index in i - cache.len()..i + 1 {
                                                    arr_down.push(zore[index].clone())
                                                }
                                            } else if status == Status::QSDOWN {
                                                level_one.push(self.build_new_line_one_down(
                                                    &arr_down,
                                                    Status::Down,
                                                ));
                                                arr_down.clear();
                                                status = Status::QSUP;
                                                for index in i - cache.len()..i + 1 {
                                                    arr_up.push(zore[index].clone())
                                                }
                                            } else {
                                                panic!("err not up or down")
                                            }
                                            self.bi_zs_check(&mut zs, &zore[cache[0]..i + 1]);
                                            cache.clear()
                                        } else {
                                            cache.push(i)
                                        }
                                    } else {
                                        if Line::check_trait_x_y(
                                            zore,
                                            i,
                                            cache.len() + 1,
                                            cache.len() + 1,
                                        ) {
                                            if status == Status::QSUP {
                                                for index in i - cache.len()..i + 1 {
                                                    arr_up.push(zore[index].clone());
                                                }
                                            } else if status == Status::QSDOWN {
                                                for index in i - cache.len()..i + 1 {
                                                    arr_down.push(zore[index].clone());
                                                }
                                            } else {
                                                panic!("err not up or down")
                                            }
                                            self.bi_zs_check(&mut zs, &zore[cache[0]..i + 1]);
                                            cache.clear()
                                        } else {
                                            cache.push(i)
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
                i += 1;
            }
            //还原指数
            if cache.len() > 0 {
                self.bi_zs_check(&mut zs, &zore[cache[0]..i]);
                if status == Status::QSUP {
                    level_one.push(self.build_new_line_one_up(&arr_up, Status::Grow));
                    for index in cache {
                        arr_down.push(zore[index].clone());
                    }
                    level_one.push(self.build_new_line_one_down(&arr_down, Status::Grow));
                } else if status == Status::QSDOWN {
                    level_one.push(self.build_new_line_one_down(&arr_down, Status::Grow));
                    for index in cache {
                        arr_up.push(zore[index].clone());
                    }
                    level_one.push(self.build_new_line_one_up(&arr_up, Status::Grow));
                } else {
                    panic!("err not up or down")
                }
                self.line.insert(1, level_one);
                self.zs = zs;
            }
        }
    }
    fn calc_zs(&mut self) {
        let mut arr_zs: Vec<Zs> = Vec::new();
        for (i, arr_list) in &self.line {
            if i == &0 {
                continue;
            }
            let mut cache: Vec<&Line> = Vec::new();
            let mut grow_center_range = (f64::MIN, f64::MAX);
            let mut center_range = (f64::MIN, f64::MAX);
            let mut wave_range = (f64::MAX, f64::MIN);
            let mut index_range: (i64, i64) = (0, 0);
            let mut last: &Line = &arr_list[0];
            let mut index = 2;
            while index < arr_list.len() {
                match cache.len() {
                    0 => {
                        if let Some(center_range_) =
                            self.zs_check_range(&arr_list[index - 2..index + 1])
                        {
                            cache.push(&arr_list[index - 1]);
                            center_range = center_range_;
                            wave_range = center_range_;
                            grow_center_range = center_range_;
                            last = &arr_list[index];
                            index_range = arr_list[index - 1].index_range
                        }
                    }
                    _ => {
                        if self.zs_check_center_range_next(center_range, &arr_list[index]) {
                            if let Some(new_range) =
                                self.zs_check_grow_center_range_next(grow_center_range, last)
                            {
                                grow_center_range = new_range
                            }
                            wave_range.0 = wave_range.0.min(last.line_range.0);
                            wave_range.1 = wave_range.1.max(last.line_range.1);
                            index_range.1 = last.index_range.1;
                            cache.push(last);
                            last = &arr_list[index];
                        } else {
                            let mut level = *i;
                            //match cache.len() {
                            //    9..=26 => level += 1,
                            //    27.. => level += 2,
                            //    _ => {}
                            //};
                            //* 暂时不考虑升级的情况 */
                            arr_zs.push(Zs {
                                level: level,
                                index_range: index_range,
                                center_range: center_range,
                                grow_center_range: grow_center_range,
                                wave_range: wave_range,
                                status: Status::Down,
                                lines: clone_from_(&cache[..]),
                            });
                            cache.clear();
                            grow_center_range = (f64::MIN, f64::MAX);
                            center_range = (f64::MIN, f64::MAX);
                            wave_range = (f64::MAX, f64::MIN);
                            index_range = (0, 0);
                        }
                    }
                }
                index += 1
            }

            if cache.len() > 1
                || (cache.len() == 1
                    && ((last.line_range.1.abs() - last.line_range.0.abs()).abs()
                        / (cache[0].line_range.1.abs() - cache[0].line_range.0.abs()).abs()
                        > 0.35))
            {
                let mut level = *i;
                //match cache.len() {
                //    9..=26 => level += 1,
                //    27.. => level += 2,
                //    _ => {}
                //};
                arr_zs.push(Zs {
                    level: level,
                    index_range: index_range,
                    center_range: center_range,
                    grow_center_range: grow_center_range,
                    wave_range: wave_range,
                    status: Status::Down,
                    lines: clone_from_(&cache[..]),
                });
            }
        }
        for i in arr_zs {
            let arr = self.zs.entry(i.level).or_insert(Vec::<Zs>::new());
            arr.push(i)
        }
    }
    fn zs_check_range(&self, arr_line: &[Line]) -> Option<(f64, f64)> {
        let mut c_range = (f64::MIN, f64::MAX);
        for i in arr_line {
            c_range.0 = c_range.0.max(i.line_range.0);
            c_range.1 = c_range.1.min(i.line_range.1);
        }
        if c_range.0 > c_range.1 {
            return None;
        }
        Some(c_range)
    }
    fn zs_check_grow_center_range_next(
        &self,
        grow_center_range: (f64, f64),
        next: &Line,
    ) -> Option<(f64, f64)> {
        let new_range = (
            grow_center_range.0.max(next.line_range.0),
            grow_center_range.1.min(next.line_range.1),
        );
        if new_range.0 > new_range.1 {
            return None;
        }
        Some(new_range)
    }
    fn zs_check_center_range_next(&self, center_range: (f64, f64), next: &Line) -> bool {
        let new_range = (
            center_range.0.max(next.line_range.0),
            center_range.1.min(next.line_range.1),
        );
        if new_range.0 > new_range.1 {
            return false;
        }
        true
    }
    fn analyze(&mut self) {
        let mut arr_status = Vec::<(usize, Status)>::new();
        for (level, arr_) in &self.zs {
            if let Some(last_zs) = arr_.last() {
                if let Some(arr_line) = self.line.get(level) {
                    let len_ = arr_line.len();
                    if arr_line[len_ - 1].line_range.1 + arr_line[len_ - 1].line_range.0
                        < last_zs.center_range.1 * 2.
                    {
                        continue;
                    }
                    if Self::range_contain(arr_line[len_ - 1].line_range, last_zs.center_range) {
                        continue;
                    } else {
                        if Self::range_contain(arr_line[len_ - 2].line_range, last_zs.center_range)
                        {
                            //证明是出中枢第一线段回踩 回踩没到30% 太早了
                            if (arr_line[len_ - 1].line_range.1 - arr_line[len_ - 1].line_range.0)
                                / (arr_line[len_ - 2].line_range.1
                                    - arr_line[len_ - 2].line_range.0)
                                < 0.3
                            {
                                continue;
                            }
                        } else {
                            //证明是出中枢第二线段次级别回升 超过30% 太晚了
                            if (arr_line[len_ - 1].line_range.1 - arr_line[len_ - 1].line_range.0)
                                / (arr_line[len_ - 2].line_range.1
                                    - arr_line[len_ - 2].line_range.0)
                                > 0.3
                            {
                                continue;
                            }
                        }

                        let status = Status::BuyTree;

                        arr_status.push((level.clone(), status));
                    }
                }
            }
        }
        if arr_status.len() > 0 {
            self.status = Some(arr_status)
        }
    }
    fn range_contain(a: (f64, f64), b: (f64, f64)) -> bool {
        let new_range = (a.0.max(b.0), a.1.min(b.1));
        return new_range.0 < new_range.1;
    }
    fn bi_zs_check(&self, zs: &mut HashMap<usize, Vec<Zs>>, arr_line: &[Line]) {
        if arr_line.len() < 9 {
            return;
        }
        let (a, b, c) = Line::calc_intersection_maxmin(arr_line);
        if arr_line.len() >= 9 && arr_line.len() < 27 {
            zs.entry(1).or_insert(Vec::new()).push(Zs {
                level: 1,
                index_range: c,
                center_range: a,
                grow_center_range: a,
                wave_range: b,
                status: Status::Down,
                lines: clone_from(arr_line),
            });
        }
        if arr_line.len() >= 27 && arr_line.len() < 81 {
            zs.entry(2).or_insert(Vec::new()).push(Zs {
                level: 2,
                index_range: c,
                center_range: a,
                grow_center_range: a,
                wave_range: b,
                status: Status::Down,
                lines: clone_from(arr_line),
            });
        };
        if arr_line.len() >= 81 {
            zs.entry(2).or_insert(Vec::new()).push(Zs {
                level: 3,
                index_range: c,
                center_range: a,
                grow_center_range: a,
                wave_range: b,
                status: Status::Down,
                lines: clone_from(arr_line),
            });
        }
    }
    fn build_new_line_one_up(&self, arr: &Vec<Line>, status: Status) -> Line {
        let len_ = arr.len();
        Line {
            level: 1,
            index_range: (arr[0].index_range.0, arr[len_ - 1].index_range.1),
            line_range: (arr[0].line_range.0, arr[len_ - 1].line_range.1),
            difference: arr[len_ - 1].line_range.1 - arr[0].line_range.0,
            direction: Direction::Up,
            status: status,
        }
    }
    fn build_new_line_one_down(&self, arr: &Vec<Line>, status: Status) -> Line {
        let len_ = arr.len();
        Line {
            level: 1,
            index_range: (arr[0].index_range.0, arr[len_ - 1].index_range.1),
            line_range: (arr[len_ - 1].line_range.0, arr[0].line_range.1),
            difference: arr[len_ - 1].line_range.0 - arr[0].line_range.1,
            direction: Direction::Down,
            status: status,
        }
    }
    fn calc_level_any(&mut self) {}
}

//-------------------------方法-------------------------------
