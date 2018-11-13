package locale

import (
	"fmt"
	"testing"
)

func Test__NEW_CREATE_LOOKUP(t *testing.T) {
	var locale = &Locale{
		Default: "ja",
		Langs: map[string][]string{
			"ja": []string{"ja"},
			"en": []string{"en", "en-*"},
			"zh": []string{"zh"},
		},
		Ext: map[string]string{
			"ja": ".ja",
			"zh": ".zh",
		},
	}

	parse, err := locale.CreateLocale()
	if err != nil {
		t.Fatal(err)
	}
	// 2回コールしても、問題はない
	parse, err = locale.CreateLocale()

	// .zh
	// ar-DZは該当しないため、無視される。次のzhは、Langs["zh"][]string{"zh"}に該当するため、zh が返却される
	// ただし、Ext["zh"]".zh" にも該当するため、返却値は、Extで定義している .zh が返却される
	if lang := parse.Lookup("ar-DZ,zh;q=0.8,ja;q=0.6,en-US;q=0.4,en;q=0.2"); lang != ".zh" {
		t.Fatalf("%s Lookup Error", lang)
	}
	// en
	if lang := parse.Lookup("en,zh;q=0.8,ja;q=0.6,en-US;q=0.4,ar-DZ;q=0.2"); lang != "en" {
		t.Fatalf("%s Lookup Error", lang)
	}
	// en
	// ar-DZは該当しないため、無視される。次のen-USは、Langs["en"][]string{"en-*"}に該当するため、en が返却される
	if lang := parse.Lookup("ar-DZ,en-US;q=0.8,ja;q=0.6,zh;q=0.4,en;q=0.2"); lang != "en" {
		t.Fatalf("%s Lookup Error", lang)
	}
	// ja
	// すべての文字列に該当しないため、jaが返却される
	if lang := parse.Lookup("ar-DZ,ar-JO;q=0.8,id;q=0.6,ug;q=0.4,ky;q=0.2"); lang != "ja" {
		t.Fatalf("%s Lookup Error", lang)
	}
}

func Test__NEW_CREATE_DEFAULT(t *testing.T) {
	var locale = &Locale{
		Default: "ja",
	}

	parse, err := locale.CreateLocale()
	if err != nil {
		t.Fatal(err)
	}

	// ja
	if lang := parse.Lookup("ar-DZ,zh;q=0.8,ja;q=0.6,en-US;q=0.4,en;q=0.2"); lang != "ja" {
		t.Fatalf("%s Lookup Error", lang)
	}
}

func Test__NEW_CREATE_LOOKUP_NOEXT(t *testing.T) {
	var locale = &Locale{
		Default: "ja",
		Langs: map[string][]string{
			"ja": []string{"ja"},
			"en": []string{"en", "en-*"},
			"zh": []string{"zh"},
		},
	}

	parse, err := locale.CreateLocale()
	if err != nil {
		t.Fatal(err)
	}

	// .zh
	if lang := parse.Lookup("ar-DZ,zh;q=0.8,ja;q=0.6,en-US;q=0.4,en;q=0.2"); lang != "zh" {
		t.Fatalf("%s Lookup Error", lang)
	}
	// en
	if lang := parse.Lookup("en,zh;q=0.8,ja;q=0.6,en-US;q=0.4,ar-DZ;q=0.2"); lang != "en" {
		t.Fatalf("%s Lookup Error", lang)
	}
	// en
	if lang := parse.Lookup("ar-DZ,en-US;q=0.8,ja;q=0.6,zh;q=0.4,en;q=0.2"); lang != "en" {
		t.Fatalf("%s Lookup Error", lang)
	}
}

func Test__NEW_CREATE_ERROR(t *testing.T) {
	var locale = &Locale{
		Default: "ja",
		Langs: map[string][]string{
			"ja": []string{"ja"},
			"en": []string{"en", "d(^-^o"}, // デタラメな値は許されない
			"zh": []string{"zh"},
		},
		Ext: map[string]string{
			"ja": ".ja",
			"zh": ".zh",
		},
	}

	// エラーとなる
	_, err := locale.CreateLocale()
	if err == nil {
		t.Fatal("CreateLocale Error")
	}
	fmt.Println(err)
}
