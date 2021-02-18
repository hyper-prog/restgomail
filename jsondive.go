package main

/*  JSON Dive functions
    (C) 2021 Péter Deák (hyper80@gmail.com)
    License: GPLv2
*/

import (
	"fmt"
	"strings"
)

func printParsedJSON(n map[string]interface{}) (out string) {
	return "JSON => {\n" + printPersedJsonNode(n, "    ") + "}\n"
}

func printPersedJsonNode(n map[string]interface{}, indent string) (out string) {
	out = ""
	for n, v := range n {
		out += indent + n + " => "
		if str, isStr := v.(string); isStr {
			out += str + "\n"
			continue
		}
		if flt, isFlt := v.(float64); isFlt {
			out += fmt.Sprintf("%f\n", flt) + "\n"
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
			out += printPersedJsonNode(m, indent+"    ")
			out += indent + "}\n"
			continue
		}
		out += "?\n"
	}
	return
}

func getStringByPath(root map[string]interface{}, path string) (string, string) {
	val, typ := getValueByPath(root, path)
	if str, isStr := val.(string); typ == "string" && isStr {
		return str, typ
	}
	return "", "none"
}

func getFloat64ByPath(root map[string]interface{}, path string) (float64, string) {
	val, typ := getValueByPath(root, path)
	if f, isFlt := val.(float64); typ == "float64" && isFlt {
		return f, typ
	}
	return 0, "none"
}

func getBoolByPath(root map[string]interface{}, path string) (bool, string) {
	val, typ := getValueByPath(root, path)
	if b, isBool := val.(bool); typ == "bool" && isBool {
		return b, typ
	}
	return false, "none"
}

func getMapByPath(root map[string]interface{}, path string) (map[string]interface{}, string) {
	val, typ := getValueByPath(root, path)
	if m, isMap := val.(map[string]interface{}); typ == "map" && isMap {
		return m, typ
	}
	return nil, "none"
}

func getValueByPath(root map[string]interface{}, path string) (interface{}, string) {
	parts := strings.Split(path, "/")
	n := root
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
