package convertutil

import (
	"errors"
	"reflect"
	"strings"
)

func Convert2Map(obj interface{}) (placeholderMap map[string]interface{}) {
	placeholderMap = make(map[string]interface{}, 16)
	recordType := reflect.TypeOf(obj)
	recordValue := reflect.ValueOf(obj)
	for i := 0; i < recordType.NumField(); i++ {
		field := recordType.Field(i)
		tag := field.Tag.Get("ep")
		if tag == "" {
			continue
		}
		if recordValue.FieldByName(field.Name).IsZero() {
			placeholderMap[tag] = ""
		} else {
			value := recordValue.FieldByName(field.Name).Interface()
			placeholderMap[tag] = value
		}
	}
	return
}

func GetMongoModelFieldValue(obj interface{}) (map[string]interface{}, error) {
	fieldValueMap := make(map[string]interface{})
	err := getAllFieldValue(obj, "", fieldValueMap)
	if err != nil {
		return nil, err
	}
	return fieldValueMap, nil
}

func getAllFieldValue(obj interface{}, prefix string, m map[string]interface{}) error {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("bson")
		var key string
		if tagVals := strings.Split(tag, ","); len(tagVals) > 0 {
			key = tagVals[0]
			if prefix != "" {
				key = prefix + "." + key
			}
		}
		if key == "" {
			return errors.New("parse field error")
		}
		fv := v.FieldByName(f.Name)
		if fv.Kind() == reflect.Struct || fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			m[key] = fv.Interface()
			fv = reflect.Indirect(fv)
			err := getAllFieldValue(fv.Addr().Interface(), key, m)
			if err != nil {
				return err
			}
			continue
		}
		if fv.Kind() == reflect.Map {
			m[key] = fv.Interface()
			for _, k := range fv.MapKeys() {
				kk := key + "." + k.String()
				m[kk] = fv.MapIndex(k).Interface()
			}
			continue
		}
		m[key] = fv.Interface()
	}
	return nil
}

type TestField struct {
	Name          string            `bson:"name"`
	Age           int64             `bson:"age"`
	Address       AddressInfo       `bson:"address"`
	ProductIds    []string          `bson:"product_ids"`
	SpuName       map[string]string `bson:"spu_name"`
	ExtraNeedings []AddressInfo     `bson:"extra_needings"`
}

type AddressInfo struct {
	Detail    string         `bson:"detail"`
	PostCode  string         `bson:"post_code"`
	Emergency ContractInfo   `bson:"emergency"`
	Contract  []ContractInfo `bson:"contract"`
}

type ContractInfo struct {
	Name  string `bson:"name"`
	Phone string `bson:"phone"`
}
