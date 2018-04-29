package funcs

import (
	"testing"
)

func TestSum(t *testing.T) {
	t.Log(sum("34", "987"))
	t.Log(div("399", "987"))
	t.Log(div("3e9", "987"))
	t.Log(div("abc", "987"))
}
