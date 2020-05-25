package functional

import (
	"reflect"
	"sync"
	"testing"
)

func BenchmarkMap(b *testing.B) {
	res := []interface{}{}
	arr := []interface{}{"a", "ab", "abc", "abcd"}
	var wg sync.WaitGroup

	for index, item := range arr {
		wg.Add(1)

		go func(index *int, item *interface{}, res *[]interface{}) {
			defer wg.Done()
			*res = append(*res, func(item *interface{}) interface{} {
				return reflect.TypeOf(item).Name
			})
		}(&index, &item, &res)
	}

	wg.Wait()
}

func BenchmarkSync(b *testing.B) {
	res := []interface{}{}
	arr := []interface{}{"a", "ab", "abc", "abcd"}
	for _, item := range arr {
		res = append(res, reflect.TypeOf(item).Name)
	}
}
