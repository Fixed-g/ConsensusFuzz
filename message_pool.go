package tbft

import (
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
)

type MessageQueue struct {
	messages []*ConsensusMsg
}

type MessagePool struct {
	cap  int
	pool map[string]*MessageQueue
}

func newMessagePool(cap int, msgs []*ConsensusMsg) *MessagePool {
	var messagePoolTemp = &MessagePool{}

	messagePoolTemp.cap = cap
	messagePoolTemp.pool = make(map[string]*MessageQueue)

	for i := range msgs {
		addMsgToPool(messagePoolTemp, msgs[i])
	}

	return messagePoolTemp
}

// add consensus message to message pool
func addMsgToPool(messagePool *MessagePool, message *ConsensusMsg) *MessagePool {
	if _, ok := messagePool.pool[message.Type.String()]; !ok {
		queue_tmp := &MessageQueue{
			messages: make([]*ConsensusMsg, 0),
		}
		messagePool.pool[message.Type.String()] = queue_tmp
	}

	messagePool.pool[message.Type.String()].messages = append(messagePool.pool[message.Type.String()].messages, message)

	for len(messagePool.pool[message.Type.String()].messages) > messagePool.cap {
		messagePool.pool[message.Type.String()].messages = messagePool.pool[message.Type.String()].messages[1:]
	}
	return messagePool
}

// find the latest message of certain type
func findLatestMsgWithType(messagePool *MessagePool, msgType tbftpb.TBFTMsgType) *ConsensusMsg {
	if val, ok := messagePool.pool[msgType.String()]; ok {
		for i := len(val.messages) - 1; i >= 0; i-- {
			if val.messages[i].Type == msgType {
				return val.messages[i]
			}
		}
	}

	return nil
}

// display the info of message pool
func displayMsgPool(consensus *ConsensusTBFTImpl, messagePool *MessagePool) {
	consensus.logger.Debugf("-----------------------------------------")
	consensus.logger.Debugf("cap of message pool: %d", messagePool.cap)

	for k, v := range messagePool.pool {
		consensus.logger.Debugf("type of this que is: %s", k)
		for i := 0; i < len(v.messages); i++ {
			consensus.logger.Debugf("[%d] message type [%s] ", i, v.messages[i].Type.String())
		}
	}

	consensus.logger.Debugf("-----------------------------------------")
}
