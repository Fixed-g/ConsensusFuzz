package fuzzlib

import (
	tbft "Fixed-g/ConsensusFuzz/v0.1"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	consensuspb "chainmaker.org/chainmaker/pb-go/v2/consensus"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	netpb "chainmaker.org/chainmaker/pb-go/v2/net"
)

const (
	Int32 = iota
	Uint64
	String
	byteArray
	stringArray
)

type Map map[string]*Value

type Value struct {
	Type    string
	Payload interface{}
}

func UnresolvedPointer(payload interface{}) *Value {
	return &Value{
		Type:    "Pointer",
		Payload: payload,
	}
}

func Int32ToValue(payload int32) *Value {
	return &Value{
		Type:    "int32",
		Payload: payload,
	}
}

func UInt64ToValue(payload uint64) *Value {
	return &Value{
		Type:    "uint64",
		Payload: payload,
	}
}

func StringToValue(payload string) *Value {
	return &Value{
		Type:    "string",
		Payload: payload,
	}
}

func ByteArrayToValue(payload []byte) *Value {
	return &Value{
		Type:    "[]byte",
		Payload: payload,
	}
}

func StringArrayToValue(payload []string) *Value {
	return &Value{
		Type:    "[]string",
		Payload: payload,
	}
}

func VoteToValue(vote *tbftpb.Vote) *Value {
	return &Value{
		Type: "map",
		Payload: Map{
			"Type":        Int32ToValue(int32(vote.Type)),
			"Voter":       StringToValue(vote.Voter),
			"Height":      UInt64ToValue(vote.Height),
			"Round":       Int32ToValue(vote.Round),
			"Hash":        ByteArrayToValue(vote.Hash),
			"InvalidTxs":  StringArrayToValue(vote.InvalidTxs),
			"Endorsement": UnresolvedPointer(vote.Endorsement),
		},
	}
}

func BlockToValue(block *common.Block) *Value {
	return &Value{
		Type: "map_Block",
		Payload: Map{
			"Header":         UnresolvedPointer(block.Header),
			"Dag":            UnresolvedPointer(block.Dag),
			"Txs":            UnresolvedPointer(block.Txs),
			"AdditionalData": UnresolvedPointer(block.AdditionalData),
		},
	}
}

func RwSetVerifyFailTxsToValue(txs *consensuspb.RwSetVerifyFailTxs) *Value {
	return &Value{
		Type: "map_RwSetVerifyFailTxs",
		Payload: Map{
			"BlockHeight": UInt64ToValue(txs.BlockHeight),
			"TxIds":       StringArrayToValue(txs.TxIds),
		},
	}
}

func ConsensusMsgToValue(consensusMsg *tbft.ConsensusMsg) *Value {
	var res = &Value{
		Type: "map_ConsensusMsg",
		Payload: Map{
			"Type": Int32ToValue(int32(consensusMsg.Type)),
		},
	}
	if consensusMsg.Type == tbftpb.TBFTMsgType_MSG_PREVOTE || consensusMsg.Type == tbftpb.TBFTMsgType_MSG_PRECOMMIT {
		res.Payload.(Map)["Msg"] = VoteToValue(consensusMsg.Msg.(*tbftpb.Vote))
	}
	return res
}

func NetMsgToValue(netMsg netpb.NetMsg) *Value {
	return &Value{
		Type: "map_netMsg",
		Payload: Map{
			"Type":    Int32ToValue(int32(netMsg.Type)),
			"Payload": UnresolvedPointer(netMsg.Payload),
			"To":      StringToValue(netMsg.To),
		},
	}
}

func TBFTMsgToValue(msg tbftpb.TBFTMsg) *Value {
	var res = &Value{
		Type: "map_TBFTMsg",
		Payload: Map{
			"Type": Int32ToValue(int32(msg.Type)),
		},
	}
	if msg.Type == tbftpb.TBFTMsgType_MSG_PROPOSE {
		res.Payload.(Map)["Msg"] = UnresolvedPointer(msg.Msg)
	} else if msg.Type == tbftpb.TBFTMsgType_MSG_PREVOTE || msg.Type == tbftpb.TBFTMsgType_MSG_PRECOMMIT {
		res.Payload.(Map)["Msg"] = UnresolvedPointer(msg.Msg)

	}
	return res
}
