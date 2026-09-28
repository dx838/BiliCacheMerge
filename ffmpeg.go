package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// m4sHeaderBytes 为缓存文件头部需要剥离的字节数（9 个 ASCII '0'）。
const m4sHeaderBytes = 9

// MergeM4S 剥离 m4s 头部后，调用 ffmpeg 将视频流与音频流合并为 MP4。
func MergeM4S(ffmpegPath string, item Item, outputPath string) error {
	tmpDir, err := os.MkdirTemp("", "bilicachemerge-")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	videoTmp := filepath.Join(tmpDir, "video.m4s")
	audioTmp := filepath.Join(tmpDir, "audio.m4s")
	if err := stripHeader(item.VideoM4S, videoTmp); err != nil {
		return fmt.Errorf("处理视频流失败: %w", err)
	}
	if err := stripHeader(item.AudioM4S, audioTmp); err != nil {
		return fmt.Errorf("处理音频流失败: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	args := []string{"-y", "-i", videoTmp, "-i", audioTmp, "-c:v", "copy", "-c:a", "copy", outputPath}
	cmd := ffmpegCommand(ffmpegPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg 合并失败: %w\n%s", err, stderr.String())
	}
	return nil
}

// stripHeader 跳过文件头部 m4sHeaderBytes 个字节，写入目标文件。
func stripHeader(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if _, err := in.Seek(m4sHeaderBytes, io.SeekStart); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// ffmpegCommand 构造 ffmpeg 命令；Windows 下 .bat/.cmd 需经 cmd /c 执行。
func ffmpegCommand(ffmpegPath string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		switch strings.ToLower(filepath.Ext(ffmpegPath)) {
		case ".bat", ".cmd":
			return exec.Command("cmd", append([]string{"/c", ffmpegPath}, args...)...)
		}
	}
	return exec.Command(ffmpegPath, args...)
}
