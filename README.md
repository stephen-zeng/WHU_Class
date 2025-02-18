# 介绍
### 这个半自动代码可以帮助你将妮呜WHU的课表导入日历中统一管理

# 0. 需要安装的依赖
+ icalendar
+ uuid
+ uncurl
+ datetime
+ requests

### 安装命令：
```
pip3 install icalendar uuid uncurl datetime requests
```
# 1. cURL的获取
### 用电脑打开[智慧珞珈的手机版](https://zhlj.whu.edu.cn/mobile/curriculum)，登录`zhlj.whu.edu.cn`之后，按`F12`打开开发者控制台，然后切换模拟设备为一台手机，看下面的GIF。
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/master/guidance/1.gif)

### 网页打开课表页面，任选一周，控制台打开Network（网络）。
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/master/guidance/2.gif)

### 刷新页面，在控制台找到`getCurriculumData?data=`，右键选择复制为`cURL(bash)`
![alt text](guidance/image.png)


# 2. 日历的生成
### 运行`main.py`，输入获取的cURL，粘贴之后可能需要按一次或两次回车，直到出现输入提示。按格式输入开始日期（第一周的周日）和结束日期（最后一周的周日），测试用的Python版本`3.13.0`和`3.11.9`，没有测试过其他版本，理论可行。运行顺利的话会出现`calendar.ics`这个日历文件。

![3.gif](guidance/3.gif)

### 恭喜！你可以将`calendar.ics`用Apple设备打开了。理论上，所有支持通用日历文件的APP都可以使用。