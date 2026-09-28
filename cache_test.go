package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestScanCache(t *testing.T) {
	root := t.TempDir()

	// 完整的一项：p=1，元数据为 videoInfo.json
	full := filepath.Join(root, "1001")
	mkdir(t, full)
	writeFile(t, filepath.Join(full, "videoInfo.json"),
		`{"groupId":"BVtest","groupTitle":"测试合集","itemId":1001,"aid":1,"cid":1001,"bvid":"BVtest","p":1,"title":"001 第一集"}`)
	writeFile(t, filepath.Join(full, ".playurl"), "x")
	writeFile(t, filepath.Join(full, "1001-1-30064.m4s"), "")
	writeFile(t, filepath.Join(full, "1001-1-30280.m4s"), "")

	// 缺音频的一项：p=2，元数据仅有 .videoInfo（验证回退）
	part := filepath.Join(root, "1002")
	mkdir(t, part)
	writeFile(t, filepath.Join(part, ".videoInfo"),
		`{"groupId":"BVtest","groupTitle":"测试合集","itemId":1002,"aid":1,"cid":1002,"bvid":"BVtest","p":2,"title":"002 第二集"}`)
	writeFile(t, filepath.Join(part, ".playurl"), "x")
	writeFile(t, filepath.Join(part, "1002-1-30080.m4s"), "")

	// 数字 groupId 的合集（UGC 单稿）
	num := filepath.Join(root, "2001")
	mkdir(t, num)
	writeFile(t, filepath.Join(num, "videoInfo.json"),
		`{"groupId":5254035,"groupTitle":"数字分组","itemId":2001,"aid":2,"cid":2001,"bvid":"BVnum","p":1,"title":"单曲"}`)
	writeFile(t, filepath.Join(num, ".playurl"), "x")
	writeFile(t, filepath.Join(num, "2001-1-30080.m4s"), "")
	writeFile(t, filepath.Join(num, "2001-1-30280.m4s"), "")

	// 非缓存目录（无元数据）应被忽略
	mkdir(t, filepath.Join(root, "junk"))

	groups, err := ScanCache(root)
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]*Group{}
	for _, g := range groups {
		byID[g.ID] = g
	}
	if len(byID) != 2 {
		t.Fatalf("分组数 = %d, want 2 (%v)", len(byID), byID)
	}

	g := byID["BVtest"]
	if g == nil {
		t.Fatal("未找到字符串 groupId 分组 BVtest")
	}
	if g.Title != "测试合集" {
		t.Fatalf("分组标题错误: %q", g.Title)
	}
	if len(g.Items) != 2 {
		t.Fatalf("分 P 数 = %d, want 2", len(g.Items))
	}
	// 结果按 p 升序
	if g.Items[0].Info.P != 1 || g.Items[1].Info.P != 2 {
		t.Fatalf("分 P 未按 p 升序: %d, %d", g.Items[0].Info.P, g.Items[1].Info.P)
	}
	if !g.Items[0].Complete() {
		t.Errorf("p=1 应完整，缺失: %v", g.Items[0].Missing)
	}
	if g.Items[1].Complete() {
		t.Errorf("p=2 应不完整")
	}
	if len(g.Items[1].Missing) != 1 {
		t.Errorf("p=2 缺失项数 = %d, want 1 (%v)", len(g.Items[1].Missing), g.Items[1].Missing)
	}

	// 数字 groupId 的合集也必须被识别
	ng := byID["5254035"]
	if ng == nil {
		t.Fatal("未找到数字 groupId 分组 5254035")
	}
	if len(ng.Items) != 1 || !ng.Items[0].Complete() {
		t.Fatalf("数字 groupId 分组内容异常: %+v", ng)
	}
}
