import pymysql
import datetime


class mysql():
    def __init__(self) -> None:
        self.db = pymysql.connect(host='127.0.0.1',
                                  user='quant',
                                  password='Quant.8.cc',
                                  database='quant')
        self.sql_builder = sql_builder()
        self.db.cursor().execute(self.sql_builder.build_create_table_swl())
        self.db.cursor().execute(self.sql_builder.build_create_table_day())

    def cursor(self):
        return self.db.cursor()

    def print_mysql_version(self):
        cursor = self.db.cursor()
        cursor.execute("SELECT VERSION()")
        data = cursor.fetchone()
        print("Database version : %s " % data)


class sql_builder():
    def __init__(self) -> None:
        self.table_day = "day"

    def build_new_recent(self, table: str):
        return """SELECT max(date),code from {} GROUP BY code""".format(table)
