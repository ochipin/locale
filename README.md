言語判定ライブラリ
===
HTTPヘッダのAccept-Languageの情報を元に、言語判定を行います。

## サンプル

```go
package main

import (
    "fmt"

    "github.com/ochipin/locale"
)

func main() {
    // 判定する言語情報を設定
	var locale = &locale.Locale{
        // デフォルトで使用される言語
        Default: ".ja",
        // Accept-Language の判定に使用する
		Langs: map[string][]string{
			"ja": []string{"ja"},
			"en": []string{"en", "en-US", "en-*"},
			"zh": []string{"zh"},
        },
		Ext: map[string]string{
			"ja": ".ja",
			"zh": ".zh",
		},
	}

	parse, err := locale.CreateLocale()
	if err != nil {
        panic(err)
	}

	// .zh
	fmt.Println(parse.Lookup("ar-DZ,zh;q=0.8,ja;q=0.6,en-US;q=0.4,en;q=0.2"))
	// en
	fmt.Println(parse.Lookup("en,zh;q=0.8,ja;q=0.6,en-US;q=0.4,ar-DZ;q=0.2"))
    // en
	fmt.Println(parse.Lookup("ar-DZ,en-US;q=0.8,ja;q=0.6,zh;q=0.4,en;q=0.2"))
	// .ja
	fmt.Println(parse.Lookup("ar-DZ,ar-JO;q=0.8,id;q=0.6,ug;q=0.4,ky;q=0.2"))
}
```

```go
package main

import (
    "net/http"

    "github.com/ochipin/locale"
)

func main() {
    // 判定する言語情報を設定
	var locale = &locale.Locale{
		Default: ".ja",
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
        panic(err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        lang := parse.Lookup(r.Header.Get("Accept-Language"))
        w.WriteHeader(200)
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(lang))
    })
    http.ListenAndServe(":8080", nil)
}
```

