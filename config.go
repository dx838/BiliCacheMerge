package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

// Config 运行配置，来自 ini 文件。
type Config struct {
	CacheDir   string
	OutputDir  string
	FFmpegPath string
}

// LoadConfig 读取并校验 ini 配置。
func LoadConfig(path string) (*Config, error) {
	f, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
	}

	cfg := &Config{
		CacheDir:   f.Section("cache").Key("dir").String(),
		OutputDir:  f.Section("output").Key("dir").String(),
		FFmpegPath: f.Section("ffmpeg").Key("path").String(),
	}

	if cfg.CacheDir == "" {
		return nil, fmt.Errorf("配置项 [cache] dir 为空")
	}
	if cfg.OutputDir == "" {
		return nil, fmt.Errorf("配置项 [output] dir 为空")
	}
	if cfg.FFmpegPath == "" {
		return nil, fmt.Errorf("配置项 [ffmpeg] path 为空")
	}
	if _, err := os.Stat(cfg.CacheDir); err != nil {
		return nil, fmt.Errorf("缓存目录不可用: %w", err)
	}
	if _, err := os.Stat(cfg.OutputDir); err != nil {
		return nil, fmt.Errorf("输出目录不存在: %s", cfg.OutputDir)
	}
	if _, err := os.Stat(cfg.FFmpegPath); err != nil {
		return nil, fmt.Errorf("ffmpeg 路径不可用: %w", err)
	}
	return cfg, nil
}
