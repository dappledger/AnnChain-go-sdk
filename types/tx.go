// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
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

package types

import (
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/dappledger/ann-go-sdk/go-hash"
	"github.com/dappledger/ann-go-sdk/merkle"
)

func BinRead(reader io.Reader, data interface{}) error {
	return binary.Read(reader, binary.BigEndian, data)
}

func BinWrite(writer io.Writer, data interface{}) error {
	return binary.Write(writer, binary.BigEndian, data)
}

// ErrUnexpectedEOF when it is the end
func ReadBytes(reader io.Reader) ([]byte, error) {
	byLen, err := ReadVarint(reader)
	if err != nil {
		return nil, err
	}
	bys := make([]byte, byLen)
	err = BinRead(reader, &bys)
	return bys, err
}

type Tx []byte

// NOTE: this is the hash of the go-wire encoded Tx.
// Tx has no types at this level, so just length-prefixed.
// Alternatively, it may make sense to add types here and let
// []byte be type 0x1 so we can have versioned txs if need be in the future.

// ethereum transaction hash
func (tx Tx) Hash() []byte {
	return hash.Keccak256Func(tx)
}

type Txs []Tx

func (txs Txs) Hash() []byte {
	// Recursive impl.
	// Copied from go-merkle to avoid allocations
	switch len(txs) {
	case 0:
		return nil
	case 1:
		return txs[0].Hash()
	default:
		left := Txs(txs[:(len(txs)+1)/2]).Hash()
		right := Txs(txs[(len(txs)+1)/2:]).Hash()
		return merkle.SimpleHashFromTwoHashes(left, right)
	}
}

func (tx Tx) Size() int {
	return 1
}

func WrapTx(prefix []byte, tx []byte) []byte {
	return append(prefix, tx...)
}

func UnwrapTx(tx []byte) []byte {
	if len(tx) > 4 {
		return tx[4:]
	}
	return tx
}

func (txs Txs) ToBytes() [][]byte {
	return *((*[][]byte)(unsafe.Pointer(&txs)))
}
