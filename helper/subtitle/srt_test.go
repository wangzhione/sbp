package subtitle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	ctx := context.Background()
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
