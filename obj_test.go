package jsonobj

import (
	"testing"
	"fmt"
)

func showUnmarshal(obj *JSONObj, indent int) {
	
	tabs := ""

	for i:=0;i<indent;i++ {
		tabs = tabs + "\t"
	}

	fmt.Print(tabs)

 	if !obj.TypeCase(JSONInt, func(in int) {
		fmt.Println("INT", in)
	}, JSONString, func(str string) {
		fmt.Println("STRING", str)
	}, JSONFloat, func(flo float64) {
		fmt.Println("FLOAT", flo)
	}, JSONArray, func(ar []*JSONObj) {
		fmt.Println("[")
		for _, v := range ar {
			showUnmarshal(v, indent+1)
		}
		fmt.Printf("%s]\n", tabs)
	}, JSONMap, func(ar ObjMap) {
		fmt.Println("{")
		for k, v := range ar {

			fmt.Printf("%s %s ->\n", tabs, k)
			showUnmarshal(v, indent+1)
		}
		fmt.Printf("%s}\n", tabs)
	}, JSONUndef, func() {
		fmt.Println("NULL")
	}, JSONBool, func(b bool) {
		fmt.Println("BOOL", b)
	}, func(obj *JSONObj) {
		fmt.Println("ELSE", obj.Value())
	}) {
		panic("uh oh")
	}


}

func TestMessage(t *testing.T) {
	obj, objm := N().MakeMap()

	objm["test"] = N().SetInt(10)
	objm["test2"] = NewObj().SetString("hello")

	var obja *[]*JSONObj
	objm["test3"], obja = NewObj().MakeArray( NewObj().SetString("one"), NewObj().SetInt(2), NewObj().SetFloat(3.141) )

	objm["test4"] = NewObj(123)
	objm["test5"] = NewObj(nil)

	var objm2 ObjMap

	objm["level"], objm2 = NewObj().MakeMap()
	objm2["int"] = NewObj(123)
	objm2["float"] = NewObj(1.23)
	objm2["string"] = NewObj("one two three")
	objm2["array"] = NewObj(1234, 1.234, "one two three four")
	objm2["map"]  = NewObj(Table(KV("one", 1), KV("two", 2.00), KV("three", "three"), KV("four", Table(KV("sub", 1, 2.001, "3", 4)))))

	objm3, _ := objm2["map"].GetMap()
	objm3["thingy"] = NewObj([]string{"one", "two", "THREE"} )

	fmt.Println("get or exists:", objm2["int"].GetIntOr(-100))
	fmt.Println("get or not exists:", In(objm2, "uwu").Value())
	fmt.Println("get or exists 2:", objm2.In("string").Value())
	fmt.Println("get or not exists 2:", objm2.In("@w@").Value())
	fmt.Println("get or wrong type:", objm2["string"].GetFloatOr(1.123))

	objm2["boolean"] = NewObj(true)

	objm2["owo"] = NewObj(Table(KV("one", false, 123, true), KV("two", true)), true, false, "one")

	*obja = append(*obja, NewObj().SetString("world"), NewObj().SetString("foo"), NewObj().SetInt(33))

	str, err  := obj.String()
	
	if err!=nil {
		t.Errorf("Marshalling failed: %v", err)
	}
	
	fmt.Println(str)

	um, err2 := NewObj().From(str)

	if err2 != nil {
		t.Errorf("Unmarshalling failed: %v", err2)
	}

	um2, err3 := um.String()

	if err3 != nil {
		t.Errorf("2nd marshalling failed: %v", err3)
	}

	fmt.Println(um2)

	if str != um2 {
		t.Errorf("Marshalled stirngs are not equal")
	}

	showUnmarshal(obj, 0)

}
