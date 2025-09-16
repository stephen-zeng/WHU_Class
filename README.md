# 介绍
这个半自动程序可以帮助你将妮呜WHU的课表导入日历中统一管理

# 步骤
## 0. 下载对应的程序
> Mac下载darwin版本
 
https://github.com/stephen-zeng/WHU_Class/releases

## 1. 进入本科教务系统的个人课表查询
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/refs/heads/main/guide/1.jpg)

## 2. 按照动图的操作获取cURL链接
由于隐私信息，左边的网页内容没有展示，此时像正常查询课表一样点击查询，完成滑动验证码，然后右边控制台就有反应了
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/refs/heads/main/guide/2.gif)

## 3.0 打开网页，粘贴cURL链接，获取课表，下载文件
https://class.0x535a.cn

## 3.1 打开下载的程序，按照动图的操作获取日志文件
> Mac在终端打开，看动图（~~忽略我下错文件~~）
> ![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/refs/heads/main/guide/4.gif)
> Linux你想怎么开就怎么开（刻板印象：都用Linux了还不会运行文件吗）

然后运行就好了
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/refs/heads/main/guide/3.gif)

## 🚀 开发者说明

### 构建项目
```bash
go build -o whu_class
```

### Docker
```bash
# amd64
docker run -d --name whu-class -p %EXPOSE_PORT%:8080 0w0w0/whu-class:%VERSION%

# arm64
docker run -d --name whu-class -p %EXPOSE_PORT%:8080 0w0w0/whu-class-arm:%VERSION%
```

### 运行模式
```bash
# 命令行模式（默认）
./whu_class

# Web服务器模式
./whu_class -web

# 查看帮助
./whu_class -h
```

### 测试
```bash
go test
```
