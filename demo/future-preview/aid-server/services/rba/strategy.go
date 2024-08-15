package rba

import (
	"aid-server/pkg/timestamp"
	"aid-server/services/user"
)

type Strategy interface {
	// Verify Risk Based Authentication implementation
	Verify(data user.IUser, input *user.Record) bool
}

type SimpleStrategy struct {
}

var SimpleAlgo SimpleStrategy

func init() {
	SimpleAlgo = SimpleStrategy{}
}

func (s *SimpleStrategy) Verify(userItem user.IUser, input *user.Record) bool {
	if !timestamp.CheckTimestampClose5000(userItem.GetTime().CurEventTime, input.Time.CurEventTime) {
		return false
	}
	if *userItem.GetSpace() != input.Space {
		return false
	}
	return true
}
