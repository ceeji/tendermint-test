package kvstore

import (
	"testing"
)

func Test_mergeKey(t *testing.T) {
	merged := mergeKey("abc", []byte("123"))

	if len(merged) != 7 {
		t.Error("length should be 7")
	}
	if string(merged[1:]) != "abc123" {
		t.Error("wrong key")
	}
}
