import log
import dao
from apscheduler.schedulers.background import BackgroundScheduler
from datetime import datetime
import time


class DataSyncer:
    def __init__(self) -> None:
        self.logger = log.Loggers()
        self.db = dao.mysql()

    def test(self):
        #a = self.db.query_day_recent()
        # print(a)
        pass


if __name__ == '__main__':
    DataSyncer().test()

    def job():
        print(datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
    # 定义BlockingScheduler
    sched = BackgroundScheduler()
    sched.add_job(job, 'cron', second='1-59', day_of_week='0-5', hour='9-19')
    sched.start()

print("yes")
time.sleep(100)
