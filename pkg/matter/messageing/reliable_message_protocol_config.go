package messageing

import (
	"sync"
	"time"
)

var ChipConfigMrpDefaultIdleRetryInterval int64 = 5000
var ChipConfigMrpDefaultActiveRetryInterval int64 = 300

type ReliableMessageProtocolConfig struct {
	mIdleRetransTimeout   time.Duration
	mActiveRetransTimeout time.Duration
}

var insRMPC *ReliableMessageProtocolConfig
var rmpcOnce = sync.Once{}

func GetLocalMRPConfig() *ReliableMessageProtocolConfig {
	rmpcOnce.Do(func() {
		insRMPC = newReliableMessageProtocolConfig()
		insRMPC.mIdleRetransTimeout = time.Duration(ChipConfigMrpDefaultIdleRetryInterval)
		insRMPC.mActiveRetransTimeout = time.Duration(ChipConfigMrpDefaultActiveRetryInterval)
	})
	return insRMPC
}

func newReliableMessageProtocolConfig() *ReliableMessageProtocolConfig {
	return &ReliableMessageProtocolConfig{}
}
