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
		// Accept-Language --> "ar-DZ,ja;q=0.8,id;q=0.6,ug;q=0.4,ky;q=0.2"
		lang := parse.Lookup(r.Header.Get("Accept-Language"))
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		// ja 
		w.Write([]byte(lang))
	})
	http.ListenAndServe(":8080", nil)
}
```

## 言語ファイルの取り扱い

言語ファイルを用いて、言語データを取り扱うことができます。

### 言語ファイル置き場
`config/locales` ディレクトリ配下に、言語ファイルを置いていることを例として説明します。

```
config
  `-- locales
        +-- ja.json <-- 言語ファイル {"name":"ja:日本語"}
        +-- en.json <-- 言語ファイル {"name":"en:English"}
        `-- dirname
              +-- ja.json <-- 言語ファイル {"name":"dirname/ja:日本語"}
              `-- en.json <-- 言語ファイル {"name":"dirname/en:English"}
```

上記ディレクトリ構造の *.json ファイルを `LocaleDir` パラメータで指定して、読み込みます。

### サンプル
```go
package main

import (
	"fmt"

	"github.com/ochipin/locale"
)

func main() {
	// 判定する言語情報を設定
	var l = &locale.Locale{
		Default: ".ja",
		Langs: map[string][]string{
			"ja": []string{"ja"},
			"en": []string{"en", "en-*"},
			"zh": []string{"zh"},
		},
		// 言語ファイルが置いてあるディレクトリパスを指定。
		// config/locales 配下にある .json ファイルのみを対象にする
		LocaleDir: "config/locales",
	}

	parse, err := l.CreateLocale()
	if err != nil {
		panic(err)
	}

	// map[name:ja:日本語]
	fmt.Println(parse.Locale("ja"))
	// map[name:dirname/ja:日本語]
	fmt.Println(parse.Locale("dirname/ja"))
	// map[name:en:English]
	fmt.Println(parse.Locale("en"))
	// map[name:dirname/en:English]
	fmt.Println(parse.Locale("dirname/en"))

	// 言語をマージする。 ja + dirname/ja を結合したデータを取り扱う
	fmt.Println(locale.Merge(parse.Locale("ja"), parse.Locale("dirname/ja")))

	// Get関数を用いることで手軽にデータ取り出すことが可能。
	data := parse.Locale("ja")
	/*
	{
		"index": {
			"app": {
				"name": "appname!"
			}
		}
	}
	 */
	// appname! を出力
	fmt.Println(data.Get("index.app.name"))
}
```