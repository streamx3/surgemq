// Copyright (c) 2014 The SurgeMQ Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package message

// PingRespMessage A PINGRESP Packet is sent by the Server to the Client in response to a PINGREQ
// Packet. It indicates that the Server is alive.
type PingRespMessage struct {
	header
}

var _ Provider = (*PingRespMessage)(nil)

// NewPingRespMessage creates a new PINGRESP message.
func NewPingRespMessage() Provider {
	msg := &PingRespMessage{}
	msg.SetType(PINGRESP) // nolint: errcheck

	return msg
}

// Decode message
func (dm *PingRespMessage) Decode(src []byte) (int, error) {
	return dm.header.decode(src)
}

// Encode message
func (dm *PingRespMessage) Encode(dst []byte) (int, error) {
	if !dm.dirty {
		if len(dst) < len(dm.dBuf) {
			return 0, ErrInsufficientBufferSize
		}

		return copy(dst, dm.dBuf), nil
	}

	return dm.header.encode(dst)
}
