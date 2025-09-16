# 介绍
这个半自动程序可以帮助你将妮呜WHU的课表导入日历中统一管理

## 🆕 Web版本（推荐）
现在支持Material Design风格的Web界面，更加美观易用！

### 启动Web服务
```bash
# 下载程序后，直接运行
./whu_class -web

# 或指定端口
./whu_class -web -port=9000
```

访问 `http://localhost:8080` 即可使用Web界面。

如果部署在服务器上，其他人可以通过服务器IP地址访问，实现多人共享使用。

## 📱 Web界面特点
- 🎨 Material Design风格界面
- 📱 响应式设计，支持移动设备
- 🔧 模块化前端架构，嵌入式静态资源
- 📝 详细的使用说明
- ✅ 实时数据验证
- 📊 课程信息预览
- 💾 一键下载日历文件
- 🌐 支持多人使用
- ⚡ 快速加载，无外部依赖

## ⌨️ 命令行版本

# 步骤
## 0. 下载对应的程序
> Mac下载darwin版本
 
https://github.com/stephen-zeng/WHU_Class/releases

## 1. 进入本科教务系统的个人课表查询
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/refs/heads/main/guide/1.jpg)

## 2. 按照动图的操作复制cURL命令的连接
由于隐私信息，左边的网页内容没有展示，此时像正常查询课表一样点击查询，完成滑动验证码，然后右边控制台就有反应了
![](https://raw.githubusercontent.com/stephen-zeng/WHU_Class/refs/heads/main/guide/2.gif)

## 3. 打开下载的程序，按照动图的操作获取日志文件
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
