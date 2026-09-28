package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// flexibleID 兼容 groupId 的两种取值：字符串（合集 bvid）或纯数字（单稿）。
type flexibleID string

// UnmarshalJSON 同时接受 JSON 字符串与数字。
func (f *flexibleID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*f = ""
		return nil
	}
	*f = flexibleID(strings.Trim(string(data), `"`))
	return nil
}

// errNotCacheDir 表示目录内不含任何元数据文件，不是缓存目录。
var errNotCacheDir = errors.New("非缓存目录")

// m4s 流文件后缀（fileId：30064/30080=视频，30280=音频）。
// TODO: 后续在此补充按文件大小判定音视频的兜底逻辑。
var (
	videoStreamSuffixes = []string{"-30064.m4s", "-30080.m4s"}
	audioStreamSuffixes = []string{"-30280.m4s"}
)

// VideoInfo 对应 videoInfo.json / .videoInfo 的扁平结构，仅保留用到的字段。
type VideoInfo struct {
	GroupID    flexibleID `json:"groupId"`
	GroupTitle string     `json:"groupTitle"`
	Aid        int64      `json:"aid"`
	Cid        int64      `json:"cid"`
	Bvid       string     `json:"bvid"`
	P          int        `json:"p"`
	Title      string     `json:"title"`
	ItemID     int64      `json:"itemId"`
}

// Item 一个已缓存的分 P。
type Item struct {
	Info     VideoInfo
	Dir      string
	VideoM4S string
	AudioM4S string
	Missing  []string
}

// Complete 报告该项缓存是否齐全。
func (it Item) Complete() bool { return len(it.Missing) == 0 }

// Group 同一合集（groupId）下的所有分 P。
type Group struct {
	ID    string
	Title string
	Items []Item
}

// ScanCache 扫描缓存根目录下的一级子目录，按 groupId 分组返回。
func ScanCache(root string) ([]*Group, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("读取缓存目录失败: %w", err)
	}

	groups := map[string]*Group{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root, e.Name())
		info, err := loadVideoInfo(dir)
		if err != nil {
			if !errors.Is(err, errNotCacheDir) {
				fmt.Fprintf(os.Stderr, "警告: 跳过 %s: %v\n", dir, err)
			}
			continue
		}

		item := buildItem(dir, info)
		id := string(info.GroupID)
		g := groups[id]
		if g == nil {
			g = &Group{ID: id, Title: info.GroupTitle}
			groups[id] = g
		}
		g.Items = append(g.Items, item)
	}

	result := make([]*Group, 0, len(groups))
	for _, g := range groups {
		sort.Slice(g.Items, func(i, j int) bool { return g.Items[i].Info.P < g.Items[j].Info.P })
		result = append(result, g)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

// loadVideoInfo 优先读取 videoInfo.json，缺失时回退到 .videoInfo。
func loadVideoInfo(dir string) (VideoInfo, error) {
	var info VideoInfo
	data, err := os.ReadFile(filepath.Join(dir, "videoInfo.json"))
	if err != nil {
		data, err = os.ReadFile(filepath.Join(dir, ".videoInfo"))
		if err != nil {
			return info, errNotCacheDir
		}
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return info, fmt.Errorf("解析元数据失败: %w", err)
	}
	return info, nil
}

// buildItem 收集该分 P 的文件并做完整性校验。
func buildItem(dir string, info VideoInfo) Item {
	item := Item{Info: info, Dir: dir}

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		name := e.Name()
		switch {
		case hasAnySuffix(name, videoStreamSuffixes):
			item.VideoM4S = filepath.Join(dir, name)
		case hasAnySuffix(name, audioStreamSuffixes):
			item.AudioM4S = filepath.Join(dir, name)
		}
	}

	if _, err := os.Stat(filepath.Join(dir, ".playurl")); err != nil {
		item.Missing = append(item.Missing, ".playurl")
	}
	if item.VideoM4S == "" {
		item.Missing = append(item.Missing, "视频流 m4s")
	}
	if item.AudioM4S == "" {
		item.Missing = append(item.Missing, "音频流 m4s")
	}
	return item
}

// hasAnySuffix 判断 name 是否以 suffixes 中任意一项结尾。
func hasAnySuffix(name string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}
