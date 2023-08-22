package reflect

import (
	"go/ast"
	"reflect"
	"strconv"
)

func Any(arg any) string {
	return PrintAny(reflect.ValueOf(arg))
}
func PrintAny(value reflect.Value) string {

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'E', -1, 64)
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Chan, reflect.Func:
		return "Type: " + value.Type().String() + "\n" + "Address: 0x" + strconv.FormatUint(uint64(value.Pointer()), 16)
	case reflect.Invalid:
		return "valid type"
	default:
		return value.Type().String()
	}
}

func Empty() {
	//var i interface{} //接口没有指向具体的值
	//v := reflect.ValueOf(i)
	//fmt.Printf("v持有值 %t, type of v is Invalid %t\n", v.IsValid(), v.Kind() == reflect.Invalid)
	//
	//var user *User = nil
	//v = reflect.ValueOf(user) //Value指向一个nil
	//if v.IsValid() {
	//	fmt.Printf("v持有的值是nil %t\n", v.IsNil()) //调用IsNil()前先确保IsValid()，否则会panic
	//}
	//
	//var u User //只声明，里面的值都是0值
	//v = reflect.ValueOf(u)
	//if v.IsValid() {
	//	fmt.Printf("v持有的值是对应类型的0值 %t\n", v.IsZero()) //调用IsZero()前先确保IsValid()，否则会panic
	//}
	println(ast.IsExported("Name"))
}
