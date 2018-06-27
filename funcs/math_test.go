package funcs

import (
	"testing"
)

func TestSum(t *testing.T) {
	t.Log(sum("34", "987"))
	t.Log(sub("34", "987"))
	t.Log(div("399", "987"))
	t.Log(div("3e9", "987"))
	t.Log(div("abc", "987"))
	t.Log(pow("2.1", "7.9"))
	t.Log(neg("-2.1e100"))
	t.Log(sum(neg("-2.1e100"), "1.1e99"))
	t.Log(odd("24"))
	t.Log(even("24"))
	t.Log(prime("100000223"))
	t.Log(prime("100000227"))
	t.Log(Call("nf.math.prime", "100000227"))
	t.Log(Call("nf.math.prime", "227"))
}
