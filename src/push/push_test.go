package push

import (
	"testing"
)

func Test_push(t *testing.T) {
	PushUnicast("xigua", "这是一条测试消息", "测试title", "测试", "iOS")
	PushUnicast("xigua", "这是一条测试消息", "测试title", "测试", "Android")
}
