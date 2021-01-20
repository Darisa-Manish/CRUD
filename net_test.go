package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"reflect"
	"testing"
)

/*func TestGetName(t *testing.T){
	testcases:=[]struct{
		input string
		output []customer
	}{
		{"Customer1",[]customer{{1,"Customer1","02/11/1999",Address{1,"K.P.T","DMM","A.P",1}}}},
		{"Customer2",[]customer{{2,"Customer2","02/11/1999",Address{2,"K.P.T","DMM","A.P",2}}}},
		{"", []customer{{1, "Customer1", "02/11/1999", Address{1, "K.P.T", "DMM", "A.P", 1}}, {2, "Customer2", "02/11/1999", Address{2, "K.P.T", "DMM", "A.P", 2}}}},
		{"Customer6",[]customer(nil)},
	}
	for i:=range testcases{
		var c []customer
		reqst := "http://192.168.1.223:8080/customer?name=" + testcases[i].input
		req := httptest.NewRequest("GET", reqst, nil)
		w := httptest.NewRecorder()
		GetName(w, req)
		//resp := w.Result()
		//val, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(w.Body.Bytes(), &c)
		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(c, testcases[i].output) {
			t.Errorf(`FAIL: %v Expected ans: %v Got: %v`, testcases[i].input, testcases[i].output, c)
		}
	}
}*/
/*func TestGetID(t *testing.T){
	testcases:=[]struct{
		input int
		output customer
	}{
		{1,customer{1,"CustomerA","02/11/1999",Address{1,"K.P.T","DMM","A.P",1}}},
		{2,customer{2,"CustomerB","02/11/1999",Address{2,"K.P.T","DMM","A.P",2}}},
		{7,customer{}},
		//{"", []customer{{1, "Customer1", "02/11/1999", Address{1, "K.P.T", "DMM", "A.P", 1}}, {2, "Customer2", "02/11/1999", Address{2, "K.P.T", "DMM", "A.P", 2}}}},
	}
	for i:=range testcases{
		var c customer
		reqst := "http://localhost:8080/customer/"
		req := httptest.NewRequest("GET", reqst, nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(testcases[i].input)})
		w := httptest.NewRecorder()
		GetID(w, req)
		resp := w.Result()
		val, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println("val is:",string(val))
		err := json.Unmarshal(val, &c)
		//fmt.Println("c is",c)
		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(c, testcases[i].output) {
			t.Errorf(`FAIL: %v Expected ans: %v Got: %v`, testcases[i].input, testcases[i].output, c)
		}
	}
}*/
func TestPost(t *testing.T){
	testcases:=[]struct{
		input customer
		output customer
	}{
		{customer{3,"Customer3","02/11/1999",Address{3,"K.P.T","DMM","A.P",3}},customer{3,"Customer3","02/11/1999",Address{3,"K.P.T","DMM","A.P",3}}},
		{customer{4,"Customer4","02/11/1999",Address{4,"K.P.T","DMM","A.P",4}},customer{4,"Customer4","02/11/1999",Address{4,"K.P.T","DMM","A.P",4}}},
		{customer{5,"Customer5","02/11/2005",Address{5,"K.P.T","DMM","A.P",5}},customer{}},
		//{customer{{3,"Customer3","02/11/1999",Address{3,"K.P.T","DMM","A.P",3}}},[]customer{{1,"Customer1","02/11/1999",Address{1,"K.P.T","DMM","A.P",1}}}},
		//{customer{{4,"Customer4","02/11/1999",Address{4,"K.P.T","DMM","A.P",4}}},[]customer{{1,"Customer1","02/11/1999",Address{1,"K.P.T","DMM","A.P",1}}}},
		//{"", []customer{{1, "Customer1", "02/11/1999", Address{1, "K.P.T", "DMM", "A.P", 1}}, {2, "Customer2", "02/11/1999", Address{2, "K.P.T", "DMM", "A.P", 2}}}},
	}
	for i:=range testcases{
		var c customer
		byte,_:=json.Marshal(testcases[i].input)
		reqst := "http://192.168.1.223:8080/customer"
		req := httptest.NewRequest("POST", reqst, bytes.NewBuffer(byte))
		w := httptest.NewRecorder()
		postcustomer(w, req)
		resp := w.Result()
		val, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(val, &c)
		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(c, testcases[i].output) {
			t.Errorf(`FAIL: %v Expected ans: %v Got: %v`, testcases[i].input, testcases[i].output, c)
		}
	}
}
/*func TestPut(t *testing.T){
	testcases:=[]struct{
		input customer
		output customer
	}{
		{customer{1,"CustomerA","",Address{1,"K.P.T","ATP","AP",1}},customer{1,"CustomerA","",Address{1,"K.P.T","ATP","AP",1}}},
		{customer{2,"CustomerB","",Address{}},customer{2,"CustomerB","",Address{}}},
		{customer{2,"CustomerC","03/11/1999",Address{}},customer{2,"CustomerB","02/11/1999",Address{2,"K.P.T","DMM","A.P",2}}},
	}
	for i:=range testcases{
		var c customer
		request:="http://192.168.1.223/8080/customer/"
		byte,_:=json.Marshal(testcases[i].input)
		req:=httptest.NewRequest("PUT",request,bytes.NewBuffer(byte))
		req = mux.SetURLVars(req, map[string]string{"id":strconv.Itoa(testcases[i].input.ID)})
		w:=httptest.NewRecorder()
		putcustomer(w,req)
		resp:=w.Result()
		val,_:=ioutil.ReadAll(resp.Body)
		fmt.Println("val is",string(val))
		err:=json.Unmarshal(val,&c)
		if err!=nil{
			log.Fatal("2 Error")
		}
		if !reflect.DeepEqual(testcases[i].output,c){
			t.Errorf(`FAIL: %v Expected ans: %v Got: %v`, testcases[i].input, testcases[i].output, c)
		}
	}
}*/
/*func TestDelete(t *testing.T){
	testcases:=[]struct{
		input int
		output customer
	}{
		{3,customer{3,"Customer3","02/11/1999",Address{3,"K.P.T","DMM","A.P",3}}},
		//{"4",customer{4,"Customer4","02/11/1999",Address{4,"K.P.T","DMM","A.P",4}}},
		//{5,customer{}},
	}
	for i:=range testcases{
		var c customer
		request:="http://192.168.1.223/8080/customer/"
		req:=httptest.NewRequest("DELETE",request,nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(testcases[i].input)})
		w:=httptest.NewRecorder()
		delcust(w,req)
		resp:=w.Result()
		val,_:=ioutil.ReadAll(resp.Body)
		err:=json.Unmarshal(val,&c)
		if err!=nil{
			log.Fatal(err)
		}
		if !reflect.DeepEqual(c,testcases[i].output){
			t.Errorf(`FAIL: %v Expected ans: %v Got: %v`, testcases[i].input, testcases[i].output, c)
		}
	}
}*/