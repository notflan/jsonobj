package jsonobj

import (
	"github.com/pkg/errors"
	"encoding/json"
	"reflect"
)

type JSONType int32


const (
	JSONUndef JSONType = iota
	JSONString
	JSONInt
	JSONMap
	JSONBool
	JSONFloat
	JSONArray
)

type ObjMap map[string]*JSONObj

type JSONObj struct {
	jtype JSONType

	tInt int64
	tString string
	tBool bool
	tFloat float64
	tArray []*JSONObj
	tMap ObjMap
}

func (m ObjMap) In(key string) *JSONObj {
	return In(m, key)
}

func (obj JSONObj)  MarshalJSON() ([]byte, error) {
	switch obj.jtype {
		case JSONUndef:
			return json.Marshal(nil)
		case JSONInt:
			return json.Marshal(obj.tInt)
		case JSONString:
			return json.Marshal(obj.tString)
		case JSONBool:
			return json.Marshal(obj.tBool)
		case JSONFloat:
			return json.Marshal(obj.tFloat)
		case JSONMap:
			rmap := make(map[string]JSONObj)
			for k, v := range obj.tMap {
				rmap[k] = *v
			}
			return json.Marshal(rmap)
		case JSONArray:
			rar := make([]JSONObj, len(obj.tArray))
			for i, v := range obj.tArray {
				rar[i] = *v
			}
			return json.Marshal(rar)
	}
	return nil, errors.New("json marshal: unknown type")
}

func typeAttempt(obj *JSONObj, v interface{}) {
	//fmt.Printf("umtype: %T\n", v)
	switch v.(type) {
		case uint:
			typeAttempt(obj, int(v.(uint)))
		case int64:
			obj.SetInt(v.(int64))
		case int:
			obj.SetInt(int64(v.(int)))
		case string:
			obj.SetString(v.(string))
		case float32:
			typeAttempt(obj, float64(v.(float32)))
		case float64:
			obj.SetFloat(v.(float64))
		case bool:
			obj.SetBool(v.(bool))
		case []interface{}:
			from := v.([]interface{})
			ar:=make([]*JSONObj, len(from))
			
			for i, val := range from {
				obj2 := NewObj()
				typeAttempt(obj2, val)
				ar[i] = obj2
			}
			obj.SetArray(ar)
		case map[string]interface{}:
			from := v.(map[string]interface{})
			mp := make(ObjMap)

			for key, val:= range from {
				obj2 := NewObj()
				typeAttempt(obj2, val)
				mp[key] = obj2
			}
			obj.SetMap(mp)
		case *JSONObj:
			*obj = *(v.(*JSONObj))
		default:
			if v!=nil && reflect.TypeOf(v).Kind() == reflect.Slice {
				val := reflect.ValueOf(v)
				rt := make([]interface{}, val.Len())
				for i:=0;i<val.Len();i++ {
					rt[i] = val.Index(i).Interface()
				}
				typeAttempt(obj, rt)
			} else {
				obj.SetUndef()
			}
	}
}

func (obj *JSONObj) UnmarshalJSON(b []byte) error {
	var js interface{}
	if err := json.Unmarshal(b, &js); err!=nil {
		return errors.Wrap(err, "json unmarshal failed")
	}

	typeAttempt(obj, js)
	return nil
}

type tuple struct {
	key string
	value interface{}
}

func KV(key string, value interface{}, rest ...interface{}) tuple {
	if rest!=nil && len(rest)>0 {
		return tuple{key:key, value:append([]interface{}{value}, rest...)}
	} else {
		return tuple{key:key, value:value}
	}
}

func Table(init ...tuple) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, val := range init {
		ret[val.key] = val.value
	}
	return ret
}

func NewObj(init ...interface{}) *JSONObj {
	obj := new(JSONObj)
	obj.jtype= JSONUndef
	if init == nil || len(init)<1 {
		return obj
	} else if len(init) == 1 {
		typeAttempt(obj, init[0])
	} else {
		typeAttempt(obj, init)
	}
	return obj
}

func N(init ...interface{}) *JSONObj {
	return NewObj(init...)
}

func (obj *JSONObj) SetInt(in int64) *JSONObj {
	obj.jtype = JSONInt
	obj.tInt = in

	return obj
}
func (obj *JSONObj) SetString(str string) *JSONObj {
	obj.jtype = JSONString
	obj.tString = str

	return obj
}
func (obj *JSONObj) SetFloat(float float64) *JSONObj {
	obj.jtype = JSONFloat
	obj.tFloat = float

	return obj
}

