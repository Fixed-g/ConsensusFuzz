package fuzzlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestMutateMap(t *testing.T) {
	filePtr, err := os.Open("./test.json")
	if err != nil {
		fmt.Println("file open fail")
		return
	}
	defer filePtr.Close()
	// var info map[string]interface{}
	convertToblock := &common.Block{}

	byteValue, _ := ioutil.ReadAll(filePtr)
	// fmt.Println(decoder)
	err = json.Unmarshal([]byte(byteValue), convertToblock)
	// fmt.Println("info", info)

	if err != nil {
		fmt.Println("解码失败", err.Error())
	} else {
		fmt.Println("解码成功")
	}

	// info, _ = MutateMap(info)
	// fmt.Println(info)
	// err = json.Unmarshal(info, &convertToblock)

	fmt.Println("block_struct")
	// fmt.Println(*convertToblock)

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("convertToMap")
	mapnow, _ := StructToMap(*convertToblock)
	// fmt.Println(mapnow)

	fmt.Println("mutate")
	info, _ := MutateMap(mapnow)
	fmt.Println(info)
}
