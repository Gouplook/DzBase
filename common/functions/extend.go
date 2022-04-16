package functions

import "reflect"

//反射Model、初始化字段名称
func ReflectModel(structPtr interface{}){
	rType := reflect.TypeOf(structPtr)
	rVal := reflect.ValueOf(structPtr)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("structPtr must be pointer struct.")
	}
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		key := t.Tag.Get("default")
		if key == ""{
			f.Set(reflect.ValueOf(""))
		}else{
			f.Set(reflect.ValueOf(key))
		}
	}
}