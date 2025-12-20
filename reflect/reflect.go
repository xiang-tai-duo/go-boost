// Package reflect
// File:        reflect.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/reflect/reflect.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: C#-style reflection functionality for Go applications
// --------------------------------------------------------------------------------
package reflect

import (
	_reflect "reflect"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	PROPERTY_INFO struct {
		Name  string
		Type  string
		Value interface{}
	}

	REFLECT struct {
		object     _reflect.Value
		objectType _reflect.Type
	}
)

//goland:noinspection GoUnusedExportedFunction
func New(obj interface{}) *REFLECT {
	result := &REFLECT{
		object:     _reflect.ValueOf(obj),
		objectType: _reflect.TypeOf(obj),
	}
	return result
}

func (r *REFLECT) GetProperties() []PROPERTY_INFO {
	result := []PROPERTY_INFO{}

	if r.object.Kind() == _reflect.Ptr {
		r.object = r.object.Elem()
		r.objectType = r.objectType.Elem()
	}

	for i := 0; i < r.objectType.NumField(); i++ {
		field := r.objectType.Field(i)
		fieldValue := r.object.Field(i)

		var value interface{}
		if fieldValue.CanInterface() {
			value = fieldValue.Interface()
		}

		result = append(result, PROPERTY_INFO{
			Name:  field.Name,
			Type:  field.Type.String(),
			Value: value,
		})
	}

	return result
}

func (r *REFLECT) GetProperty(name string) (PROPERTY_INFO, bool) {
	var result PROPERTY_INFO
	var found bool

	if r.object.Kind() == _reflect.Ptr {
		r.object = r.object.Elem()
		r.objectType = r.objectType.Elem()
	}

	field, found := r.objectType.FieldByName(name)
	if !found {
		return result, found
	}

	fieldValue := r.object.FieldByName(name)
	var value interface{}
	if fieldValue.CanInterface() {
		value = fieldValue.Interface()
	}

	result = PROPERTY_INFO{
		Name:  field.Name,
		Type:  field.Type.String(),
		Value: value,
	}
	return result, found
}

func (r *REFLECT) SetProperty(name string, value interface{}) bool {
	result := false

	if r.object.Kind() == _reflect.Ptr {
		r.object = r.object.Elem()
		r.objectType = r.objectType.Elem()
	}

	fieldValue := r.object.FieldByName(name)
	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return result
	}

	val := _reflect.ValueOf(value)
	if val.Type() != fieldValue.Type() {
		return result
	}

	fieldValue.Set(val)
	result = true
	return result
}

func (r *REFLECT) GetType() string {
	result := r.objectType.String()
	return result
}

func (r *REFLECT) GetTypeName() string {
	result := r.objectType.Name()
	return result
}
