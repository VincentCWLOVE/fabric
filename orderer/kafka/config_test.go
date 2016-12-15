/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kafka

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hyperledger/fabric/orderer/localconfig"
	cb "github.com/hyperledger/fabric/protos/common"
	ab "github.com/hyperledger/fabric/protos/orderer"
)

var (
	testBrokerID     = int32(0)
	testOldestOffset = int64(100)                                    // The oldest block available on the broker
	testNewestOffset = int64(1100)                                   // The offset that will be assigned to the next block
	testMiddleOffset = (testOldestOffset + testNewestOffset - 1) / 2 // Just an offset in the middle

	// Amount of time to wait for block processing when doing time-based tests
	// We generally want this value to be as small as possible so as to make tests execute faster
	// But this may have to be bumped up in slower machines
	testTimePadding = 200 * time.Millisecond
)

var testConf = &config.TopLevel{
	General: config.General{
		OrdererType:   "kafka",
		LedgerType:    "ram",
		BatchTimeout:  500 * time.Millisecond,
		BatchSize:     100,
		QueueSize:     100,
		MaxWindowSize: 100,
		ListenAddress: "127.0.0.1",
		ListenPort:    7050,
		GenesisMethod: "provisional",
	},
	Kafka: config.Kafka{
		Brokers: []string{"127.0.0.1:9092"},
		Retry: config.Retry{
			Period: 3 * time.Second,
			Stop:   60 * time.Second,
		},
		Verbose: false,
		Version: sarama.V0_9_0_1,
	},
}

func testClose(t *testing.T, x Closeable) {
	if err := x.Close(); err != nil {
		t.Fatal("Cannot close mock resource:", err)
	}
}

func newTestEnvelope(content string) *cb.Envelope {
	return &cb.Envelope{Payload: []byte(content)}
}

func testNewSeekMessage(startLabel string, seekNo, windowNo uint64) *ab.DeliverUpdate {
	var startVal ab.SeekInfo_StartType
	switch startLabel {
	case "oldest":
		startVal = ab.SeekInfo_OLDEST
	case "newest":
		startVal = ab.SeekInfo_NEWEST
	default:
		startVal = ab.SeekInfo_SPECIFIED

	}
	return &ab.DeliverUpdate{
		Type: &ab.DeliverUpdate_Seek{
			Seek: &ab.SeekInfo{
				Start:           startVal,
				SpecifiedNumber: seekNo,
				WindowSize:      windowNo,
			},
		},
	}
}

func testNewAckMessage(ackNo uint64) *ab.DeliverUpdate {
	return &ab.DeliverUpdate{
		Type: &ab.DeliverUpdate_Acknowledgement{
			Acknowledgement: &ab.Acknowledgement{
				Number: ackNo,
			},
		},
	}
}
