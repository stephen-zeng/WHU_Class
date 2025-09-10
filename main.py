from icalendar import Calendar, Event
from get_calendar import get_json
import uuid
import json

get_json()

with open('calendar.json', 'r', encoding="utf-8") as zongbiao:
    KeBiao = json.load(zongbiao)
    allkb = KeBiao['data']

cal = Calendar()
cal.add('VERSION','2.0')
cal.add('X-WR-CALNAME','生成ics文件测试')
cal.add('X-APPLE-CALENDAR-COLOR','#540EB9')
cal.add('X-WR-TIMEZONE','Asia/Shanghai')

bgtime = ['0','080000','085000','095000','104000','113000','140500','145500','154500','164000','173000','183000','192000','201000']
edtime = ['0','084500','093500','103500','112500','121500','145000','154000','163000','172500','181500','191500','200500','205500']


# 月份天数（考虑闰年）
def is_leap_year(year):
    """判断是否为闰年"""
    return (year % 4 == 0 and year % 100 != 0) or (year % 400 == 0)

days_in_month = [0, 31, 29 if is_leap_year(KeBiao['year']) else 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]

year = KeBiao['year']
month = KeBiao['month']
wkday = KeBiao['day']
wkday = wkday - 7

for week in allkb:
    wkday = wkday + 7
    if (wkday > days_in_month[month]):
        wkday = wkday - days_in_month[month]
        month = month + 1 
        if (month > 12):
            month = 1
            year = year + 1
    strday = "%02d" % wkday
    strmonth = "%02d" % month
    basedate = int(str(year) + strmonth + strday)
    for day in week['data'][0]['data']:
        if (day['day']==8):
            continue
        date = basedate + day['day']
        courses = day['curriculumList']
        for course in courses:
            if (course['hasClass']==False):
                continue
            event = Event()
            bg = bgtime[course['fromClass']]
            ed = edtime[course['endClass']]
            classname = course['name']
            if (classname == None):
                classname = '没有上课地点'
            classroom = course['classroom']
            teacher = course['teacher']
            event.add('UID',str.upper(str(uuid.uuid4())))
            event.add('DTSTART;TZID=Asia/Shanghai', str(date) + 'T' + bg)
            event.add('DTEND;TZID=Asia/Shanghai', str(date) + 'T' + ed)
            event.add('SUMMARY', classname)
            event.add('SEQUENCE','0')
            event.add('DESCRIPTION','Teacher: ' + teacher)
            event.add('LOCATION', classroom)
            cal.add_component(event)

f = open('calendar.ics', 'wb')
f.write(cal.to_ical())
f.close()

# 要插入的内容
new_content = """BEGIN:VCALENDAR
CALSCALE:GREGORIAN
PRODID:-//Apple Inc.//macOS 15.2//EN
VERSION:2.0
X-APPLE-CALENDAR-COLOR:#CC73E1
BEGIN:VTIMEZONE
TZID:Asia/Shanghai
BEGIN:STANDARD
DTSTART:19890917T020000
RRULE:FREQ=YEARLY;UNTIL=19910914T170000Z;BYMONTH=9;BYDAY=3SU
TZNAME:GMT+8
TZOFFSETFROM:+0900
TZOFFSETTO:+0800
END:STANDARD
BEGIN:DAYLIGHT
DTSTART:19910414T020000
RDATE:19910414T020000
TZNAME:GMT+8
TZOFFSETFROM:+0800
TZOFFSETTO:+0900
END:DAYLIGHT
END:VTIMEZONE
"""

with open('calendar.ics', 'r', encoding='utf-8') as file:
    lines = file.readlines()

updated_lines = new_content.splitlines(keepends=True) + lines[5:]

with open('calendar.ics', 'w', encoding='utf-8') as file:
    file.writelines(updated_lines)

print("calendar.ics创建成功！")