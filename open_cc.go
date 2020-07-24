package opencc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	occ "github.com/ApesPlan/prefixtree-OpenCC"
)

// s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, hk2t, t2hk, t2jp, jp2t, tw2t

// s2t ==> Simplified Chinese to Traditional Chinese 簡體到繁體
// t2s ==> Traditional Chinese to Simplified Chinese 繁體到簡體
// s2tw ==> Simplified Chinese to Traditional Chinese (Taiwan Standard) 簡體到臺灣正體
// tw2s ==> Traditional Chinese (Taiwan Standard) to Simplified Chinese 臺灣正體到簡體
// s2hk ==> Simplified Chinese to Traditional Chinese (Hong Kong variant) 簡體到香港繁體
// hk2s ==> Traditional Chinese (Hong Kong variant) to Simplified Chinese 香港繁體到簡體
// s2twp ==> Simplified Chinese to Traditional Chinese (Taiwan Standard) with Taiwanese idiom 簡體到繁體（臺灣正體標準）並轉換爲臺灣常用詞彙
// tw2sp ==> Traditional Chinese (Taiwan Standard) to Simplified Chinese with Mainland Chinese idiom 繁體（臺灣正體標準）到簡體並轉換爲中國大陸常用詞彙
// t2tw ==> Traditional Chinese (OpenCC Standard) to Taiwan Standard 繁體（OpenCC 標準）到臺灣正體
// hk2t ==> Traditional Chinese (Hong Kong variant) to Traditional Chinese 香港繁體到繁體（OpenCC 標準）
// t2hk ==> Traditional Chinese (OpenCC Standard) to Hong Kong variant 繁體（OpenCC 標準）到香港繁體
// t2jp ==> Traditional Chinese Characters (Kyūjitai) to New Japanese Kanji (Shinjitai) 繁體（OpenCC 標準，舊字體）到日文新字體
// jp2t ==> New Japanese Kanji (Shinjitai) to Traditional Chinese Characters (Kyūjitai) 日文新字體到繁體（OpenCC 標準，舊字體）
// tw2t ==> Traditional Chinese (Taiwan standard) to Traditional Chinese 臺灣正體到繁體（OpenCC 標準）

var (
	appDir               = defaultDir()
	supportedConversions = "s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, hk2t, t2hk, t2jp, jp2t, tw2t"
	configDir            = "config"
	dictDir              = "dictionary"
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
)

func defaultDir() string {
	_, p, _, _ := runtime.Caller(1)
	return path.Dir(p)
}

// Group holds a sequence of dicts
type Group struct {
	Files []string
	Dicts []*occ.Dict
}

func (g *Group) String() string {
	return fmt.Sprintf("%+v", g.Files)
}

// OpenCC contains the converter
type OpenCC struct {
	Conversion  string   // 14种转换方式 简写 s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, hk2t, t2hk, t2jp, jp2t, tw2t
	Description string   // 选择的config/*.json文件中的方式名称 eg: s2t.json中 "name": "Simplified Chinese to Traditional Chinese",
	DictChains  []*Group // 解析config/*.json文件中的group 将字典文件名称放入切片中 eg: s2t.json中 STPhrases.txt STCharacters.txt
}

// New construct an instance of OpenCC.
// Supported conversions: s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, hk2t, t2hk, t2jp, jp2t, tw2t
func New(conversion string) (*OpenCC, error) {
	if strings.TrimSpace(conversion) == "" {
		return nil, fmt.Errorf("Please select a conversion mode: %s", supportedConversions)
	}
	if _, has := conversions[conversion]; !has {
		return nil, fmt.Errorf("%s The conversion mode does not exist", conversion)
	}
	cc := &OpenCC{Conversion: conversion}
	err := cc.initDict()
	if err != nil {
		return nil, err
	}
	return cc, nil
}

