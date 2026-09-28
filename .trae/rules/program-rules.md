# SM3SUM 项目规则


## 项目目标

制作一个将 bilibili 缓存合并为一个视频。

## 项目功能

使用 ini 配置 bili 视频缓存位置、ffmpeg 安装位置、视频输出目录。

通过配置 bilibili 的缓存目录， 读取 videoInfo.json 文件 获取 groupId 、groupTitle 、aid、cid、bvid、p、title、itemId，根据 p 数据统计已经缓存的视频数据。p为分组中按顺序每个视频对应一个编号。在统计数据时，要检查缓存文件是否齐全（m4s、playurl、videoInfo），两个文件中视频文件一般较大，如果没有合适的方式判断为视频或音频文件，可按大小进行判断，一般音频文件较小。

输出时，文件名称为 p 字段扩充至 3 位，不够的用 0 进行填充；使用下划线连接；title 字段数据 ，title 最长 32位，超过的截断。

调用 ffmpeg 来合成视频。


## 缓存内容说明

见 `biliCache.txt` 中的说明。


