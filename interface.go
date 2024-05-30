package tbft

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/consensus"
	consensuspb "chainmaker.org/chainmaker/pb-go/v2/consensus"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	"github.com/mitchellh/mapstructure"
)

func MapToBlock(m map[string]interface{}) *common.Block {
	block := &common.Block{}
	bt, err := json.Marshal(m)
	err = json.Unmarshal(bt, block)
	if err != nil {
		return nil
	}
	return block
}

func MapToProposal(m map[string]interface{}) *TBFTProposal {
	proposal := &TBFTProposal{}
	err := mapstructure.Decode(m, proposal)
	if err != nil {
		return nil
	}
	return proposal
}

func MapToVote(m map[string]interface{}) *tbftpb.Vote {
	vote := &tbftpb.Vote{}
	err := mapstructure.Decode(m, vote)
	if err != nil {
		return nil
	}
	return vote
}

func MapToVoteMsg(m map[string]interface{}, t tbftpb.TBFTMsgType) *ConsensusMsg {
	vote := &tbftpb.Vote{}
	err := mapstructure.Decode(m["Msg"].(map[string]interface{}), vote)
	if err != nil {
		return nil
	}
	return &ConsensusMsg{
		// Type: tbftpb.TBFTMsgType(m["Type"].(tbftpb.TBFTMsgType)),
		Type: t,
		Msg:  vote,
	}
}

func MapToProposalMsg(m map[string]interface{}, t tbftpb.TBFTMsgType) *ConsensusMsg {
	proposal := &tbftpb.Proposal{}
	err := mapstructure.Decode(m["Msg"].(map[string]interface{}), proposal)
	if err != nil {
		return nil
	}
	return &ConsensusMsg{
		// Type: tbftpb.TBFTMsgType(m["Type"].(tbftpb.TBFTMsgType)),
		Type: t,
		Msg:  proposal,
	}
}

func MapToTxs(m map[string]interface{}) *consensuspb.RwSetVerifyFailTxs {
	txs := &consensuspb.RwSetVerifyFailTxs{}
	err := mapstructure.Decode(m, txs)
	if err != nil {
		return nil
	}
	return txs
}

func MapToGossipState(m map[string]interface{}) *tbftpb.GossipState {
	state := &tbftpb.GossipState{}
	err := mapstructure.Decode(m, state)
	if err != nil {
		return nil
	}
	return state
}

func MapToRoundQC(m map[string]interface{}) *tbftpb.RoundQC {
	roundQC := &tbftpb.RoundQC{}
	err := mapstructure.Decode(m, roundQC)
	if err != nil {
		return nil
	}
	return roundQC
}

func MutateBool(b bool) bool {
	return Generate_random_bool()
}

func MutateBlock(block *common.Block) (*common.Block, error) {
	var err error
	var block_map map[string]interface{}
	block_map, err = StructToMap(*block)
	if err != nil {
		return nil, err
	}
	block_map, err = MutateMap(block_map)
	if err != nil {
		return nil, err
	}
	block = MapToBlock(block_map)
	return block, nil
}

func MutateProposal(proposal *TBFTProposal) (*TBFTProposal, error) {
	var err error
	var proposal_map map[string]interface{}
	proposal_map, err = StructToMap(*proposal)
	if err != nil {
		return nil, err
	}
	proposal_map, err = MutateMap(proposal_map)
	if err != nil {
		return nil, err
	}
	proposal = MapToProposal(proposal_map)
	return proposal, nil
}

func MutateVote(vote *tbftpb.Vote) (*tbftpb.Vote, error) {
	var err error
	var vote_map map[string]interface{}
	vote_map, err = StructToMap(*vote)
	if err != nil {
		return nil, err
	}
	vote_map, err = MutateMap(vote_map)
	if err != nil {
		return nil, err
	}
	vote = MapToVote(vote_map)
	return vote, nil
}

func MutateTxs(txs *consensus.RwSetVerifyFailTxs) (*consensus.RwSetVerifyFailTxs, error) {
	var err error
	var txs_map map[string]interface{}
	txs_map, err = StructToMap(*txs)
	if err != nil {
		return nil, err
	}
	txs_map, err = MutateMap(txs_map)
	if err != nil {
		return nil, err
	}
	txs = MapToTxs(txs_map)
	return txs, nil
}

// 需要指定consensusMsg类型
func MutateVoteMsg(msg *ConsensusMsg) (*ConsensusMsg, error) {
	var err error
	var msg_map map[string]interface{}
	msg_map, err = StructToMap(*msg)

	if err != nil {
		return nil, err
	}
	msg_map, err = MutateMap(msg_map)
	if err != nil {
		return nil, err
	}
	msg = MapToVoteMsg(msg_map, msg.Type)
	return msg, nil
}

