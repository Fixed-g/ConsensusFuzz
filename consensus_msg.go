package tbft

import (
	"chainmaker.org/chainmaker/common/v2/msgbus"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	netpb "chainmaker.org/chainmaker/pb-go/v2/net"
	"github.com/gogo/protobuf/proto"
)

// sendConsensusMsg
// @Description: send consensus msg,If to is an empty string, send to all validators
// @receiver consensus
// @param msg
// @param to
func (consensus *ConsensusTBFTImpl) sendConsensusMsg(msg proto.Message, to string) {
	if msg == nil {
		return
	}

	var validators []string
	if to != "" {
		validators = append(validators, to)
	} else {
		validators = append(validators, consensus.validatorSet.Validators...)
	}

	consensus.logger.Infof("%s ready send consensus message to %v ", consensus.Id, validators)
	for _, v := range validators {
		// The recipient is yourself
		if v == consensus.Id {
			continue
		}
		go func(validator string) {
			netMsg := &netpb.NetMsg{
				Payload: mustMarshal(msg),
				Type:    netpb.NetMsg_CONSENSUS_MSG,
				To:      validator,
			}
			consensus.logger.Infof("%s send consensus message to %s succeeded", consensus.Id, validator)
			consensus.msgbus.Publish(msgbus.SendConsensusMsg, netMsg)
		}(v)
	}
}

// send consensus proposal
func (consensus *ConsensusTBFTImpl) sendConsensusProposal(proposal *TBFTProposal, to string) {
	if proposal == nil || proposal.Bytes == nil {
		return
	}

	// TODO: mutate proposal here
	new_proposal, err := MutateProposal(proposal)
	if err != nil {
		consensus.logger.Errorf(err.Error())
		return
	}

	msg := createProposalTBFTMsg(new_proposal)

	consensus.logger.Infof("%s send consensus proposal", consensus.Id)
	consensus.logger.Debugf("we mutate proposal message in sendConsensusProposal function and send it ")

	// msg := createProposalTBFTMsg(proposal)
	// consensus.logger.Debugf("we don't mutate proposal message in sendConsensusProposal function and send it ")

	consensus.sendConsensusMsg(msg, to)
}

// send consensus vote
// prevote or precommit
func (consensus *ConsensusTBFTImpl) sendConsensusVote(vote *tbftpb.Vote, to string) {
	if vote == nil {
		return
	}
	// TODO: mutate vote here
	new_vote, err := MutateVote(vote)
	if err != nil {
		consensus.logger.Errorf(err.Error())
		return
	}

	// new_vote := vote

	var msg *tbftpb.TBFTMsg
	switch vote.Type {
	case tbftpb.VoteType_VOTE_PREVOTE:
		msg = createPrevoteTBFTMsg(new_vote)
	case tbftpb.VoteType_VOTE_PRECOMMIT:
		msg = createPrecommitTBFTMsg(new_vote)
	}

	consensus.logger.Infof("%s send consensus %s", consensus.Id, vote.String())
	consensus.logger.Debugf("we mutate vote message in sendConsensusVote function and send it ")
	// consensus.logger.Debugf("we dont mutate vote message in sendConsensusVote function and send it ")

	consensus.sendConsensusMsg(msg, to)
}
