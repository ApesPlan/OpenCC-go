# OpenCC-go
golang 简体繁体中文互转


## go mod demo

(```)
    package main

    import (
        "fmt"

        ccgo "github.com/ApesPlan/OpenCC-go"
    )

    func main() {
        // 简体转繁体
        str := "刘德华"
        s2t, err := ccgo.New("s2t")
        if err != nil {
            fmt.Printf("Error: %s\n", err)
        }
        out, err := s2t.Convert(str)
        if err != nil {
            fmt.Printf("Error: %s\n", err)
        }
        fmt.Println(out)

        // 繁体转简体
        // str := "劉德華"
        // t2s, err := ccgo.New("t2s")
        // if err != nil {
        // 	fmt.Printf("Error: %s\n", err)
        // }
        // out, err := t2s.Convert(str)
        // if err != nil {
        // 	fmt.Printf("Error: %s\n", err)
        // }
        // fmt.Println(out)
    }
(```)

### 使用包的时候，https://github.com/BYVoid/OpenCC有最新的简体繁体互转的配置文件和字典文件
资源路径 OpenCC/data/config/*.json 和 OpenCC/data/dictionary/*.txt 
可根据实际情况替换本包的 config/*.json 和 dictionary/*.txt 

OpenCC/data/config/*.json文件中 默认匹配的是.ocd2文件 （"type": "ocd2", "file": "TSPhrases.ocd2"），全部替换为txt即可

### 现目前支持14种
s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, hk2t, t2hk, t2jp, jp2t, tw2t

1. s2t ==> Simplified Chinese to Traditional Chinese 簡體到繁體
2. t2s ==> Traditional Chinese to Simplified Chinese 繁體到簡體
3. s2tw ==> Simplified Chinese to Traditional Chinese (Taiwan Standard) 簡體到臺灣正體
4. tw2s ==> Traditional Chinese (Taiwan Standard) to Simplified Chinese 臺灣正體到簡體
5. s2hk ==> Simplified Chinese to Traditional Chinese (Hong Kong variant) 簡體到香港繁體
6. hk2s ==> Traditional Chinese (Hong Kong variant) to Simplified Chinese 香港繁體到簡體
7. s2twp ==> Simplified Chinese to Traditional Chinese (Taiwan Standard) with Taiwanese idiom 簡體到繁體（臺灣正體標準）並轉換爲臺灣常用詞彙
8. tw2sp ==> Traditional Chinese (Taiwan Standard) to Simplified Chinese with Mainland Chinese idiom 繁體（臺灣正體標準）到簡體並轉換爲中國大陸常用詞彙
9. t2tw ==> Traditional Chinese (OpenCC Standard) to Taiwan Standard 繁體（OpenCC 標準）到臺灣正體
10. hk2t ==> Traditional Chinese (Hong Kong variant) to Traditional Chinese 香港繁體到繁體（OpenCC 標準）
11. t2hk ==> Traditional Chinese (OpenCC Standard) to Hong Kong variant 繁體（OpenCC 標準）到香港繁體
12. t2jp ==> Traditional Chinese Characters (Kyūjitai) to New Japanese Kanji (Shinjitai) 繁體（OpenCC 標準，舊字體）到日文新字體
13. jp2t ==> New Japanese Kanji (Shinjitai) to Traditional Chinese Characters (Kyūjitai) 日文新字體到繁體（OpenCC 標準，舊字體）
14. tw2t ==> Traditional Chinese (Taiwan standard) to Traditional Chinese 臺灣正體到繁體（OpenCC 標準）


### 如果有新添加的 需要到源码包open_cc.go文件中修改如下常量即可：
(```)
    supportedConversions = "s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, hk2t, t2hk, t2jp, jp2t, tw2t"
(```)

(```)
    conversions          = map[string]struct{}{
        "s2t":   {},
        "t2s":   {},
        "s2tw":  {},
        "tw2s":  {},
        "s2hk":  {},
        "hk2s":  {},
        "s2twp": {},
        "tw2sp": {},
        "t2tw":  {},
        "hk2t":  {},
        "t2hk":  {},
        "t2jp":  {},
        "jp2t":  {},
        "tw2t":  {},
	}
(```)