func MutateProposalMsg(msg *ConsensusMsg) (*ConsensusMsg, error) {
	var err error
	var msg_map map[string]interface{}
	msg_map, err = StructToMap(*msg)
	if err != nil {
		return nil, err
	}
	msg_map, err = MutateMap(msg_map)
	if err != nil {
		return nil, err
	}
	msg = MapToProposalMsg(msg_map, msg.Type)
	return msg, nil
}

// Broadcaster Mutations:
func MutateGossipState(state *tbftpb.GossipState) (*tbftpb.GossipState, error) {
	var err error
	var state_map map[string]interface{}
	state_map, err = StructToMap(*state)
	if err != nil {
		return nil, err
	}
	state_map, err = MutateMap(state_map)
	if err != nil {
		return nil, err
	}
	state = MapToGossipState(state_map)
	return state, nil
}

func MutateRoundQC(roundQC *tbftpb.RoundQC) (*tbftpb.RoundQC, error) {
	var err error
	var roundQC_map map[string]interface{}
	roundQC_map, err = StructToMap(*roundQC)
	if err != nil {
		return nil, err
	}
	roundQC_map, err = MutateMap(roundQC_map)
	if err != nil {
		return nil, err
	}
	roundQC = MapToRoundQC(roundQC_map)
	return roundQC, nil
}

const (
	tagMeta = "meta"
)

// StructToMap convert struct to map[string]interface{}
// if value of field is nil, result doesn't contain related field and value
// for example: input is struct {a: 1, b: "abc", c: nil},
// convert to a--1, b--"abc", result doesn't contain field c
// notice: field must be exported, unexported field will panic
func StructToMap(config interface{}) (map[string]interface{}, error) {
	if config == nil {
		return nil, nil
	}
	if reflect.TypeOf(config).Kind() != reflect.Struct {
		return nil, errors.New("incorrect config type: config type should be struct")
	}
	result := make(map[string]interface{})
	configValue := reflect.ValueOf(config)
	for i := 0; i < configValue.NumField(); i++ {
		field := reflect.TypeOf(config).Field(i)
		if !parseMetaTag(field) {
			continue
		}
		configField := field.Name
		rv := configValue.Field(i)
		// changed here
		if rv.Kind() == reflect.Interface {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Int:
			result[configField] = rv.Interface()
		case reflect.Int8:
			result[configField] = rv.Interface()
		case reflect.Int16:
			result[configField] = rv.Interface()
		case reflect.Int32:
			result[configField] = rv.Interface()
		case reflect.Int64:
			result[configField] = rv.Interface()
		case reflect.Uint:
			result[configField] = rv.Interface()
		case reflect.Uint8:
			result[configField] = rv.Interface()
		case reflect.Uint16:
			result[configField] = rv.Interface()
		case reflect.Uint32:
			result[configField] = rv.Interface()
		case reflect.Uint64:
			result[configField] = rv.Interface()
		case reflect.Float32:
			result[configField] = rv.Interface()
		case reflect.Float64:
			result[configField] = rv.Interface()
		case reflect.String:
			result[configField] = rv.Interface()
		case reflect.Bool:
			result[configField] = rv.Interface()
		case reflect.Ptr:
			v, err := parsePtr(rv)
			if err != nil {
				errMsg := errors.New(fmt.Sprintf("structToMap fail, field is %s, value is %v, err is %s",
					configField, rv, err))
				return nil, errMsg
			}
			if v == nil {
				continue
			}
			result[configField] = v
		case reflect.Map:
			v, err := parseMap(rv)
			if err != nil {
				return nil, err
			}
			result[configField] = v
		case reflect.Slice:
			v, err := parseSlice(rv)
			if err != nil {
				errMsg := errors.New(fmt.Sprintf("structToMap fail, field is %s, value is %v, err is %s",
					configField, rv, err))
				return nil, errMsg
			}
			if v == nil {
				continue
			}
			result[configField] = v
		case reflect.Struct:
			v, err := StructToMap(rv.Interface())
			if err != nil {
				return nil, err
			}
			result[configField] = v
		default:
			errMsg := fmt.Sprintf("structToMap fail, unknow value type, type is %s, value is %v\n",
				rv.Kind(), rv)
			return nil, errors.New(errMsg)
		}
	}
	return result, nil
}

