package locale

import (
	"regexp"
	"strings"
)

// Parse :
type Parse interface {
	Lookup(string) string
}

// Locale :
type Locale struct {
	Default string                      // デフォルト言語
	Langs   map[string][]string         // ["en"][]string{"en", "en-US", ...}
	Ext     map[string]string           //
	match   map[string][]*regexp.Regexp //
}

// Lookup : Accept-Language のデータを解析し、言語判定を行う
func (locale *Locale) Lookup(language string) string {
	// 言語未設定の場合、デフォルト言語設定を返却する
	if locale.Langs == nil {
		return locale.Default
	}

	// Accept-Languageのデータを解析する
	// ex) ar-DZ,zh;q=0.8,ja;q=0.6,en-US;q=0.4,en;q=0.2
	for _, lang := range strings.Split(language, ",") {
		// q=...では判定しない。リストの並びで優先順位を選定する
		if idx := strings.Index(lang, ";"); idx != -1 {
			lang = lang[:idx]
		}
		// ex) lang = ar-DZ, zh, ja, en-US, en の順に処理
		// 正規表現を使用し、該当する言語を検索する
		for langname, exps := range locale.match {
			for _, exp := range exps {
				// 正規表現に該当しない場合、次の言語へ
				if !exp.Copy().MatchString(lang) {
					continue
				}
				// ex) map["en"][]{"^en\-US$", "^en\-.*$"} に使用しているキー名の "en" を返却する
				if locale.Ext == nil {
					return langname
				}
				if ext, ok := locale.Ext[langname]; ok {
					return ext
				}
				return langname
			}
		}
	}
	return locale.Default
}

// CreateLocale : 言語判定を行うために必要な情報を構築する
func (locale *Locale) CreateLocale() (Parse, error) {
	if locale.match != nil {
		return locale, nil
	}

	locale.match = make(map[string][]*regexp.Regexp)
	// 設定されている言語数分、正規表現オブジェクトを生成する
	for langname, langs := range locale.Langs {
		for _, lang := range langs {
			// ex) en-* => en\-.* 正規表現文字列へ置き換える
			lang = strings.NewReplacer(
				"-", "\\-",
				"*", ".*",
			).Replace(lang)
			// ex) ^en\-.*$ 正規表現オブジェクトを生成
			exp, err := regexp.Compile("^" + lang + "$")
			if err != nil {
				return nil, err
			}
			// 生成した正規表現オブジェクトを登録する
			locale.match[langname] = append(locale.match[langname], exp)
			// ex) map["ja"][]{"^ja$"}
			// ex) map["en"][]{"^en$", "^en\-US$", "^en\-.*"}
		}
	}

	return locale, nil
}
