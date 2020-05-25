package functional

import (
	"reflect"
	"sync"
)

func Map(it func(*interface{}, *int) interface{}, arr interface{}) []interface{} {
	res := []interface{}{}
	var wg sync.WaitGroup

	finalArr := InterfaceSlice(arr)

	for index, item := range finalArr {
		wg.Add(1)

		go func(index *int, item *interface{}, res *[]interface{}) {
			defer wg.Done()
			*res = append(*res, it(item, index))
		}(&index, &item, &res)
	}

	wg.Wait()

	return res
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
