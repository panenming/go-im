package exec

import (
	"strings"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	err := ExecTimeout(1*time.Second, "sleep", "10")
	// 这个判断不能通过，linux kill的返回值和windows不同
	// windows 需要判断 !strings.Contains(err.Error(), "exit status 1")
	// 这个最好可以自定义kill的返回值？？？
	if err != nil && !strings.Contains(err.Error(), "signal: killed") {
		t.Fatal(err)
	}
}
