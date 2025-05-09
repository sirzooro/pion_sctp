// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sctp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testChunkReconfigParamA() []byte {
	return []byte{
		0x00, 0x0d, 0x00, 0x16, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06,
	}
}

func testChunkReconfigParamB() []byte {
	return []byte{0x0, 0xd, 0x0, 0x10, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3}
}

func TestParamOutgoingResetRequest_Success(t *testing.T) {
	tt := []struct {
		binary []byte
		parsed *paramOutgoingResetRequest
	}{
		{
			testChunkReconfigParamA(),
			&paramOutgoingResetRequest{
				paramHeader: paramHeader{
					typ: outSSNResetReq,
					len: 22,
					raw: testChunkReconfigParamA()[4:],
				},
				reconfigRequestSequenceNumber:  1,
				reconfigResponseSequenceNumber: 2,
				senderLastTSN:                  3,
				streamIdentifiers:              []uint16{4, 5, 6},
			},
		},
		{
			testChunkReconfigParamB(),
			&paramOutgoingResetRequest{
				paramHeader: paramHeader{
					typ: outSSNResetReq,
					len: 16,
					raw: testChunkReconfigParamB()[4:],
				},
				reconfigRequestSequenceNumber:  1,
				reconfigResponseSequenceNumber: 2,
				senderLastTSN:                  3,
				streamIdentifiers:              []uint16{},
			},
		},
	}

	for i, tc := range tt {
		actual := &paramOutgoingResetRequest{}
		_, err := actual.unmarshal(tc.binary)
		assert.NoErrorf(t, err, "failed to unmarshal #%d", i)
		assert.Equal(t, tc.parsed, actual)

		b, err := actual.marshal()
		assert.NoErrorf(t, err, "failed to marshal #%d", i)
		assert.Equal(t, tc.binary, b)
	}
}

func TestParamOutgoingResetRequest_Failure(t *testing.T) {
	tt := []struct {
		name   string
		binary []byte
	}{
		{"packet too short", testChunkReconfigParamA()[:8]},
		{"param too short", []byte{0x0, 0xd, 0x0, 0x4}},
	}

	for i, tc := range tt {
		actual := &paramOutgoingResetRequest{}
		_, err := actual.unmarshal(tc.binary)
		assert.Errorf(t, err, "expected unmarshal #%d: '%s' to fail.", i, tc.name)
	}
}
