package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// stdin 全局复用同一个缓冲读取器，避免缓冲内容在多次读取间丢失。
var stdin = bufio.NewReader(os.Stdin)

func main() {
	configPath := flag.String("config", "config.ini", "配置文件路径")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		fatal(err)
	}

	groups, err := ScanCache(cfg.CacheDir)
	if err != nil {
		fatal(err)
	}
	if len(groups) == 0 {
		fmt.Println("未在缓存目录中发现任何视频")
		pressAnyKeyToExit()
		return
	}

	printGroups(groups)

	group, err := selectGroup(groups)
	if err != nil {
		fatal(err)
	}

	processGroup(cfg, group)

	pressAnyKeyToExit()
}

// fatal 打印错误，等待用户按任意键后退出。
func fatal(err error) {
	fmt.Fprintln(os.Stderr, "错误:", err)
	pressAnyKeyToExit()
	os.Exit(1)
}

// pressAnyKeyToExit 暂停，等待用户按任意键。
func pressAnyKeyToExit() {
	fmt.Print("\n按任意键退出...")
	_, _ = stdin.ReadByte()
}

// printGroups 输出已发现合集的序号、groupId 与 groupTitle。
func printGroups(groups []*Group) {
	fmt.Println("已发现的缓存合集：")
	for i, g := range groups {
		fmt.Printf("  %d) %s  %s\n", i+1, g.ID, g.Title)
	}
}

// selectGroup 读取用户输入，按序号或 groupId 选择要处理的合集。
func selectGroup(groups []*Group) (*Group, error) {
	fmt.Print("请输入序号或 groupId 选择要处理的合集: ")
	line, err := stdin.ReadString('\n')
	if err != nil && line == "" {
		return nil, fmt.Errorf("读取输入失败: %w", err)
	}

	input := strings.TrimSpace(line)
	if input == "" {
		return nil, fmt.Errorf("未输入任何内容")
	}
	if n, err := strconv.Atoi(input); err == nil {
		if n < 1 || n > len(groups) {
			return nil, fmt.Errorf("序号超出范围: %s", input)
		}
		return groups[n-1], nil
	}
	for _, g := range groups {
		if g.ID == input {
			return g, nil
		}
	}
	return nil, fmt.Errorf("未找到 groupId: %s", input)
}

// processGroup 处理指定合集下的所有分 P。
func processGroup(cfg *Config, g *Group) {
	fmt.Printf("\n合集: %s\n  groupId: %s  分 P: %d\n", g.Title, g.ID, len(g.Items))
	outDir := filepath.Join(cfg.OutputDir, Sanitize(g.Title))

	var okCount, skipCount, incompleteCount, failCount int
	for _, item := range g.Items {
		name := OutputFileName(item.Info.P, item.Info.Title)
		target := filepath.Join(outDir, name)

		if !item.Complete() {
			incompleteCount++
			fmt.Printf("  [不完整] p=%d 缺失: %s\n", item.Info.P, strings.Join(item.Missing, "、"))
			continue
		}
		if _, err := os.Stat(target); err == nil {
			skipCount++
			fmt.Printf("  [已存在] %s\n", name)
			continue
		}
		if err := MergeM4S(cfg.FFmpegPath, item, target); err != nil {
			failCount++
			fmt.Printf("  [失败] p=%d %v\n", item.Info.P, err)
			continue
		}
		okCount++
		fmt.Printf("  [完成] %s\n", name)
	}

	fmt.Printf("\n合计: 成功 %d, 已存在 %d, 不完整 %d, 失败 %d\n", okCount, skipCount, incompleteCount, failCount)
}
