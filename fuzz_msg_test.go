package tbft

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"os"
	"testing"
)

func TestFunc(t *testing.T) {
	msg := &ConsensusMsg{
		Type: tbftpb.TBFTMsgType_MSG_PREVOTE,
		Msg: &tbftpb.Vote{
			Type:        1,
			Voter:       "2",
			Height:      3,
			Round:       4,
			Hash:        make([]byte, 5),
			InvalidTxs:  make([]string, 6),
			Endorsement: nil,
		},
	}
	vote_map := &map[string]interface{}{
		"Type":        1,
		"Voter":       "2",
		"Height":      3,
		"Round":       4,
		"Hash":        make([]byte, 5),
		"InvalidTxs":  make([]string, 6),
		"Endorsement": nil,
	}
	proposal_map := &map[string]interface{}{
		"Voter":    "1",
		"Height":   2,
		"Round":    3,
		"PolRound": 4,
		"Block": &map[string]interface{}{
			"Header":         &map[string]interface{}{},
			"Dag":            nil,
			"Txs":            nil,
			"AdditionalData": nil,
		},
		"Endorsement": nil,
		"TxsRwSet":    nil,
		"Qc":          nil,
	}
	vote := &tbftpb.Vote{}
	json_proposal := &tbftpb.Proposal{}
	ms_proposal := &tbftpb.Proposal{}
	fmt.Println("proposal_map:", proposal_map)
	bt, err := json.Marshal(proposal_map)
	err = mapstructure.Decode(proposal_map, ms_proposal)
	fmt.Println("ms_proposal:", ms_proposal)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("proposal_bytes:", string(bt))
	err = json.Unmarshal(bt, json_proposal)
	if err != nil {
		return
	}
	fmt.Println("json_proposal:", json_proposal)
	*proposal_map, err = StructToMap(*json_proposal)
	fmt.Println("proposal_map:", proposal_map)
	msg = msg
	vote = vote
	vote_map = vote_map
}

func TestBlock(t *testing.T) {
	jsonFile, err := os.Open("block.json")
	out, _ := os.Create("out.txt")
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			return
		}
	}(jsonFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	block := &common.Block{}
	err = json.Unmarshal([]byte(byteValue), block)
	if err != nil {
		return
	}
	//_, _ = fmt.Fprintln(out, block)
	msg_map, _ := StructToMap(*block)
	fmt.Println(msg_map)
	msg_json, err := json.Marshal(msg_map)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = fmt.Fprintln(out, string(msg_json))
}
