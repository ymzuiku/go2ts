package go2ts

var typeMap = map[string]string{
	"string":                    "string",
	"bool":                      "boolean",
	"int":                       "number",
	"int8":                      "number",
	"int16":                     "number",
	"int32":                     "number",
	"int64":                     "number",
	"uint":                      "number",
	"uint8":                     "number",
	"uint16":                    "number",
	"uint32":                    "number",
	"uint64":                    "number",
	"float32":                   "number",
	"float64":                   "number",
	"error":                     "string",
	"[]any":                     "any[]",
	"[]interface {}":            "any[]",
	"[]string":                  "string[]",
	"[]int":                     "number[]",
	"[]float32":                 "number[]",
	"[]map[string]any":          "Record<string, any>[]",
	"[]map[string]interface {}": "Record<string, any>[]",
	"map[string]interface {}":   "Record<string, any>",
	"map[string]string":         "Record<string, string>",
	"map[string]bool":           "Record<string, boolean>",
	"time.Time":                 "string",
}

var tsBaseType = getTsBaseType()

func getTsBaseType() map[string]bool {
	out := map[string]bool{}
	for _, v := range typeMap {
		out[v] = true
	}
	return out
}
