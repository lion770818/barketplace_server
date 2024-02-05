package utils

import (
	"errors"
	"reflect"
)

// =========================================================================================================================================================================
// 一种可应用于任意slice的remove方法实现 http://weilin.me/articles/geneicRemoveForSlice.html
type sliceHelper struct {
	slicePtr interface{}
}

// =========================================================================================================================================================================
// 只能用于移除slice的数组中的元素
func (t *sliceHelper) Remove(index int) error {
	if t.slicePtr == nil {
		return errors.New("slice ptr is nil!")
	}

	slicePtrValue := reflect.ValueOf(t.slicePtr)
	// 必须为指针
	if slicePtrValue.Type().Kind() != reflect.Ptr {
		return errors.New("should be slice ptr!")
	}

	sliceValue := slicePtrValue.Elem()
	// 必须为slice
	if sliceValue.Type().Kind() != reflect.Slice {
		return errors.New("should be slice ptr!")
	}

	if index < 0 || index >= sliceValue.Len() {
		return errors.New("index out of range!")
	}
	sliceValue.Set(reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, sliceValue.Len())))
	return nil
}

// =========================================================================================================================================================================
// 参数slicePtr必须是指向slice的指针
func SliceHelper(slicePtr interface{}) *sliceHelper {
	return &sliceHelper{slicePtr}
}
