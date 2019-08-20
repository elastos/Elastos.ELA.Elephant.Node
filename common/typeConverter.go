package common

import (
	"reflect"
	"strconv"
	"strings"
)

//Map2Struct convert map into struct
func Map2Struct(src map[string]interface{}, destStrct interface{}) {
	value := reflect.ValueOf(destStrct)
	e := value.Elem()
	for k, v := range src {
		f := e.FieldByName(strings.ToUpper(k[:1]) + k[1:])
		if !f.IsValid() {
			continue
		}
		if !f.CanSet() {
			continue
		}
		mv := reflect.ValueOf(v)
		mvt := mv.Kind().String()
		sft := f.Kind().String()
		if sft != mvt {
			if mvt == "string" && (strings.Index(sft, "int") != -1) {
				if sft == "int64" {
					i, err := strconv.ParseInt(v.(string), 10, 64)
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "int32" {
					i, err := strconv.ParseInt(v.(string), 10, 32)
					r := int32(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				} else if sft == "int" {
					i, err := strconv.Atoi(v.(string))
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "uint64" {
					i, err := strconv.ParseUint(v.(string), 10, 64)
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "uint32" {
					i, err := strconv.ParseUint(v.(string), 10, 32)
					r := uint32(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				} else if sft == "uint" {
					i, err := strconv.ParseUint(v.(string), 10, 0)
					r := uint(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				}
			}

			if mvt == "int" && (strings.Index(sft, "int") != -1) {
				if sft == "int64" {
					r := int64(v.(int))
					f.Set(reflect.ValueOf(r))
				} else if sft == "int32" {
					r := int32(v.(int))
					f.Set(reflect.ValueOf(r))
				} else if sft == "uint64" {
					r := uint64(v.(int))
					f.Set(reflect.ValueOf(r))
				} else if sft == "uint32" {
					r := uint32(v.(int))
					f.Set(reflect.ValueOf(r))
				} else if sft == "uint" {
					r := uint(v.(int))
					f.Set(reflect.ValueOf(r))
				}
			}

			// make string and string[] more friendly
			if mvt == "string" && sft == "slice" {
				_, ok := f.Interface().([]string)
				if ok {
					f.Set(reflect.ValueOf(strings.Split(v.(string), ",")))
				}
			}

			// make float64 and int64 more friendly
			if mvt == "float64" && sft == "int64" {
				f.Set(reflect.ValueOf(int64(v.(float64))))
			}
			continue
		}
		f.Set(mv)
	}
}
