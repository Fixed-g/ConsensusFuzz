package tbft

import (
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestMutateMap(t *testing.T) {
	filePtr, err := os.Open("./block.json")
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

func TestNilPanic(t *testing.T) {
	block := &common.Block{}
	bt, err := json.Marshal(nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(bt, block)
	return
}

func TestGossipState(t *testing.T) {
	proposal := []byte("123")
	roundVoteSet := &tbftpb.RoundVoteSet{
		Height:     1,
		Round:      2,
		Prevotes:   nil,
		Precommits: nil,
	}
	gossipProto := &tbftpb.GossipState{
		Id:               "1",
		Height:           2,
		Round:            3,
		Step:             4,
		Proposal:         proposal,
		VerifingProposal: proposal,
		RoundVoteSet:     roundVoteSet,
	}
	m, err := StructToMap(*gossipProto)
	fmt.Println(m)
	if err != nil {
		fmt.Println(err.Error())
	}
	new_state := MapToGossipState(m)
	new_m, err := StructToMap(*new_state)
	fmt.Println(new_m)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type A struct {
	Val int
}

type B struct {
	A map[string]*A
}

type C struct {
	B map[string]*B
}

func TestMutatePointerMap(t *testing.T) {
	aa := &A{
		Val: 1,
	}
	bb := &B{
		A: map[string]*A{
			"aa": aa,
		},
	}
	cc := &C{
		B: map[string]*B{
			"bb": bb,
		},
	}

	ident := reflect.ValueOf(*cc).Kind()
	fmt.Println(ident)

	m, err := StructToMap(*cc)
	fmt.Println(m)
	if err != nil {
		fmt.Println(err)
	}
	m, err = MutateMap(m)
	fmt.Println(m)
	if err != nil {
		fmt.Println(err)
	}
}
