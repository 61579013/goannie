## goannie
goannie 是一个视频资源采集下载的实用工具。目前它还没有 gui ，通过命令行交互的方式操作。


## 获取方式
https://gitee.com/rock_rabbit/goannie/releases

## 特点
* 视频批量采集
* 下载cookie设置
* 重复下载过滤

## 应用
* 自媒体视频制作者

## 环境
开发测试：
`windows10 x64`

## 支持
```

                                        __
   __     ___      __      ___     ___ /\_\     __
 /'_ `\  / __`\  /'__`\  /' _ `\ /' _ `\/\ \  /'__`\
/\ \L\ \/\ \L\ \/\ \L\.\_/\ \/\ \/\ \/\ \ \ \/\  __/
\ \____ \ \____/\ \__/.\_\ \_\ \_\ \_\ \_\ \_\ \____\
 \/___L\ \/___/  \/__/\/_/\/_/\/_/\/_/\/_/\/_/\/____/
   /\____/
   \_/__/
        version: v0.0.20        updateTime: 2020-10-07

支持平台
|-----------------       优酷视频       -----------------|
task: one       info: 单视频            https://v.youku.com/v_show/id_XNDg2MTM3MjMyMA==.html
task: userList  info: 作者视频          http://i.youku.com/i/UNjMwMTY2MDUyMA==
cookie 设置：在goannie.exe同级目录中新建 douyin.txt 写入name=value;name=value....格式即可。
ccode 和 ckey 设置：在goannie.exe同级目录中新建 ccode.txt 和 ckey.txt 写入其中即可。

|-----------------       抖音视频       -----------------|
task: one       info: 单视频            https://www.iesdouyin.com/share/video/6877354382132808971
task: userList  info: 作者视频          https://www.iesdouyin.com/share/user/2836383897749943?sec_uid=xxxxx
task: shortURL  info: 短链接            https://v.douyin.com/JDq8uv7/
cookie 设置：在goannie.exe同级目录中新建 douyin.txt 写入name=value;name=value....格式即可。

|-----------------       腾讯视频       -----------------|
task: one       info: 单视频            https://v.qq.com/x/cover/mzc00200agq0com/r31376lllyf.html
task: detail    info: 腾讯剧集页        https://v.qq.com/detail/5/52852.html
task: userList  info: 作者视频          https://v.qq.com/s/videoplus/1790091432
task: lookList  info: 看作者作品列表    look https://v.qq.com/s/videoplus/1790091432
cookie 设置：在goannie.exe同级目录中新建 tengxun.txt 写入name=value;name=value....格式即可。

|-----------------       火锅视频       -----------------|
task: userList  info: 作者视频          https://huoguo.qq.com/m/person.html?userid=18590596
task: lookList  info: 看作者作品列表    look https://huoguo.qq.com/m/person.html?userid=18590596
cookie 设置：在goannie.exe同级目录中新建 tengxun.txt 写入name=value;name=value....格式即可。

|-----------------       爱奇艺视频     -----------------|
task: one       info: 单视频            https://www.iqiyi.com/v_1fr4mggxzpo.html
task: userList  info: 作者视频          https://www.iqiyi.com/u/2182689830
task: detail    info: 爱奇艺剧集页      https://www.iqiyi.com/a_19rrht2ok5.html
cookie 设置：在goannie.exe同级目录中新建 iqiyi.txt 写入name=value;name=value....格式即可。

|-----------------       西瓜视频       -----------------|
task: one       info: 单视频            https://www.ixigua.com/6832194590221533707
task: userList  info: TA的视频          https://www.ixigua.com/home/85383446500/video/
task: lookList  info: 看作者作品列表    look https://www.ixigua.com/home/85383446500/video/
cookie 设置：在goannie.exe同级目录中新建 xigua.txt 写入name=value;name=value....格式即可。

|-----------------       好看视频       -----------------|
task: one       info: 单视频            https://haokan.baidu.com/v?vid=3881011031260239591
task: userList  info: 作者视频          https://haokan.baidu.com/author/1649278643844524
cookie 设置：在goannie.exe同级目录中新建 haokan.txt 写入name=value;name=value....格式即可。

|-----------------       哔哩哔哩       -----------------|
task: one       info: 单视频            https://www.bilibili.com/video/BV1iK4y1e7uL
task: userList  info: TA的视频          https://space.bilibili.com/337312411
cookie 设置：在goannie.exe同级目录中新建 bilibili.txt 写入name=value;name=value....格式即可。

下载统计
腾讯视频：0  爱奇艺视频：0  好看视频：0  哔哩哔哩：0  西瓜视频：0  抖音视频：0  优酷视频：51

$ 请输入保存路径：
```

## 附属程序
annie.exe  
aria2c.exe  
redis-server.exe  
ffmpeg.exe  
以上会在启动时请求下载。  
存储位置：`%APPDATA%/goannie/bin`

## 感谢
https://github.com/iawia002/annie

## 截图
![goannie 截图](http://image.68wu.cn/blog/goannie_20201007.png)
## 挖坑
计划会写一个GUI版本

![videoSpade 视频铲](http://image.68wu.cn/blog/videoSpade_20201007.png)