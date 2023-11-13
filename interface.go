package tbft

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/consensus"
	consensuspb "chainmaker.org/chainmaker/pb-go/v2/consensus"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	"github.com/mitchellh/mapstructure"
)

func MapToBlock(m map[string]interface{}) *common.Block {
	block := &common.Block{}
	err := mapstructure.Decode(m, block)
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

func MapToVoteMsg(m map[string]interface{}) *ConsensusMsg {
	vote := &tbftpb.Vote{}
	err := mapstructure.Decode(m["Msg"].(map[string]interface{}), vote)
	if err != nil {
		return nil
	}
	return &ConsensusMsg{
		Type: m["Type"].(tbftpb.TBFTMsgType),
		Msg:  mustMarshal(vote),
	}
}

func MapToProposalMsg(m map[string]interface{}) *ConsensusMsg {
	proposal := &tbftpb.Proposal{}
	err := mapstructure.Decode(m["Msg"].(map[string]interface{}), proposal)
	if err != nil {
		return nil
	}
	return &ConsensusMsg{
		Type: m["Type"].(tbftpb.TBFTMsgType),
		Msg:  mustMarshal(proposal),
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

func MutateBool(b bool) bool {
	return Generate_random_bool()
}

func MutateBlock(block *common.Block) (*common.Block, error) {
	var err error
	var block_map map[string]interface{}
	block_map, err = StructToMap(block)
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
	proposal_map, err = StructToMap(proposal)
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
	vote_map, err = StructToMap(vote)
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
	txs_map, err = StructToMap(txs)
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

func MutateVoteMsg(msg *ConsensusMsg) (*ConsensusMsg, error) {
	var err error
	var msg_map map[string]interface{}
	msg_map, err = StructToMap(msg)
	if err != nil {
		return nil, err
	}
	msg_map, err = MutateMap(msg_map)
	if err != nil {
		return nil, err
	}
	msg = MapToVoteMsg(msg_map)
	return msg, nil
}

func MutateProposalMsg(msg *ConsensusMsg) (*ConsensusMsg, error) {
	var err error
	var msg_map map[string]interface{}
	msg_map, err = StructToMap(msg)
	if err != nil {
		return nil, err
	}
	msg_map, err = MutateMap(msg_map)
	if err != nil {
		return nil, err
	}
	msg = MapToProposalMsg(msg_map)
	return msg, nil
}
