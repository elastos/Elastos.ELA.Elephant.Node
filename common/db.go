package common

import (
	"container/list"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"strconv"
)

type Dba struct {
	*sql.DB
}

func NewInstance(filePath string) (*Dba, error) {
	db, err := sql.Open("sqlite3", filePath+"/dpos/dpos.db")
	if err != nil {
		return nil, err
	}
	return &Dba{db}, nil
}

func (d *Dba) Qu(s string) (*list.List, error) {
	rows, err := d.Query(s)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([][]byte, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	retList := list.New()

	// Fetch rows
	for rows.Next() {
		retMap := make(map[string]interface{})
		retList.PushBack(retMap)
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				// skip nil value
				continue
			}
			retMap[columns[i]] = string(col)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retList, nil
}

//
func (d *Dba) ToStruct(sql string, strct interface{}) ([]interface{}, error) {
	l, err := d.Qu(sql)
	if err != nil {
		return nil, err
	}
	i := 0
	r := make([]interface{}, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		v := reflect.New(reflect.TypeOf(strct))
		src := e.Value.(map[string]interface{})
		vi := v.Interface()
		Map2Struct(src, vi)
		r[i] = reflect.ValueOf(vi).Interface()
		i++
	}
	return r, nil
}

func (d *Dba) ToInt(sql string) (int, error) {

	l, err := d.Qu(sql)
	if err != nil || l.Len() == 0 {
		return -1, err
	}
	m := l.Front().Value.(map[string]interface{})
	for _, v := range m {
		return strconv.Atoi(v.(string))
	}

	return -1, err
}

func (d *Dba) ToFloat(sql string) (float64, error) {

	l, err := d.Qu(sql)
	if err != nil || l.Len() == 0 {
		return -1, err
	}
	m := l.Front().Value.(map[string]interface{})
	for _, v := range m {
		return strconv.ParseFloat(v.(string), 64)
	}

	return -1, err
}

//
func (d *Dba) ToString(sql string) (string, error) {

	l, err := d.Qu(sql)
	if err != nil || l.Len() == 0 {
		return "", err
	}
	m := l.Front().Value.(map[string]interface{})
	for _, v := range m {
		return v.(string), nil
	}

	return "", err
}
