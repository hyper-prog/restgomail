package main

/*  Smart JSON functions - Helper functions to handle JSON
    (C) 2021 Péter Deák (hyper80@gmail.com)
    License: GPLv2
*/

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SmartJSON struct {
	parsed map[string]interface{}
}

func parseSmartJSON(rawdata []byte) (SmartJSON, error) {
	s := SmartJSON{}
	s.parsed = make(map[string]interface{})
	err := json.Unmarshal(rawdata, &s.parsed)
	return s, err
}

func (sJSON SmartJSON) toFormattedString() (out string) {
	return "JSON => {\n" + jsonNodeToString(sJSON.parsed, "    ") + "}\n"
}

func jsonNodeToString(n map[string]interface{}, indent string) (out string) {
	out = ""
	for n, v := range n {
		out += indent + n + " => "
		if str, isStr := v.(string); isStr {
			out += str + "\n"
			continue
		}
		if flt, isFlt := v.(float64); isFlt {
			out += fmt.Sprintf("%f\n", flt)
			continue
		}
		if b, isB := v.(bool); isB {
			if b {
				out += "true\n"
			} else {
				out += "false\n"
			}
			continue
		}

		m, okm := v.(map[string]interface{})
		if okm {
			out += "{\n"
			out += jsonNodeToString(m, indent+"    ")
			out += indent + "}\n"
			continue
		}
		out += "?\n"
	}
	return
}

func (sJSON SmartJSON) getValueByPath(path string) (interface{}, string) {
	parts := strings.Split(path, "/")
	n := sJSON.parsed
	for i := 0; i < len(parts); i++ {
		v, ok := n[parts[i]]
		if !ok {
			return nil, "none"
		}
		if i == len(parts)-1 {
			if s, isStr := v.(string); isStr {
				return s, "string"
			}
			if f, isFlt := v.(float64); isFlt {
				return f, "float64"
			}
			if b, isBool := v.(bool); isBool {
				return b, "bool"
			}
			if m, isMap := v.(map[string]interface{}); isMap {
				return m, "map"
			}
		}
		if m, isMap := v.(map[string]interface{}); isMap {
			n = m
			continue
		}
		return nil, "none"
	}
	return nil, "none"
}

func (sJSON SmartJSON) getStringByPath(path string) (string, string) {
	val, typ := sJSON.getValueByPath(path)
	if str, isStr := val.(string); typ == "string" && isStr {
		return str, typ
	}
	return "", "none"
}

func (sJSON SmartJSON) getFloat64ByPath(path string) (float64, string) {
	val, typ := sJSON.getValueByPath(path)
	if f, isFlt := val.(float64); typ == "float64" && isFlt {
		return f, typ
	}
	return 0, "none"
}

func (sJSON SmartJSON) getBoolByPath(path string) (bool, string) {
	val, typ := sJSON.getValueByPath(path)
	if b, isBool := val.(bool); typ == "bool" && isBool {
		return b, typ
	}
	return false, "none"
}

func (sJSON SmartJSON) getMapByPath(path string) (map[string]interface{}, string) {
	val, typ := sJSON.getValueByPath(path)
	if m, isMap := val.(map[string]interface{}); typ == "map" && isMap {
		return m, typ
	}
	return nil, "none"
}

func (sJSON SmartJSON) getStringByPathWithDefault(path string, def string) string {
	val, typ := sJSON.getValueByPath(path)
	if str, isStr := val.(string); typ == "string" && isStr {
		return str
	}
	return def
}

func (sJSON SmartJSON) getFloat64ByPathWithDefault(path string, def float64) float64 {
	val, typ := sJSON.getValueByPath(path)
	if f, isFlt := val.(float64); typ == "float64" && isFlt {
		return f
	}
	return def
}

func (sJSON SmartJSON) getBoolByPathWithDefault(path string, def bool) bool {
	val, typ := sJSON.getValueByPath(path)
	if b, isBool := val.(bool); typ == "bool" && isBool {
		return b
	}
	return def
}
