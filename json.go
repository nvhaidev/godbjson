package godbjson

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"reflect"
)

var _filename = ""

type DataArray struct {
	Value []interface{}
}

func setFilename(filename string) {
	_filename = filename
}
func uid() string {
	id := uuid.New()
	return id.String()
}
func readFile(filename string) []interface{} {
	jsonFile, _ := os.Open(filename)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data DataArray
	err := json.Unmarshal(byteValue, &data.Value)
	if err != nil {
		return nil
	}
	return data.Value
}

func write(filename string, data []interface{}) {
	jsonData, _ := json.Marshal(data)
	err := ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return
	}
}
func getData() []interface{} {
	data := readFile(_filename)
	if data == nil {
		return nil
	}

	return data
}
func setData(data interface{}) {
	content := getData()
	if content == nil {
		content = []interface{}{}
	}
	content = append(content, data)
	write(_filename, content)
}
func create(data map[string]interface{}) interface{} {
	data["_id"] = uid()
	setData(data)
	return data
}
func findById(id string) interface{} {
	data := getData()
	if data == nil {
		return nil
	}
	for _, v := range data {
		if v.(map[string]interface{})["_id"] == id {
			return v
		}
	}
	return nil
}
func findOne(query map[string]interface{}) interface{} {
	data := getData()
	if data == nil {
		return nil
	}
	var keyWhere string
	var valueWhere any
	var valueExclude []string
	var result interface{}
	for k, v := range query {
		if k == "where" {
			if rec, ok := v.(map[string]interface{}); ok {
				for key, val := range rec {
					keyWhere = key
					valueWhere = val
					break
				}
			} else {
				fmt.Printf("where not a map[string]interface{}: %v\n", v)
			}
		}
		if k == "exclude" {
			if rec, ok := v.([]string); ok {
				valueExclude = rec
			} else {
				fmt.Printf("exclude not a []string: %v\n", v)
			}
		}
	}
	for _, v := range data {
		if v.(map[string]interface{})[keyWhere] == valueWhere {
			result = v
			break
		}
	}
	if result == nil {
		return nil
	}
	if valueExclude != nil {
		for _, v := range valueExclude {
			delete(result.(map[string]interface{}), v)
		}
	}

	return result
}
func findAll(query map[string]interface{}) []interface{} {
	data := getData()
	if data == nil {
		return nil
	}
	var keyWhere string = ""
	var valueWhere any = nil
	var valueExclude []string = nil
	var valueLimit int = -1
	var valueOffset int = -1
	var result []interface{}

	for k, v := range query {
		if k == "where" {
			if rec, ok := v.(map[string]interface{}); ok {
				for key, val := range rec {
					keyWhere = key
					valueWhere = val
					break
				}
			} else {
				fmt.Printf("where not a map[string]interface{}: %v\n", v)
			}
		}
		if k == "exclude" {
			if rec, ok := v.([]string); ok {
				valueExclude = rec
			} else {
				fmt.Printf("exclude not a []string: %v\n", v)
			}
		}
		if k == "limit" {
			if rec, ok := v.(int); ok {
				valueLimit = rec
			} else {
				fmt.Printf("limit not a int: %v\n", v)
			}
		}
		if k == "offset" {
			if rec, ok := v.(int); ok {
				valueOffset = rec
			} else {
				fmt.Printf("offset not a int: %v\n", v)
			}
		}
	}
	if valueOffset != -1 {
		data = data[valueOffset:]
	}
	if valueWhere != nil && keyWhere != "" {
		for _, v := range data {
			if v.(map[string]interface{})[keyWhere] == valueWhere {
				result = append(result, v)
			}
		}
	} else {
		result = data
	}

	if valueLimit != -1 {
		result = result[:valueLimit]
	}
	if valueExclude != nil {
		result = exclude(result, valueExclude)
	}
	return result
}
func update(id string, data map[string]interface{}) interface{} {
	content := getData()
	if content == nil {
		return nil
	}
	for i, v := range content {
		if v.(map[string]interface{})["_id"] == id {
			for k, vv := range data {
				content[i].(map[string]interface{})[k] = vv
			}
			write(_filename, content)
			return content[i]
		}
	}
	return nil
}
func remove(s []interface{}, i int) []interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
func exclude(s []interface{}, remote []string) []interface{} {
	if s == nil {
		return nil
	}
	if remote == nil {
		return s
	}
	for _, v := range remote {
		for i, vv := range s {
			delete(vv.(map[string]interface{}), v)
			s[i] = vv
		}
	}
	return s
}
func deleteById(id string) interface{} {
	content := getData()
	if content == nil {
		return nil
	}
	for i, v := range content {
		if v.(map[string]interface{})["_id"] == id {
			value := remove(content, i)
			write(_filename, value)
			return v
		}
	}
	return nil
}
func findAllAndCount(query map[string]interface{}) interface{} {
	content := findAll(query)
	data := getData()
	if data == nil {
		return nil
	}
	if content == nil {
		return nil
	}
	return map[string]interface{}{
		"data":  content,
		"count": len(content),
	}
}
func findIndex(query map[string]interface{}) int {
	content := getData()
	var keyWhere string
	var valueWhere any
	if content == nil {
		return -1
	}
	for k, v := range query {
		if k == "where" {
			if rec, ok := v.(map[string]interface{}); ok {
				for key, val := range rec {
					keyWhere = key
					valueWhere = val
					break
				}
			} else {
				fmt.Printf("where not a map[string]interface{}: %v\n", v)
			}
		}
	}
	for i, v := range content {
		if reflect.TypeOf(v.(map[string]interface{})[keyWhere]).String() == "float64" {
			float64ToInt := int(v.(map[string]interface{})[keyWhere].(float64))
			if float64ToInt == valueWhere {
				return i
			}
		}
		if v.(map[string]interface{})[keyWhere] == valueWhere {
			return i
		}
	}
	return -1
}
func getIndex(data []interface{}, index int) interface{} {
	if data == nil {
		return nil
	}
	if index < 0 || index >= len(data) {
		return nil
	}
	return data[index]
}
func includes(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func filter(s []interface{}, f func(interface{}) bool) []interface{} {
	var r []interface{}
	for _, e := range s {
		if f(e) {
			r = append(r, e)
		}
	}
	return r
}