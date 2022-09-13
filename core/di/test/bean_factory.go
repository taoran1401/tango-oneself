package test

import (
	"reflect"
)

type BeanFactoryImpl struct {
	beanMapper BeanMapper
	ExprMap    map[string]interface{}
}

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactory()
}

func NewBeanFactory() *BeanFactoryImpl {
	return &BeanFactoryImpl{
		beanMapper: make(BeanMapper),
		ExprMap:    make(map[string]interface{}),
	}
}

//设置
func (this *BeanFactoryImpl) Set(vlist ...interface{}) {
	if vlist == nil || len(vlist) == 0 {
		return
	}

	for _, v := range vlist {
		this.beanMapper.Add(v)
	}
}

//获取
func (this *BeanFactoryImpl) Get(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	get_v := this.beanMapper.Get(v)
	if get_v.IsValid() {
		return get_v.Interface()
	}
	return nil
}

func (this *BeanFactoryImpl) Apply(bean interface{}) {
	if bean == nil {
		return
	}

	v := reflect.ValueOf(bean)
	//判断是否指针
	if v.Kind() == reflect.Ptr {
		// Elem returns the value that the interface v contains
		// or that the pointer v points to.
		// It panics if v's Kind is not Interface or Ptr.
		// It returns the zero Value if v is nil.
		v = v.Elem()
	}
	//判断是否结构体
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		// Field returns a struct type's i'th field.
		// It panics if the type's Kind is not Struct.
		// It panics if i is not in the range [0, NumField()).
		field := v.Type().Field(i)
		//判断是否可以修改并且tag指定的值不为空
		if v.Field(i).CanSet() && field.Tag.Get("inject") != "" {
			//判断是否存在
			if field.Tag.Get("inject") != "-" {
				//多例支持
				//表达式方式支持
				//ret := expr.BeanExpr(field.Tag.Get("inject"), this.ExprMap)
				/*for k, v := range ret {
					fmt.Printf("%k - %v \n", k, v)
				}*/
				/*if ret != nil && !ret.IsEmpty() {
					retValue := ret[0]
					if retValue != nil {
						v.Field(i).Set(reflect.ValueOf(ret[0]))
						this.Apply(retValue)
					}
				}*/
			} else {
				//单例
				getV := this.Get(field.Type)
				if getV != nil {
					v.Field(i).Set(reflect.ValueOf(getV))
					this.Apply(getV)
				}
			}
		}
	}
}

func (this *BeanFactoryImpl) Config(cfgs ...interface{}) {
	for _, cfg := range cfgs {
		t := reflect.TypeOf(cfg)
		//判断是否指针
		if t.Kind() != reflect.Ptr {
			panic("required ptr object")
		}

		//把自己加入到bean
		this.Set(cfg)
		t = t.Elem()
		this.ExprMap[t.Name()] = cfg //自动构建
		v := reflect.ValueOf(cfg).Elem()
		for i := 0; i < t.NumMethod(); i++ {
			method := v.Method(i)
			callRet := method.Call(nil)
			if callRet != nil && len(callRet) == 1 {
				this.Set(callRet[0].Interface())
			}
		}

	}
}
