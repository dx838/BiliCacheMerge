# BiliCacheMerge 编码方案

> 将 B 站客户端缓存合并为可播放的 MP4。

## 一、项目目标

读取 B 站客户端缓存目录，按合集统计已缓存的分 P，校验文件完整性，剥离 m4s 头部后调用 ffmpeg 合并为 MP4。

## 二、技术方案

- **语言/构建**：Go（`CGO_ENABLED=0`，沿用现有 `build_*.bat` 跨平台脚本，产物名取 `env` = `BiliCacheMerge`）。
- **配置**：ini 文件，含缓存目录、输出目录、ffmpeg 路径三项；解析用 `gopkg.in/ini.v1`。
- **形态**：命令行批处理 + 控制台日志，不做 TUI/数据库。

## 三、已确认的输入事实（基于真实缓存实测）

| 项 | 结论 |
|---|---|
| `videoInfo.json` | 单行、扁平 JSON，字段全在顶层，无嵌套 |
| 关键字段类型 | `groupId` = string（合集 bvid）**或 number（单稿数字 ID）**；`bvid`/`title`/`groupTitle` = string；`itemId`/`aid`/`cid`/`p` = number（`aid` 超 int32，需 int64） |
| `groupId` | 两种取值：字符串时为合集 `bvid`（如 `BV12NTi6cEVV`），数字时为单稿数字 ID（如 `5254035`） |
| `p` | 组内序号，实测 `1..108` |
| `title` | 与 `tabName` 相同，且开头自带 p 序号，如 `"019 【邪恶的计划】  An Evil Plan"` |
| 元数据文件 | 目录内同时存在 `videoInfo.json` 与 `.videoInfo`，内容逐字相同 |
| 完整性凭据 | `.playurl` 存在 + 视频/音频 m4s 均存在 |
| m4s 命名 | `{itemId}-{n}-{fileId}.m4s`；本缓存中间段恒为 `1`；`30064`/`30080`=视频、`30280`=音频 |
| m4s 头部 | 前 9 字节是 9 个 ASCII 字符 `0`（`0x30`），其后才是 `00 00 00 20 ftyp`；**必须剥离**，否则 ffmpeg 报 `Invalid data found` |
| ffmpeg | 实测环境 `C:\ProFiles\cmd\ffmpeg.bat`（v8.0，`.bat` 需经 `cmd /c` 执行） |

### videoInfo.json 结构示例（节选）

```json
{
  "type": "ugc",
  "codecid": 7,
  "groupId": "BV12NTi6cEVV",
  "itemId": 39556222659,
  "aid": 116843612280112,
  "cid": 39556222659,
  "bvid": "BV12NTi6cEVV",
  "p": 1,
  "title": "001 【石猴出世】 The Monkey",
  "groupTitle": "【英语听力暴涨】108集《西游记》英文版动画！中英字幕｜边看边学｜合集珍藏 | 看动画学英语 | 英语启蒙"
}
```

## 四、已确定的规则细节

| 项 | 决定 |
|---|---|
| title 数字前缀 | 若开头数字（含前导 0）等于 `p`，则去掉该前缀及后随空白 |
| 非法字符 | `\ / : * ? " < > |` 及 ASCII 控制符替换为 `_`；再去除结尾的点和空格；目录名与文件名都清洗 |
| title 截断 | 按 **字符（rune）** 截断到 32 位 |
| 文件名 | `{p 补零 3 位}_{title}.mp4`；title 处理后为空则退化为 `{p 补零 3 位}.mp4` |
| 输出路径 | `{输出目录}/{清洗后的 groupTitle}/{文件名}`；**输出目录须已存在**，否则提示并退出；合集子目录自动创建 |
| 已存在输出 | 跳过（支持增量重跑） |
| 音视频区分 | 按 fileId 后缀精确匹配：`-30064.m4s`/`-30080.m4s`=视频、`-30280.m4s`=音频 |
| 合并参数 | `-c:v copy -c:a copy`（源音频为 AAC，直接封装，兼容 Android） |

## 五、文件结构

| 文件 | 职责 |
|---|---|
| `go.mod` | 模块 `bilicachemerge`，依赖 `gopkg.in/ini.v1` |
| `config.go` | 读 ini（缓存目录/输出目录/ffmpeg 路径），校验并给出中文错误 |
| `cache.go` | 扫描子目录、读元数据 JSON、按 `groupId` 分组、按 `p` 排序、完整性校验 |
| `naming.go` | 去前缀、清洗非法字符、rune 截断、命名拼接 |
| `ffmpeg.go` | 剥离 m4s 头部到临时文件 → 调 ffmpeg 合并 → 清理临时文件 |
| `main.go` | 串联流程、输出汇总报告 |
| `config.ini` | 配置模板 |
| `naming_test.go` / `cache_test.go` | 单元测试 |

## 六、核心流程

1. 读 ini 配置并校验输出目录、ffmpeg、bilibili 缓存目录是否存在，异常时提示错误并等待用户按任意键退出。
2. 扫描缓存根目录的一级子目录，读取元数据（`videoInfo.json` 优先，缺失回退 `.videoInfo`），跳过无元数据目录。
3. 按 `groupId` 分组，按 `p` 升序排列。
4. 列出合集（序号、`groupId`、`groupTitle`），等待用户输入序号或 `groupId` 选择要处理的合集。
5. 对选中合集的每项做完整性校验：`.playurl`、视频 m4s、音频 m4s 是否齐全，缺失则记录原因。
6. 对完整项：目标 MP4 已存在则跳过；否则剥离 9 字节头到临时文件，调用 ffmpeg 合并。
7. 打印汇总：成功 / 已存在 / 不完整 / 失败；随后等待用户按任意键退出。

## 七、任务分解与验证

| # | 任务 | 验证标准 |
|---|---|---|
| 1 | 初始化 `go.mod` 与依赖 | `go build` 通过 |
| 2 | `config.go` | 缺项/路径不存在报中文错误；正常样例正确读出三项 |
| 3 | `cache.go` 扫描+解析+分组 | 对真实缓存正确识别合集与 `p=1..108` |
| 4 | `cache.go` 完整性校验 | 缺音频/缺 playurl 被标为不完整并给出原因 |
| 5 | `naming.go` | `p=19` + `"019 【邪恶的计划】  An Evil Plan"` → `019_【邪恶的计划】  An Evil Plan.mp4`；`|` 被替换；超 32 rune 截断 |
| 6 | `ffmpeg.go` | 对 1 个完整 P 产出可播放 MP4（前 9 字节已剥离） |
| 7 | `main.go` | 全量跑通并打印汇总；再次运行全部命中「跳过」 |
| 8 | `config.ini` 模板 | 注释说明三项配置 |

## 八、运行方式

```
go build -o bin\BiliCacheMerge.exe .
bin\BiliCacheMerge.exe -config config.ini
```

启动后先校验三项路径；通过则列出合集，输入序号或 `groupId` 选择要处理的合集，随后开始合并。

## 九、后续可调整项

- **fileId 兜底**：`cache.go` 顶部已标注 `TODO`，后缀集中在 `videoStreamSuffixes` / `audioStreamSuffixes`，可在 `buildItem` 中改为按文件大小判定。
- **音频参数**：当前 `-c:a copy`（源音频为 AAC，直接封装，兼容 Android；AC3 在 Android 上无声）；若需统一编码可改 `-c:a aac`。
