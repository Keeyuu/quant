import jqdatasdk
import copy
import log
import dao
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.schedulers.blocking import BlockingScheduler
from datetime import date, datetime
import time
from jqdatasdk import finance, query


class DataSyncer:
    def __init__(self) -> None:
        # self.__auth()
        self.logger = log.Loggers()
        self.db = dao.mysql()
        # self.job_async = BackgroundScheduler()
        # self.job = BlockingScheduler()
        # self.load_job()

    def __heartbeat(self):
        self.logger.info('heartbeat')

    def __auth(self):
        jqdatasdk.auth('17675677591', 'Keeyu.cc.9')

    def load_job(self):
        self.job.add_job(self.__auth, 'cron', day_of_week='0-5', hour='16-19')
        self.job.add_job(self.check_load, 'cron',
                         day_of_week='0-5', hour='16-19')
        self.job_async.add_job(self.__heartbeat, 'cron', minute='*')
        self.job_async.start()
        self.job.start()

    def insert_code(self):
        codes = jqdatasdk.get_all_securities(
            types=['etf', 'stock', 'lof'], date=None)
        for code, row in codes.iterrows():
            self.db.insert_many_code(
                [(row['start_date'], row['end_date'], code, row['display_name'], row['type'])], self.logger)
        self.logger.info('has insert_code')

    def pull_bar_15m(self, code, count):
        rsb = jqdatasdk.get_bars(code, count=count, fields=[
            'date', 'open', 'close', 'high', 'low', 'volume'], unit='15m', include_now=False, df=True)
        list = []
        for i in rsb.index:
            data = dict(rsb.iloc[i])
            data['code'] = code
            list.append((data['date'], data['code'], data['open'],
                         data['high'], data['low'], data['close'], data['volume'], data['volume']))
        self.db.insert_many_15m(list, self.logger)

    def pull_bar_day(self, code, count):
        rsb = jqdatasdk.get_bars(code, count=count, fields=[
            'date', 'open', 'close', 'high', 'low', 'volume'], unit='1d', include_now=False, df=True)
        list = []
        for i in rsb.index:
            data = dict(rsb.iloc[i])
            data['code'] = code
            list.append((data['date'], data['code'], data['open'],
                         data['high'], data['low'], data['close'], data['volume'], data['volume']))
        self.db.insert_many_day(list, self.logger)

    def check_load(self):
        new = jqdatasdk.get_trade_days(count=1)[0]
        code_list = [i[0] for i in self.db.get_all_code()]
        day_need = {}
        m15_need = {}
        m15_recent = {i[0]: i[1] for i in self.db.query_15m_recent()}
        day_recent = {i[0]: i[1] for i in self.db.query_day_recent()}
        for code in code_list:
            if code not in day_recent.keys():
                day_need[code] = 1200
            elif day_recent[code] < new:
                day_need[code] = min((new-day_recent[code]).days, 1200)

            if code not in m15_recent.keys():
                m15_need[code] = 1500
            elif m15_recent[code].date() < new:
                m15_need[code] = min(
                    ((new-m15_recent[code].date()).days)*16, 1500)
        for i in day_need:
            self.pull_bar_day(i, day_need[i])
            self.logger.info('fin day update {}'.format(i))
        for i in m15_need:
            self.pull_bar_15m(i, m15_need[i])
            self.logger.info('fin 15m update {}'.format(i))

    def init_time_timestamp(self):
        list = []
        for i in self.db.get_all_15m():
            list.append((int(i[1].timestamp() * 1000), i[0]))
            if len(list) >= 50000:
                self.db.update_many_15m(list, self.logger)
                list = []
        if len(list) >= 0:
            self.db.update_many_15m(list, self.logger)
            list = []

        for i in self.db.get_all_day():
            list.append((int(datetime.strptime(
                str(i[1]), '%Y-%m-%d').timestamp() * 1000), i[0]))
            if len(list) >= 50000:
                self.db.update_many_day(list, self.logger)
                list = []
        if len(list) >= 0:
            self.db.update_many_day(list, self.logger)
            list = []


if __name__ == '__main__':
    DataSyncer().init_time_timestamp()
