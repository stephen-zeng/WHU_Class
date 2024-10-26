# 介绍
这个半自动代码可以帮助你将妮呜WHU的课表导入日历中统一管理

# 需要安装的包
+ icalendar
+ uuid
+ json

安装命令：
```
pip3 install icalendar uuid json
```
# 课表的获取
首先，用电脑打开智慧珞珈的手机版，登录`zhlj.whu.edu.cn`之后，按`F12`打开开发者控制台，然后切换模拟设备为一台手机，看下面的GIF。
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/master/guidance/1.gif)
然后网页打开课表页面，控制台打开Network（网络），准备抓取json。
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/master/guidance/2.gif)
刷新页面，在控制台找到`getCurriculumData?data=`，双击打开，这就是你的这一周课表的JSON数据了。
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/master/guidance/3.gif)
要查看下一周的课表JSON，只要在网页端查看对应周的课表，然后在控制台找到`getCurriculumDate?data=2024-10-27`的文件，其中`2024-10-27`为这一周的星期天，同样双击打开
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/master/guidance/4.gif)
注意，复制的时候勾上“Pretty-print（格式化）”
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/master/guidance/5.gif)

# 使用说明
首先，更改`all.json`里面的内容，json里面不能放注释，所以我就在下面放格式吧
```json
{
    // year, month, day分别对应着第一周星期天的年份，月份和日期。注意，妮呜的一周从星期天开始
    "year": 2024, 
    "month": 10,
    "day": 13,
    "data": [
        {
            "week": 1,
            "data": [
                // 这里放对应周的课表JSON，直接将JSON里的所有内容  
            ]
        },
        {
            "week": 2,
            "data": [
                // 以此类推，注意星期数要连续
            ]
        } // 注意，最后一周的大括号不要有逗号结尾
    ]
}
```
更改完json之后，运行`main.py`，Python版本的话我本地是3.13.0，没有测试过其他版本。运行顺利的话会出现`all.ics`这个日历文件，但还需要处理。用文本编辑器打开`head.ics`和`all.ics`，将`all.ics`的这一些内容替换成`head.ics`的所有内容：
```ics
BEGIN:VCALENDAR
VERSION:2.0
X-APPLE-CALENDAR-COLOR:#540EB9
X-WR-CALNAME:生成ics文件测试
X-WR-TIMEZONE:Asia/Shanghai
```
（幽默果子的日历头部要求极其严格hhh
然后，你就可以将更改好的日历文件用Apple设备打开了。理论上所有支持通用日历文件的APP都可以使用。