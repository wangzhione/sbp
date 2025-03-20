package subtitle

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangzhione/sbp/util/chain"
)

var ctx = chain.Context()

// 测试 `DownloadSRT`
func TestDownloadSRT(t *testing.T) {
	// 定义模拟 SRT 文件内容
	mockSRT := strings.TrimSpace(`
1
00:00:01,500 --> 00:00:04,000
Hello, world!

2
00:00:05,000 --> 00:00:07,500
This is a test subtitle.
`)

	// 创建一个模拟 HTTP 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-subrip")
		_, _ = w.Write([]byte(mockSRT))
	}))
	defer server.Close() // 关闭服务器

	// 执行 `DownloadSRT`
	subtitles, err := DownloadSRT(ctx, server.URL)

	// 确保没有错误
	assert.NoError(t, err, "DownloadSRT should not return an error")
	assert.NotNil(t, subtitles, "Subtitles should not be nil")

	// 验证解析出的字幕条目
	assert.Equal(t, 2, len(subtitles.Items), "Expected 2 subtitle entries")

	// 验证第一条字幕
	assert.NotEmpty(t, subtitles.Items[0].Lines, "First subtitle should have at least one line")
	assert.NotEmpty(t, subtitles.Items[0].Lines[0].Items, "First subtitle line should have at least one item")

	item1 := subtitles.Items[0]
	assert.Equal(t, "Hello, world!", item1.Lines[0].Items[0].Text, "First subtitle text should match")
	assert.Equal(t, "1.5s", item1.StartAt.String(), "First subtitle start time should match")
	assert.Equal(t, "4s", item1.EndAt.String(), "First subtitle end time should match")

	// 验证第二条字幕
	assert.NotEmpty(t, subtitles.Items[1].Lines, "Second subtitle should have at least one line")
	assert.NotEmpty(t, subtitles.Items[1].Lines[0].Items, "Second subtitle line should have at least one item")

	item2 := subtitles.Items[1]
	assert.Equal(t, "This is a test subtitle.", item2.Lines[0].Items[0].Text, "Second subtitle text should match")
	assert.Equal(t, "5s", item2.StartAt.String(), "Second subtitle start time should match")
	assert.Equal(t, "7.5s", item2.EndAt.String(), "Second subtitle end time should match")
}

// 测试 LoadSRT 正常解析 SRT 文件
func TestLoadSRT(t *testing.T) {
	// 定义测试 SRT 内容
	mockSRT := `1
00:00:01,500 --> 00:00:04,000
Hello, world!

2
00:00:05,000 --> 00:00:07,500
This is a test subtitle.`

	// 创建临时 SRT 文件
	tempFile, err := os.CreateTemp("", "test_*.srt")
	assert.NoError(t, err, "Failed to create temporary SRT file")
	defer os.Remove(tempFile.Name()) // 测试完成后删除文件

	// 写入 SRT 内容
	_, err = tempFile.WriteString(mockSRT)
	assert.NoError(t, err, "Failed to write to temporary SRT file")

	// 关闭文件，使其可用于读取
	err = tempFile.Close()
	assert.NoError(t, err, "Failed to close temporary SRT file")

	// 执行 LoadSRT
	subtitles, err := LoadSRT(ctx, tempFile.Name())

	// 确保解析没有错误
	assert.NoError(t, err, "LoadSRT should not return an error")
	assert.NotNil(t, subtitles, "Subtitles should not be nil")

	// 验证字幕条目数量
	assert.Equal(t, 2, len(subtitles.Items), "Expected 2 subtitle entries")

	// 验证第一条字幕内容
	item1 := subtitles.Items[0]
	assert.Equal(t, "Hello, world!", item1.Lines[0].String(), "First subtitle text should match")
	assert.Equal(t, "Hello, world!", item1.String(), "First subtitle text should match")
	assert.Equal(t, "1.5s", item1.StartAt.String(), "First subtitle start time should match")
	assert.Equal(t, "4s", item1.EndAt.String(), "First subtitle end time should match")

	// 验证第二条字幕内容
	item2 := subtitles.Items[1]
	assert.Equal(t, "This is a test subtitle.", item2.Lines[0].String(), "Second subtitle text should match")
	assert.Equal(t, "This is a test subtitle.", item2.String(), "Second subtitle text should match")
	assert.Equal(t, "5s", item2.StartAt.String(), "Second subtitle start time should match")
	assert.Equal(t, "7.5s", item2.EndAt.String(), "Second subtitle end time should match")
}
