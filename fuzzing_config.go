package tbft

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var addr = "./fuzzing_config.yml"

type NodeConfig struct {
	IsFuzzNode bool `yaml:"isFuzzNode"`

	ProcProposeFuzz      bool `yaml:"procProposeFuzz"`
	SendProposeStateFuzz bool `yaml:"sendProposeStateFuzz"`
	CommitBlockFuzz      bool `yaml:"commitBlockFuzz"`
	DelInvalidTxsFuzz    bool `yaml:"delInvalidTxsFuzz"`
	EnterPrevoteFuzz     bool `yaml:"enterPrevoteFuzz"`
	EnterPrecommitFuzz   bool `yaml:"enterPrecommitFuzz"`

	SendStateFuzz            bool `yaml:"sendStateFuzz"`
	SendRoundQCFuzz          bool `yaml:"sendRoundQCFuzz"`
	sendPrecommitInStateFuzz bool `yaml:"sendPrecommitInStateFuzz"`
	sendPrevoteInStateFuzz   bool `yaml:"sendPrevoteInStateFuzz"`
	sendProposalInStateFuzz  bool `yaml:"sendProposalInStateFuzz"`
	sendPrecommitOfRoundFuzz bool `yaml:"sendPrecommitOfRoundFuzz"`
	sendPrevoteOfRoundFuzz   bool `yaml:"sendPrevoteOfRoundFuzz"`
	sendProposalOfRoundFuzz  bool `yaml:"sendProposalOfRoundFuzz"`

	OthersFuzz bool `yaml:"othersFuzz"`

	Delay     bool `yaml:"delay"`
	DelayBase int  `yaml:"delayBase"`
	DelayLim  int  `yaml:"delayLim"`
}

func GetConfig() *NodeConfig {
	config := NodeConfig{}
	dataBytes, err := os.ReadFile(addr)
	if err != nil {
		return nil
	}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		return nil
	}
	return &config
}

func (c *NodeConfig) ToString() string {
	return fmt.Sprintf("FuzzConfig{IsFuzzNode:%t, procProposeFuzz:%t, SendProposeStateFuzz:%t, CommitBlockFuzz:%t, delInvalidTxsFuzz:%t, enterPrevoteFuzz:%t, enterPrecommitFuzz:%t, delay:%t}",
		c.IsFuzzNode, c.ProcProposeFuzz, c.SendProposeStateFuzz, c.CommitBlockFuzz, c.DelInvalidTxsFuzz, c.EnterPrevoteFuzz, c.EnterPrecommitFuzz, c.Delay)
}
