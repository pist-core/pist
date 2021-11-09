// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"git.taiyue.io/pist/go-pist/common"
	"math/big"
)

// NewTxsEvent is posted when a batch of transactions enter the transaction pool.
type NewTxsEvent struct{ Txs []*Transaction }

// NewFastBlocksEvent is posted when a block has been imported.
type NewFastBlocksEvent struct{ FastBlocks []*Block }

// PendingLogsEvent is posted pre mining and notifies of pending logs.
type PendingLogsEvent struct {
	Logs []*Log
}

// PendingStateEvent is posted pre mining and notifies of pending state changes.
type PendingStateEvent struct{}

// RemovedLogsEvent is posted when a reorg happens
type RemovedLogsEvent struct{ Logs []*Log }

type FastChainEvent struct {
	Block *Block
	Hash  common.Hash
	Logs  []*Log
}

type FastChainSideEvent struct {
	Block *Block
}

type FastChainHeadEvent struct{ Block *Block }

type ElectionEvent struct {
	Option           uint
	CommitteeID      *big.Int
	CommitteeMembers []*CommitteeMember
	BackupMembers    []*CommitteeMember
	BeginFastNumber  *big.Int
	EndFastNumber    *big.Int
}

type PbftSignEvent struct {
	Block    *Block
	PbftSign *PbftSign
}

// NodeInfoEvent is posted when nodeInfo send
type NodeInfoEvent struct{ NodeInfo *EncryptNodeMessage }
