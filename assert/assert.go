package assert

func Assert(condition bool, msg string) {
    if !condition {
        panic("Assertion Failure: " + msg);
    }
}