func parsePtr(v reflect.Value) (map[string]interface{}, error) {
	if v.IsNil() {
		return nil, nil
	}
	v = v.Elem()
	if !v.IsValid() {
		return nil, nil
	}
	return StructToMap(v.Interface())
}

func parseMap(v reflect.Value) (map[string]interface{}, error) {
	if v.IsNil() {
		return nil, nil
	}
	result := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		// changed here
		rv := v.MapIndex(key)
		if rv.Kind() == reflect.Interface {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Int:
			result[key.String()] = rv.Interface()
		case reflect.Int8:
			result[key.String()] = rv.Interface()
		case reflect.Int16:
			result[key.String()] = rv.Interface()
		case reflect.Int32:
			result[key.String()] = rv.Interface()
		case reflect.Int64:
			result[key.String()] = rv.Interface()
		case reflect.Uint:
			result[key.String()] = rv.Interface()
		case reflect.Uint8:
			result[key.String()] = rv.Interface()
		case reflect.Uint16:
			result[key.String()] = rv.Interface()
		case reflect.Uint32:
			result[key.String()] = rv.Interface()
		case reflect.Uint64:
			result[key.String()] = rv.Interface()
		case reflect.Float32:
			result[key.String()] = rv.Interface()
		case reflect.Float64:
			result[key.String()] = rv.Interface()
		case reflect.String:
			result[key.String()] = rv.Interface()
		case reflect.Bool:
			result[key.String()] = rv.Interface()
		case reflect.Ptr:
			v, err := parsePtr(rv)
			if err != nil {
				errMsg := errors.New(fmt.Sprintf("structToMap fail, key is %s, value is %v, err is %s",
					key.String(), rv, err))
				return nil, errMsg
			}
			if v == nil {
				continue
			}
			result[key.String()] = v
		case reflect.Map:
			v, err := parseMap(rv)
			if err != nil {
				return nil, err
			}
			result[key.String()] = v
		case reflect.Slice:
			v, err := parseSlice(rv)
			if err != nil {
				return nil, err
			}
			result[key.String()] = v
		case reflect.Struct:
			v, err := StructToMap(rv.Interface())
			if err != nil {
				return nil, err
			}
			result[key.String()] = v
		default:
			errMsg := fmt.Sprintf("structToMap fail, unknow value type, type is %s, value is %v\n",
				rv.Kind(), rv)
			return nil, errors.New(errMsg)
		}
		result[key.String()] = rv.Interface()
	}
	return result, nil
}

func parseSlice(v reflect.Value) ([]interface{}, error) {
	if v.Type().Kind() != reflect.Slice {
		return nil, errors.New("incorrect config type: config type should be slice")
	}

	if v.Len() <= 0 || !v.IsValid() || v.IsNil() {
		return nil, nil
	}

	res := make([]interface{}, v.Len())

	for i := 0; i < v.Len(); i++ {
		switch v.Index(i).Type().Kind() {
		case reflect.Int:
			res[i] = v.Index(i).Interface()
		case reflect.Int8:
			res[i] = v.Index(i).Interface()
		case reflect.Int16:
			res[i] = v.Index(i).Interface()
		case reflect.Int32:
			res[i] = v.Index(i).Interface()
		case reflect.Int64:
			res[i] = v.Index(i).Interface()
		case reflect.Uint:
			res[i] = v.Index(i).Interface()
		case reflect.Uint8:
			res[i] = v.Index(i).Interface()
		case reflect.Uint16:
			res[i] = v.Index(i).Interface()
		case reflect.Uint32:
			res[i] = v.Index(i).Interface()
		case reflect.Uint64:
			res[i] = v.Index(i).Interface()
		case reflect.Float32:
			res[i] = v.Index(i).Interface()
		case reflect.Float64:
			res[i] = v.Index(i).Interface()
		case reflect.String:
			res[i] = v.Index(i).Interface()
		case reflect.Bool:
			res[i] = v.Index(i).Interface()
		case reflect.Ptr:
			value, err := parsePtr(v.Index(i))
			if err != nil {
				return nil, err
			}
			res[i] = value
		case reflect.Struct:
			curStruct, err := StructToMap(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			res[i] = curStruct
		default:
			errMsg := fmt.Sprintf("unknow slice type %s", v.Index(i).Type().Kind().String())
			return nil, errors.New(errMsg)
		}
	}
	return res, nil
}

func parseMetaTag(f reflect.StructField) bool {
	metaTag := f.Tag.Get(tagMeta)
	parseCurrentField, err := strconv.ParseBool(metaTag)
	if err != nil {
		return true
	}
	if parseCurrentField {
		return true
	}
	return false
}
