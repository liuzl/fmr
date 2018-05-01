package funcs

import (
	"testing"
)

func TestSum(t *testing.T) {
	t.Log(sum("34", "987"))
	t.Log(div("399", "987"))
	t.Log(div("3e9", "987"))
	t.Log(div("abc", "987"))
	t.Log(pow("2.1", "7.9"))
	t.Log(neg("-2.1e100"))
	t.Log(odd("24"))
	t.Log(even("24"))
}
