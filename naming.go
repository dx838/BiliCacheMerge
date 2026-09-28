package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const titleMaxRunes = 32

var (
	// 匹配 title 开头的数字序号（含前导 0），如 "019 "。
	leadingIndex = regexp.MustCompile(`^\s*0*(\d+)(\s+|$)`)
	// Windows 文件名非法字符及 ASCII 控制符。
	illegalChars = regexp.MustCompile(`[\\/:*?"<>|\x00-\x1f]`)
)

// StripIndexPrefix 去掉 title 开头与 p 相同的数字前缀。
// 例如 p=19、title="019 【邪恶的计划】  An Evil Plan" 得到 "【邪恶的计划】  An Evil Plan"。
func StripIndexPrefix(title string, p int) string {
	m := leadingIndex.FindStringSubmatch(title)
	if m == nil {
		return title
	}
	if n, err := strconv.Atoi(m[1]); err == nil && n == p {
		return title[len(m[0]):]
	}
	return title
}

// Sanitize 替换文件名非法字符，并去掉结尾的点和空格。
func Sanitize(name string) string {
	return strings.TrimRight(illegalChars.ReplaceAllString(name, "_"), ". ")
}

// Truncate 按字符（rune）截断到 max 个字符。
func Truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

// OutputFileName 生成输出文件名 "{p 补零 3 位}_{title}.mp4"。
func OutputFileName(p int, title string) string {
	title = Sanitize(Truncate(StripIndexPrefix(title, p), titleMaxRunes))
	if title == "" {
		return fmt.Sprintf("%03d.mp4", p)
	}
	return fmt.Sprintf("%03d_%s.mp4", p, title)
}
