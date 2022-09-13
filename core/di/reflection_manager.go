package di

import (
	"reflect"
)

type ReflectionManagerImpl struct {
	Container Container
	ExprMap   map[string]interface{}
}

var ReflectionManager *ReflectionManagerImpl

//定义时自动创建初始化ReflectionManager
func init() {
	ReflectionManager = NewReflectionManager()
}

//构造方法
func NewReflectionManager() *ReflectionManagerImpl {
	return &ReflectionManagerImpl{}
}

//设置
func (this *ReflectionManagerImpl) Set(vlist ...interface{}) {
	if vlist == nil || len(vlist) == 0 {
		return
	}
	for _, v := range vlist {
		this.Container.Add(v)
	}
}

//获取
func (this *ReflectionManagerImpl) Get(module interface{}) interface{} {
	if module == nil {
		return nil
	}
	getV := this.Container.Get(module)
	if getV.IsValid() {
		return getV.Interface()
	}
	return nil
}

func (this *ReflectionManagerImpl) Apply(module interface{}) {
	if module == nil {
		return
	}
	//返回一个 reflect.Value 对象
	v := reflect.ValueOf(module)

	//判断是否指针
	if v.Kind() == reflect.Ptr {
		//Elem返回接口v所包含的值
		//或v指向的指针。
		//如果v的种类不是接口或Ptr，它会panic。
		//如果v为空，则返回零值。
		v = v.Elem()
	}

	//判断是否结构体
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		//返回结构类型的第i个字段。
		//如果类型的Kind不是Struct，它会panic。
		field := v.Type().Field(i)
		//判断是否可以修改并且tag指定的值不为空
		if v.Field(i).CanSet() && field.Tag.Get("inject") != "" {
			//fmt.Printf("%T %v \n", field.Tag.Get("inject"), field.Tag.Get("inject"))

			//判断tag是否存在
			if field.Tag.Get("inject") != "-" {
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

//构建
func (this *ReflectionManagerImpl) Build(cfgs ...interface{}) {
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

//提取结构体标签
/*func (this *ReflectionManagerImpl) ExtractStructTag(bean interface{}) {
	//通过反射获取type定义
	s := reflect.TypeOf(bean).Elem()
	for i := 0; i < s.NumField(); i++ {
		//获取tag存入
		tag := s.Field(i).Tag.Get("inject")
		this.Container.Add(tag)
	}
}*/
