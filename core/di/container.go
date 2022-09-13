package di

import (
	"reflect"
)

//定义一个map类型
type Mapper map[reflect.Type]reflect.Value

//容器
type Container struct {
	//存储映射关系
	Mapper Mapper
}

//add
func (this Container) Add(name interface{}) {
	t := reflect.TypeOf(name)
	//判断是否指针内类型
	if t.Kind() != reflect.Ptr {
		panic("require ptr object")
	}
	this.Mapper[t] = reflect.ValueOf(name)
}

//get
func (this Container) Get(module interface{}) reflect.Value {
	var t reflect.Type
	if bt, ok := module.(reflect.Type); ok {
		t = bt
	} else {
		reflect.TypeOf(module)
	}

	if v, ok := this.Mapper[t]; ok {
		return v
	}

	//处理接口继承
	for k, v := range this.Mapper {
		if k.Implements(t) {
			return v
		}
	}

	return reflect.Value{}
}
