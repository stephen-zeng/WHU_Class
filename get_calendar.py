import requests
import json
from datetime import datetime, timedelta

# ----------------------------------------------
# 请求的cookies和headers(请根据README.md中的方法获取)
cookies = {
    
}

headers = {
    
}

#----------------------------------------------

# 函数用于生成两个日期之间的所有日期
def get_date_range(start_date, end_date):
    start = datetime.strptime(start_date, "%Y-%m-%d")
    end = datetime.strptime(end_date, "%Y-%m-%d")
    delta = timedelta(days=7)
    
    current = start
    while current <= end:
        yield current.strftime("%Y-%m-%d")
        current += delta

def get_json():
    start_date = input("请输入开始日期 (格式: YYYY-MM-DD): ")
    end_date = input("请输入结束日期 (格式: YYYY-MM-DD): ")
    date_range = list(get_date_range(start_date, end_date))
    first_date = date_range[0]
    first_date_obj = datetime.strptime(first_date, "%Y-%m-%d")
    year = first_date_obj.year
    month = first_date_obj.month
    day = first_date_obj.day
    data_structure = {
        "year": year,
        "month": month,
        "day": day,
        "data": []
    }
    week = 0
    for date in get_date_range(start_date, end_date):
        week += 1
        params = {
            'date': date,
        }
        response = requests.get(
            'https://zhlj.whu.edu.cn/mobile/homepageapi/getCurriculumData',
            params=params,
            cookies=cookies,
            headers=headers,
        )
        week_data = {
            "week": week,
            "data": [json.loads(response.text)]
        }
        data_structure["data"].append(week_data)
    with open("calendar.json", "w", encoding="utf-8") as json_file:
        json.dump(data_structure, json_file, ensure_ascii=False, indent=4)
