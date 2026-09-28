package main

import (
	"strings"
	"testing"
)

func TestOutputFileName(t *testing.T) {
	cases := []struct {
		p     int
		title string
		want  string
	}{
		{19, "019 【邪恶的计划】  An Evil Plan", "019_【邪恶的计划】  An Evil Plan.mp4"},
		{1, "001 【石猴出世】 The Monkey", "001_【石猴出世】 The Monkey.mp4"},
		{2, "第 2 集：开始", "002_第 2 集：开始.mp4"},
		{3, "合集珍藏 | 看动画", "003_合集珍藏 _ 看动画.mp4"},
		{4, "004", "004.mp4"},
	}
	for _, c := range cases {
		if got := OutputFileName(c.p, c.title); got != c.want {
			t.Errorf("OutputFileName(%d, %q) = %q, want %q", c.p, c.title, got, c.want)
		}
	}
}

func TestStripIndexPrefix(t *testing.T) {
	// 前缀数字与 p 不同则保留
	if got := StripIndexPrefix("001 第一集", 2); got != "001 第一集" {
		t.Errorf("前缀与 p 不同时应保留，得到 %q", got)
	}
	// 前缀非纯数字则保留
	if got := StripIndexPrefix("第 1 集", 1); got != "第 1 集" {
		t.Errorf("非数字前缀应保留，得到 %q", got)
	}
}

func TestTruncate(t *testing.T) {
	got := Truncate(strings.Repeat("中", 40), titleMaxRunes)
	if n := len([]rune(got)); n != titleMaxRunes {
		t.Errorf("截断后长度 = %d, want %d", n, titleMaxRunes)
	}
}

func TestSanitize(t *testing.T) {
	if got, want := Sanitize(`a<b>c:d"e/f\g|h?i*j`), "a_b_c_d_e_f_g_h_i_j"; got != want {
		t.Errorf("Sanitize() = %q, want %q", got, want)
	}
	if got, want := Sanitize("abc. "), "abc"; got != want {
		t.Errorf("结尾点空格应被去除，得到 %q", got)
	}
}
