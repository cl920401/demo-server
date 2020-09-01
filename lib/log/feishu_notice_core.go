package log

import (
	"fmt"
	"strings"

	"demo-server/lib/feishu"
	"go.uber.org/zap/zapcore"
)

type FeiShuNoticeCore struct {
	enc zapcore.Encoder
}

func NewFeiShuNoticeCore(enc zapcore.Encoder) zapcore.Core {
	return &FeiShuNoticeCore{
		enc: enc,
	}
}
func (c *FeiShuNoticeCore) Enabled(level zapcore.Level) bool {
	return level >= zapcore.ErrorLevel
}
func (c *FeiShuNoticeCore) With([]zapcore.Field) zapcore.Core { return c }
func (c *FeiShuNoticeCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}
func (c *FeiShuNoticeCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	// 不用实际写入任何东西，发送通知即可

	var title = fmt.Sprintf("%s %s:%d", strings.ToUpper(ent.Level.String()),
		ent.Caller.File, ent.Caller.Line)

	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	msg := buf.String()
	buf.Free()

	feishu.Send(title, msg)
	return nil
}
func (c *FeiShuNoticeCore) Sync() error { return nil }
