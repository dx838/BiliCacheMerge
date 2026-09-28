# BiliCacheMerge

将哔哩哔哩（bilibili）Windows 客户端的离线缓存（分离的音频 / 视频 `.m4s` 文件）批量合并为可直接播放的 MP4。

- 自动扫描缓存目录，按合集（`groupId`）分组、按分 P（`p`）排序
- 校验每个分 P 的缓存完整性（`videoInfo.json` / `.videoInfo`、`.playurl`、音频与视频 m4s）
- 交互式选择要处理的合集，输出处理汇总
- 自动剥离 m4s 文件开头的 9 个占位字节，调用 ffmpeg 无损合并（`-c:v copy -c:a copy`）
- 生成的 MP4 在 Windows、Android、iOS 均可正常播放

## 运行环境

- Windows（缓存结构基于 bilibili Windows 客户端）
- [bilibili 客户端](https://app.bilibili.com/?spm_id_from=333.1007.0.0)（用于离线缓存视频）
- [ffmpeg](https://ffmpeg.org/download.html)（用于音视频合成，需自行下载并在配置文件中填写路径）
- 从源码构建需要 [Go](https://go.dev/dl/) 1.21+

## 缓存目录结构

bilibili 客户端缓存路径可在客户端「设置 → 离线设置」中查看，默认为 `C:\Users\用户名\Videos\bilibili\`。其下每个以 `itemId` 命名的子目录对应一个已缓存的分 P：

```
bilibili\
└── 39556222659\
    ├── videoInfo.json          # 视频元数据（缺失时回退读取 .videoInfo）
    ├── .videoInfo
    ├── .playurl                # 播放地址信息（完整性校验依据之一）
    ├── 39556222659-1-30064.m4s # 视频流（30064 / 30080 均为视频）
    └── 39556222659-1-30280.m4s # 音频流（30280）
```

## 构建

直接使用 Go 编译：

```bat
go build -o BiliCacheMerge.exe .
```

或使用仓库内的跨平台构建脚本（产物名取自 `env` 文件，当前为 `BiliCacheMerge.exe`）：

| 脚本 | 目标平台 |
|---|---|
| `build_windowsx64.bat` | Windows amd64 |
| `build_windowsarm64.bat` | Windows arm64 |
| `build_linuxx64.bat` | Linux amd64 |
| `build_linuxarm64.bat` | Linux arm64 |

## 配置

复制并修改程序同目录下的 `config.ini`（也可通过 `-config` 指定其他路径）：

```ini
; bilibili 缓存根目录（其下每个子目录为一个已缓存视频）
[cache]
dir = D:\Binai\Videos\bilibili

; MP4 输出目录（必须已存在，否则程序提示并退出；其下按合集标题自动建子目录）
[output]
dir = D:\Binai\Downloads\BiliCacheMergeOutput

; ffmpeg 可执行文件路径（也支持 ffmpeg.bat 这类包装脚本）
[ffmpeg]
path = C:\ProFiles\cmd\ffmpeg.bat
```

启动时会依次检查缓存目录、输出目录、ffmpeg 是否存在，任一不可用都会提示错误，按任意键退出。

## 使用

```bat
BiliCacheMerge.exe -config config.ini
```

运行流程：

1. 校验配置中的三个路径；
2. 扫描缓存，列出所有合集的序号、`groupId`、`groupTitle`：

   ```
   已发现的缓存合集：
     1) BV12NTi6cEVV  【英语听力暴涨】108集《西游记》英文版动画！...
     2) 5254035       【循环歌单】遇到喜欢的歌不妨循环一下
   ```

3. 输入**序号**或 **groupId** 选择要处理的合集；
4. 程序逐个处理该合集的分 P，输出形如：

   ```
   [完成] 001_【石猴出世】 The Monkey.mp4
   [已存在] 002_....mp4
   [不完整] p=15 缺失: 音频流 m4s
   ```

5. 打印汇总（成功 / 已存在 / 不完整 / 失败），按任意键退出。

## 输出命名规则

- 输出路径：`{输出目录}/{合集标题}/{分P文件名}.mp4`
- 文件名：`p` 补零至 3 位，与标题用下划线连接，例如 `p=1` → `001_标题.mp4`
- 标题开头与 `p` 相同的数字前缀会被去除（避免出现 `019_019 xxx`）
- 标题最长 32 个字符（按 rune 截断，不会截断中文产生乱码）
- Windows 文件名非法字符（`\ / : * ? " < > |` 等）会被替换为 `_`
- 目标文件已存在时自动跳过，可安全地增量重跑

## 说明

- m4s 开头有 9 个占位字节，程序会先剥离到临时文件再交给 ffmpeg，合并结束自动清理临时文件。
- 音频直接复制源 AAC 音轨（`-c:a copy`）。早期版本曾转码为 AC3，会导致 **Android 端有画面无声音**，现已改为复制。
- 重新生成的 MP4 若与旧文件同名会被跳过；需要重做请先删除旧的输出文件。
