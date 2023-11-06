package fuzzlib

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/consensus"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	tbft "github.com/fixed-g/consensusfuzz/v2"
)

func MutateBool(b bool) bool {
	return generate_random_bool()
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

func MutateProposal(proposal *tbft.TBFTProposal) (*tbft.TBFTProposal, error) {
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

func MutateVoteMsg(msg *tbft.ConsensusMsg) (*tbft.ConsensusMsg, error) {
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

func MutateProposalMsg(msg *tbft.ConsensusMsg) (*tbft.ConsensusMsg, error) {
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
