package app

import "github.com/incognitochain/incognito-chain/common"

type ShardV2Logger struct {
	log common.Logger
}

func (self *ShardV2Logger) Init(inst common.Logger) {
	self.log = inst
}

// Global instant to use
var Logger = ShardV2Logger{}
