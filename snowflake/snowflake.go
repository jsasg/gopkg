package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	machineBits  int64 = 10
	snBits       int64 = 12
	machineIDMax int64 = -1 ^ (-1 << machineBits)
	snMax        int64 = -1 ^ (-1 << snBits)
	timeLeft     int64 = machineBits + snBits
	machineLeft  int64 = snBits
)

var initTime int64 = time.Now().UnixMilli()

type SnowFlake struct {
	mut       sync.Mutex
	timestamp int64
	machineID int64
	sn        int64
}

func New(machineID int64) (*SnowFlake, error) {
	if machineID >= machineIDMax || machineID < 0 {
		return nil, errors.New("机器ID过大或不能小于0")
	}
	return &SnowFlake{
		timestamp: 0,
		machineID: machineID,
		sn:        0,
	}, nil
}

func (s *SnowFlake) Next() (id int64) {
	s.mut.Lock()
	defer s.mut.Unlock()
	now := time.Now().UnixMilli()
	if s.timestamp == now {
		s.sn++
		if s.sn > snMax {
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sn = 0
		s.timestamp = now
	}
	id = ((now - initTime) << timeLeft) | (s.machineID << machineLeft) | (s.sn)
	return
}
