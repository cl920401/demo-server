package redis

import (
	"testing"
)

func TestDefault(t *testing.T) {
	conn := Default()
	defer conn.Close()
	t.Log(conn.SetExString("aaa", "11111", 11111))
	t.Log(conn.Get("aaa"))
}