// 解析字典文件到group切片中
func (cc *OpenCC) initDict() error {
	configFile := filepath.Join(appDir, configDir, cc.Conversion+".json")
	body, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	var m interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	config := m.(map[string]interface{})
	name, has := config["name"]
	if !has {
		return fmt.Errorf("name not found in %s", configFile)
	}
	cc.Description = name.(string)
	chain, has := config["conversion_chain"]
	if !has {
		return fmt.Errorf("conversion_chain not found in %s", configFile)
	}
	if dictChain, ok := chain.([]interface{}); ok {
		for _, v := range dictChain {
			if d, ok := v.(map[string]interface{}); ok {
				if gdict, has := d["dict"]; has {
					if dict, is := gdict.(map[string]interface{}); is {
						group, err := cc.addDictChain(dict) // 获取json中的配置
						if err != nil {
							return err
						}
						cc.DictChains = append(cc.DictChains, group)
					}
				} else {
					return fmt.Errorf("should have dict inside conversion_chain")
				}
			} else {
				return fmt.Errorf("should be map inside conversion_chain")
			}
		}
	} else {
		return fmt.Errorf("format %+v not correct in %s",
			reflect.TypeOf(dictChain), configFile)
	}
	return nil
}

// 递归解析*.json中group字段的值，即：dictionary/*.txt的文件名到group切片中
// OpenCC/data/config/*.json 和 OpenCC/data/dictionary/*.txt https://github.com/BYVoid/OpenCC
// OpenCC/data/config/*.json文件中 默认匹配的是.ocd2文件 （"type": "ocd2", "file": "TSPhrases.ocd2"），全部替换为txt
func (cc *OpenCC) addDictChain(d map[string]interface{}) (*Group, error) {
	t, has := d["type"]
	if !has {
		return nil, fmt.Errorf("type not found in %+v", d)
	}
	if typ, ok := t.(string); ok {
		ret := &Group{}
		switch typ {
		case "group":
			ds, has := d["dicts"]
			if !has {
				return nil, fmt.Errorf("no dicts field found")
			}
			dicts, is := ds.([]interface{})
			if !is {
				return nil, fmt.Errorf("dicts field invalid")
			}

			for _, dict := range dicts {
				if d, is := dict.(map[string]interface{}); is {
					group, err := cc.addDictChain(d)
					if err != nil {
						return nil, err
					}
					ret.Files = append(ret.Files, group.Files...)
					ret.Dicts = append(ret.Dicts, group.Dicts...)
				} else {
					return nil, fmt.Errorf("dicts items invalid")
				}
			}
		case "txt":
			file, has := d["file"]
			if !has {
				return nil, fmt.Errorf("no file field found")
			}
			daDict, err := occ.BuildFromFile(filepath.Join(appDir, dictDir, file.(string))) // 获取txt中数据
			if err != nil {
				return nil, err
			}
			ret.Files = append(ret.Files, file.(string))
			ret.Dicts = append(ret.Dicts, daDict)
		default:
			return nil, fmt.Errorf("type should be txt or group")
		}
		return ret, nil
	}
	return nil, fmt.Errorf("type should be string")
}

// Convert string from Simplified Chinese to Traditional Chinese or vice versa
// 将字符串从简体中文转换为繁体中文，反之亦然
// 外部调用的文本转换方法
func (cc *OpenCC) Convert(in string) (string, error) {
	var token string
	for _, group := range cc.DictChains {
		r := []rune(in)
		var tokens []string
		for i := 0; i < len(r); {
			s := r[i:]
			max := 0
			for _, dict := range group.Dicts {
				ret, err := dict.PrefixMatch(string(s))
				if err != nil {
					return "", err
				}
				if len(ret) > 0 {
					o := ""
					for k, v := range ret {
						if len(k) > max {
							max = len(k)
							token = v[0]
							o = k
						}
					}
					i += len([]rune(o))
					break
				}
			}
			if max == 0 { //no match
				token = string(r[i])
				i++
			}
			tokens = append(tokens, token)
		}
		in = strings.Join(tokens, "")
	}
	return in, nil
}
