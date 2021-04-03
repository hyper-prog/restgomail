package main

/*  Smart JSON functions - Helper functions to query JSON
    (C) 2021 Péter Deák (hyper80@gmail.com)
    License: GPLv2
*/

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SmartJSON struct {
	parsed interface{}
}

func parseSmartJSON(rawdata []byte) (SmartJSON, error) {
	s := SmartJSON{}
	err := json.Unmarshal(rawdata, &s.parsed)
	return s, err
}

func (sJSON SmartJSON) toFormattedString() (out string) {
	return jsonNodeToString(sJSON.parsed, "") + "\n"
}

func (sJSON SmartJSON) toPrettify() (out string) {
	return jsonNodeToString(sJSON.parsed, "") + "\n"
}

func jsonNodeToString(v interface{}, indent string) (out string) {
	out = ""
	if m, isMap := v.(map[string]interface{}); isMap {
		out += "{\n" + indent + "  "
		c := 0
		for n, v := range m {
			sep := ""
			if c > 0 {
				sep = ",\n  " + indent
			}
			out += sep + "\"" + n + "\":" + jsonNodeToString(v, indent + "  ")
			c++
		}
		out += "\n" + indent + "}"
		return out
	}
	if arr, isArray := v.([]interface{}); isArray {
		out += "[\n" + indent + "  "
		l := len(arr)
		for i := 0; i < l; i++ {
			sep := ""
			if i > 0 {
				sep = ",\n  " + indent
			}
			out += sep + jsonNodeToString(arr[i], indent + "  ")
		}
		out += "\n" + indent + "]"
		return out
	}
	if str, isStr := v.(string); isStr {
		out += "\"" + str + "\""
		return out
	}
	if flt, isFlt := v.(float64); isFlt {
		out += fmt.Sprintf("%f", flt)
		return out
	}
	if b, isB := v.(bool); isB {
		if b {
			out += "true"
		} else {
			out += "false"
		}
		return out
	}
	if v == nil {
		out += "null"
		return out
	}
	return "" //should not happend
}

func pathEvalNode(last interface{}) (interface{}, string) {
	if str, isStr := last.(string); isStr {
		return str, "string"
	}
	if flt, isFlt := last.(float64); isFlt {
		return flt, "float64"
	}
	if bo, isBool := last.(bool); isBool {
		return bo, "bool"
	}
	if mp, isMap := last.(map[string]interface{}); isMap {
		return mp, "map"
	}
	if ar, isArr := last.([]interface{}); isArr {
		return ar, "array"
	}
	if last == nil {
		return nil, "null"
	}
	return nil, "none"
}

func (sJSON SmartJSON) getValueByPath(path string) (interface{}, string) {
	parts := strings.Split(path, "/")
	n := sJSON.parsed
	for i := 0; i < len(parts); i++ {
		if map_node, isMap_node := n.(map[string]interface{}); isMap_node {
			map_node_value, ok := map_node[parts[i]]
			if !ok {
				return nil, "none"
			}
			if i == len(parts)-1 {
				return pathEvalNode(map_node_value)
			}
			n = map_node_value
			continue
		}
		if arr_node, isArr_node := n.([]interface{}); isArr_node {
			if len(arr_node) == 0 {
				return nil, "none"
			}
			var arr_node_item interface{}
			if parts[i] == "[]" {
				arr_node_item = arr_node[0]
			} else {
				r, rxerr := regexp.Compile(`^\[([0-9]+)\]$`)
				if rxerr != nil {
					return nil, "none"
				}
				matches := r.FindStringSubmatch(parts[i])
				if len(matches) != 2 {
					return nil, "none"
				}
				index, erratoi := strconv.Atoi(matches[1])
				if erratoi != nil {
					return nil, "none"
				}
				if index >= len(arr_node) {
					return nil, "none"
				}
				arr_node_item = arr_node[index]
			}

			if i == len(parts)-1 {
				return pathEvalNode(arr_node_item)
			}
			n = arr_node_item
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
