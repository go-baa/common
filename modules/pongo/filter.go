package pongo

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/go-baa/baa"
	"github.com/go-baa/common/util"
	"github.com/micate/pongo2"
)

// Filters ...
func Filters(b *baa.Baa) map[string]pongo2.FilterFunction {
	return map[string]pongo2.FilterFunction{
		"duration":  duration,
		"json":      jsonEncode,
		"html":      html,
		"trim":      trim,
		"inline":    inline,
		"index":     index,
		"string":    parseValueToString,
		"striptags": striptags,
	}
}

func parseValueToString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.String()), nil
}

func duration(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	second := in.Integer()
	return pongo2.AsValue(util.SecondToTime(second)), nil
}

func jsonEncode(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	enc, err := json.Marshal(in.Interface())
	if err != nil {
		return nil, nil
	}
	return pongo2.AsSafeValue(string(enc)), nil
}

func html(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsSafeValue(in.String()), nil
}

func trim(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	str := in.String()
	str = strings.TrimSpace(str)
	return pongo2.AsValue(str), nil
}

func inline(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	str := in.String()
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.TrimSpace(str)
	return pongo2.AsValue(str), nil
}

func index(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	found := new(pongo2.Value)
	in.Iterate(func(idx, count int, key, value *pongo2.Value) bool {
		if key.Interface() == param.Interface() {
			found = value
			return false
		}
		return true
	}, func() {})
	return found, nil
}

func striptags(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	str := in.String()
	str = bytes.NewBuffer(util.StripTags(bytes.NewBufferString(str).Bytes())).String()
	return pongo2.AsValue(str), nil
}
