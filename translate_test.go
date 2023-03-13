package deeplx

import (
	"testing"
)

func TestTranslate(t *testing.T) {
	result, err := Translate("Go是一种语言层面支持并发（Go最大的特色、天生支持并发）\n内置runtime、iiii支持垃圾回收（GC）、静态强类型，快速编译的语言", "auto", "en")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(result.Text)
}
