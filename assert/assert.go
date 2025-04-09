package assert

import (
	"runtime"
	"strconv"
)

func Assert(condition bool, msg ...string) {
    if condition {
        return
    }
    message := "Assertion Failure: "
    for _, m := range msg {
        message += m
    }
    _, file, line, ok := runtime.Caller(1)
    if !ok {
        panic(message)
    }
    message += "\nFile: " + file + "\nLine: " + strconv.FormatInt(int64(line), 10)
    panic(message)
}