func (obj *JSONObj) SetBool(bl bool) *JSONObj {
	obj.jtype = JSONBool
	obj.tBool = bl

	return obj
}

func (obj *JSONObj) SetArray(from []*JSONObj) *JSONObj {
	obj.jtype = JSONArray
	if from == nil {
		obj.tArray = make([]*JSONObj, 0)
	} else {
		obj.tArray = from
	}
	return obj
}
func (obj *JSONObj) SetMap(from ObjMap) *JSONObj {
	obj.jtype=  JSONMap
	if from == nil {
		obj.tMap = make(ObjMap)
	} else {
		obj.tMap = from
	}

	return obj
}

func (obj *JSONObj) MakeArray(init ...*JSONObj) (*JSONObj, *[]*JSONObj) {
	if init == nil || len(init)<1 {
		obj.SetArray(nil)
	} else {
		obj.SetArray(init)
	}
	return obj, &obj.tArray
}

func (obj *JSONObj) MakeMap() (*JSONObj,ObjMap) {
	obj.SetMap(nil)
	return obj, obj.tMap
}

func (obj *JSONObj) SetUndef() *JSONObj {
	obj.jtype = JSONUndef
	return obj
}

func (obj *JSONObj) Type() JSONType {
	return obj.jtype
}

func (obj *JSONObj) Value() interface{} {
	switch obj.jtype {
		case JSONUndef:
			return nil
		case JSONInt:
			return obj.tInt
		case JSONString:
			return obj.tString
		case JSONFloat:
			return obj.tFloat
		case JSONBool:
			return obj.tBool
		case JSONArray:
			ar := make([]interface{}, len(obj.tArray))
			for i, v := range obj.tArray {
				ar[i] = v.Value()
			}
			return ar
		case JSONMap:
			mp := make(map[string]interface{})
			for k, v := range obj.tMap {
				mp[k] = v.Value()
			}
			return mp
	}
	return nil
}

func (obj *JSONObj) GetInt() (int64, bool) {
	if obj.jtype == JSONInt {
		return obj.tInt, true
	} else if obj.jtype == JSONFloat {
		return int64(obj.tFloat), true
	} else { return 0, false }
}

func (obj *JSONObj) GetIntOr(or int64) int64 {
	if v, ok := obj.GetInt(); ok {
		return v
	}
	return or
}

func (obj *JSONObj) GetBool() (bool, bool) {
	if obj.jtype == JSONBool {
		return obj.tBool, true
	} else {
		return false,false
	}
}
func (obj *JSONObj) GetBoolOr(or bool) bool {
	if v, ok := obj.GetBool(); ok {
		return v
	}
	return or
}
func (obj *JSONObj) GetString() (string, bool) {
	if obj.jtype == JSONString {
		return obj.tString, true
	} else { return "", false }
}
func (obj *JSONObj) GetStringOr(or string) string {
	if v, ok := obj.GetString(); ok {
		return v
	}
	return or
}
func (obj *JSONObj) GetFloat() (float64, bool) {
	if obj.jtype == JSONFloat {
		return obj.tFloat, true
	} else if obj.jtype == JSONInt {
		return float64(obj.tInt), true
	} else { return 0.00, false }
}
func (obj *JSONObj) GetFloatOr(or float64) float64 {
	if v, ok := obj.GetFloat(); ok {
		return v
	}
	return or
}
func (obj *JSONObj) GetArray() ([]*JSONObj, bool) {
	if obj.jtype == JSONArray {
		return obj.tArray, true
	} else { return nil, false }
}
func (obj *JSONObj) GetArrayOr(or []*JSONObj) []*JSONObj {
	if v, ok := obj.GetArray(); ok {
		return v
	}
	return or
}
func (obj *JSONObj) GetMap() (ObjMap, bool) {
	if obj.jtype == JSONMap {
		return obj.tMap, true
	} else { return nil, false }
}
func (obj *JSONObj) GetMapOr(or ObjMap) ObjMap {
	if v, ok := obj.GetMap(); ok {
		return v
	}
	return or
}
func (obj *JSONObj) IsUndefined() bool {
	return obj.jtype ==JSONUndef
}

