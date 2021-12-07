import pymysql
import datetime


class mysql():
    def __init__(self) -> None:
        self.db = pymysql.connect(host='127.0.0.1',
                                  user='quant',
                                  password='Quant.8.cc',
                                  database='quant')
        self.sql_builder = sql_builder()

    def __cursor(self):
        return self.db.cursor()

    def query_day_recent(self):
        cursor = self.__cursor()
        cursor.execute(self.sql_builder.build_recent('day'))
        return cursor.fetchall()

    def query_15m_recent(self):
        cursor = self.__cursor()
        cursor.execute(self.sql_builder.build_recent('15m'))
        return cursor.fetchall()

    def insert_many_code(self, list, log):
        try:
            with self.db.cursor() as cursor:
                cursor.executemany(
                    self.sql_builder.build_insert_many_code(), list)
                self.db.commit()
                cursor.close()
        except BaseException as err:
            log.error('insert_many_code err {}'.format(err))

    def insert_many_day(self, list, log):
        try:
            with self.db.cursor() as cursor:
                cursor.executemany(
                    self.sql_builder.build_insert_many('day'), list)
                self.db.commit()
                cursor.close()
        except BaseException as err:
            log.error('insert_many_day err {}'.format(err))

    def insert_many_15m(self, list, log):
        try:
            with self.db.cursor() as cursor:
                cursor.executemany(
                    self.sql_builder.build_insert_many('15m'), list)
                self.db.commit()
                cursor.close()
        except BaseException as err:
            log.error('insert_many_15m err {}'.format(err))

    def get_all_code(self):
        with self.db.cursor() as cursor:
            cursor.execute(
                self.sql_builder.build_query_all_code())
            return cursor.fetchall()

    def get_all_day(self):
        with self.db.cursor() as cursor:
            cursor.execute(
                self.sql_builder.build_query("day"))
            return cursor.fetchall()

    def get_all_15m(self):
        with self.db.cursor() as cursor:
            cursor.execute(
                self.sql_builder.build_query("15m"))
            return cursor.fetchall()

    def update_many_day(self, list, log):
        try:
            with self.db.cursor() as cursor:
                cursor.executemany(
                    self.sql_builder.build_update_timestamp('day'), list)
                self.db.commit()
                cursor.close()
        except BaseException as err:
            log.error('update_many_day err {}'.format(err))
        log.info('update_many_day fin len {}'.format(len(list)))

    def update_many_15m(self, list, log):
        try:
            with self.db.cursor() as cursor:
                cursor.executemany(
                    self.sql_builder.build_update_timestamp('15m'), list)
                self.db.commit()
                cursor.close()
        except BaseException as err:
            log.error('update_many_15m err {}'.format(err))
        log.info('update_many_15m fin len {}'.format(len(list)))


class sql_builder():

    def build_recent(self, table: str):
        return """SELECT code, max(date) from {} GROUP BY code""".format(table)

    def build_query(self, table: str):
        return """SELECT id,date from {}""".format(table)

    def build_update_timestamp(self, table: str):
        return """UPDATE {} SET timestamp=%s WHERE id=%s""".format(table)

    def build_all_code(self, table: str):
        return """SELECT DISTINCT(code) FROM {}""".format(table)

    def build_insert_many_code(self):
        return """INSERT INTO code(start_date, end_date, code, name, type_,timestamp)VALUES ( %s, %s, %s, %s, % s,%s)"""

    def build_query_all_code(self):
        return """SELECT code FROM code where end_date>=NOW()"""

    def build_insert_many(self, table):
        return """INSERT INTO {}(date, code, open, high, low, close, volume,timestamp)VALUES ( % s, % s, % s, % s, % s, % s, % s,% s)ON DUPLICATE KEY UPDATE  open=%s,high=%s,low=%s,close=%s,volume=%s""".format(table)
