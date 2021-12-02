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

    def query_day_all_code(self):
        cursor = self.__cursor()
        cursor.execute(self.sql_builder.build_all_code('day'))
        return cursor.fetchall()


class sql_builder():

    def build_recent(self, table: str):
        return """SELECT code, max(date) from {} GROUP BY code""".format(table)

    def build_all_code(self, table: str):
        return """SELECT DISTINCT(code) FROM {}""".format(table)