func In(mappu map[string]*JSONObj, str string) *JSONObj {
	if val, ok := mappu[str];ok {
		return val
	}
	return NewObj().SetUndef()
}

type TCaseInt		func(int64)
type TCaseString	func(string)
type TCaseFloat		func(float64)
type TCaseBool		func(bool)
type TCaseArray		func([]*JSONObj)
type TCaseMap		func(ObjMap)
type TCaseNull		func()
type TCaseElse		func(*JSONObj)

func (obj *JSONObj) TypeCaseFixed(ci TCaseInt, cs TCaseString, cf TCaseFloat, cb TCaseBool, ca TCaseArray, cm TCaseMap, cn TCaseNull, ce TCaseElse) bool {
	handled := false

	switch obj.jtype {
		case JSONInt:
			if ci != nil {
				vl , _ :=obj.GetInt()
				ci(vl)
				handled = true
			}
		case JSONBool:
			if cb != nil {
				vl, _ :=obj.GetBool()
				cb(vl)
				handled = true
			}
		case JSONFloat:
			if cf != nil {
				vl ,_ := obj.GetFloat()
				cf(vl)
				handled = true
			}
		case JSONString:
			if cs != nil {
				vl ,_:=obj.GetString()
				cs(vl)
				handled = true
			}
		case JSONArray:
			if ca != nil {
				vl, _ := obj.GetArray()
				ca(vl)
				handled = true
			}
		case JSONMap:
			if cm != nil {
				vl ,_ :=obj.GetMap()
				cm(vl)
				handled = true
			}
		default:
			if cn != nil {
				cn()
				handled = true
			}
	}

	if !handled && ce!=nil {
		ce(obj)
		handled=true
	}

	return handled
}

func (obj *JSONObj) TypeCase(rest ...interface{}) bool {
	i := 0

	if rest == nil || len(rest)<1{
		return false
	}

	if len(rest) % 2 != 0 {
		return obj.TypeCase(append(rest, nil)...)
	}

	var ci TCaseInt = nil
	var cf TCaseFloat = nil
	var cb TCaseBool = nil
	var cs TCaseString = nil
	var ca TCaseArray = nil
	var cm TCaseMap = nil
	var cn TCaseNull = nil
	var ce TCaseElse = nil

	for ;i<len(rest);i+=2 {
		ty := rest[i]
		vl := rest[i+1]


		var ok bool = false
		switch ty.(type) {
			case JSONType:
			case func(*JSONObj):
				ce, ok = ty.(func(*JSONObj))
				continue
			default:
				continue
		}

		switch ty.(JSONType) {
			case JSONInt:
				ci, ok = vl.(func(int64))
				if !ok {
					ci32, ok2 := vl.(func(int))
					if ok2 {
						ci = func(sixfour int64) {
							ci32(int(sixfour))
						}
					}
					ok = ok2
				}
			case JSONString:
				cs, ok = vl.(func(string))
			case JSONFloat:
				cf, ok = vl.(func(float64))
				if !ok {
					cf32, ok2 := vl.(func(float32))
					if ok2 {
						cf = func(sixfour float64) {
							cf32(float32(sixfour))
						}
					}
					ok = ok2
				}
			case JSONArray:
				ca, ok = vl.(func([]*JSONObj))
			case JSONMap:
				cm, ok = vl.(func(ObjMap))
				if !ok {
					cmo, ok2 := vl.(func(map[string]*JSONObj))
					if ok2 {
						cm = func(om ObjMap) {
							cmo(ObjMap(om))
						}
					}
					ok = ok2
				}
			case JSONBool:
				cb, ok = vl.(func(bool))
			case JSONUndef:
				cn, ok = vl.(func())
		}
		if !ok {
			ok = ok
		}
	}


	return obj.TypeCaseFixed(ci, cs,cf, cb,ca,cm,cn,ce)
}

func (obj *JSONObj) String() (string, error) {
	by, err := json.Marshal(*obj)

	if err != nil {
		return "undefined", errors.Wrap(err, "json encode failed")
	}

	return string(by), nil
}

func (obj *JSONObj) From(js string) (*JSONObj, error) {
	
	err := json.Unmarshal([]byte(js), obj)

	if err != nil {
		return nil, errors.Wrap(err,"json decode failed")
	}

	return obj, nil

}
