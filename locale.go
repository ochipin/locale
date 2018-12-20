package locale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Parse :
type Parse interface {
	Lookup(string) string
	Locale(string) map[string]interface{}
}

// Locale : 言語環境管理構造体
type Locale struct {
	Default   string                                 // デフォルト言語
	Langs     map[string][]string                    // ["en"][]string{"en", "en-US", ...}
	Ext       map[string]string                      //
	match     map[string][]*regexp.Regexp            //
	LocaleDir string                                 // 言語ファイル置き場
	locales   map[string]interface{}                 // 言語ファイル群
	Walk      func(string, os.FileInfo, error) error // 言語設定ファイル解析関数ポインタ
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

// Locale : 設定済みの言語設定情報を取得する
func (locale *Locale) Locale(name string) map[string]interface{} {
	if v, ok := locale.locales[name]; ok {
		if mapdata, ok := v.(map[string]interface{}); ok {
			return mapdata
		}
	}
	return nil
}

func (locale *Locale) setLocale() error {
	// 設定項目が未設定の場合は何もせず復帰する
	if locale.LocaleDir == "" {
		return nil
	}
	// 言語設定情報を格納するマップを初期化する
	locale.locales = make(map[string]interface{})
	if locale.Walk == nil {
		locale.Walk = locale.DefaultWalk
	}

	return filepath.Walk(locale.LocaleDir, locale.Walk)
}

// DefaultWalk : 言語設定ファイル群を読み込みデータに保持する関数
func (locale *Locale) DefaultWalk(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if f.IsDir() {
		return nil
	}

	// JSON ファイルか否かをチェックする
	idx := strings.LastIndex(path, ".")
	if idx == -1 {
		return nil
	}
	if path[idx:] != ".json" {
		return nil
	}

	// フルパスを整形する ex) /path/to/conf.ja.json => conf.ja
	var name = path
	if len(locale.LocaleDir)+1 < len(path[:idx]) {
		name = path[len(locale.LocaleDir)+1 : idx]
	}

	// 言語設定ファイルからデータを抽出する
	var data = make(map[string]interface{})
	buf, _ := ioutil.ReadFile(path)
	if len(buf) == 0 {
		buf = []byte("{}")
	}
	if err := json.Unmarshal(buf, &data); err != nil {
		return fmt.Errorf("%s: %v", name, err)
	}
	locale.locales[name] = data
	return nil
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
	// 言語設定ファイルをデータとして保持する
	if err := locale.setLocale(); err != nil {
		return nil, err
	}

	return locale, nil
}

// Merge : src に dst を追加した値を取得する
func Merge(src, dst map[string]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	if src != nil {
		merge(result, src)
	}
	if dst != nil {
		merge(result, dst)
	}
	// 引数に渡されたsrc,dst双方がnilの場合は、nilを返却する
	if src == nil && dst == nil {
		return nil
	}
	// マージしたデータを返却する
	return result
}

// src に dst をマージする
func merge(src map[string]interface{}, dst map[string]interface{}, keys ...string) {
	// マージするデータをループで全データ処理
	for key, value := range dst {
		keys = append(keys, key)
		if v, ok := value.(map[string]interface{}); ok {
			// map[app]などのデータがmapの場合、再帰する
			merge(src, v, keys...)
		} else {
			// mapの終端(map[app][key][name])へ辿り着いた時点で、データをマージ先と結合する
			set(src, value, keys)
		}
		// map[app][key][name]処理後、map[app][key][????]など別データがある可能性もあるため、
		// 参照先を map[app][key] までの値に戻す
		if len(keys) > 0 {
			keys = keys[:len(keys)-1]
		}
	}
}

// mapにデータを追加/上書きする
func set(src map[string]interface{}, dst interface{}, keys []string) {
	last := keys[len(keys)-1] // [app key name] の最後尾にある name のみを格納
	keys = keys[:len(keys)-1] // [app key]の2つの要素のみにする
	// [app key]の値のみを検証
	for _, key := range keys {
		if v, ok := src[key].(map[string]interface{}); ok {
			// src[key]がmapの場合、次の要素へ
			src = v
		} else {
			// src[key]が存在しない場合、mapを生成して次の要素へ
			src[key] = make(map[string]interface{})
			src = src[key].(map[string]interface{})
		}
	}
	src[last] = dst
}
