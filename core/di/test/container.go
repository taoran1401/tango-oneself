package test

import (
	"reflect"
)

type BeanMapper map[reflect.Type]reflect.Value

/*type Container struct {
	Mapper map[reflect.Type]reflect.Value
}*/

//插入值
func (this BeanMapper) Add(bean interface{}) {
	t := reflect.TypeOf(bean)
	//判断是否指针内类型
	if t.Kind() != reflect.Ptr {
		panic("require ptr object")
	}
	this[t] = reflect.ValueOf(bean)
}

//获取
func (this BeanMapper) Get(bean interface{}) reflect.Value {
	var t reflect.Type
	if bt, ok := bean.(reflect.Type); ok {
		t = bt
	} else {
		reflect.TypeOf(bean)
	}
	if v, ok := this[t]; ok {
		return v
	}

	//处理接口继承
	for k, v := range this {
		if k.Implements(t) {
			return v
		}
	}

	return reflect.Value{}
}
