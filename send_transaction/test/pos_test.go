package test

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"git.taiyue.io/pist/go-pist/core"
	"git.taiyue.io/pist/go-pist/core/state"
	"git.taiyue.io/pist/go-pist/core/types"
	"git.taiyue.io/pist/go-pist/core/vm"
	"git.taiyue.io/pist/go-pist/log"
	"git.taiyue.io/pist/go-pist/params"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))
}

///////////////////////////////////////////////////////////////////////
func TestGetLockedAsset(t *testing.T) {
	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(number uint64, gen *core.BlockGen, blockchain *core.BlockChain, header *types.Header, statedb *state.StateDB) {
		sendTranction(number, gen, statedb, mAccount, saddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)

		sendDepositTransaction(number, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendCancelTransaction(number-types.GetEpochFromID(2).BeginHeight, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		if number == 90 {
			stateDb := gen.GetStateDB()
			impawn := vm.NewImpawnImpl()
			impawn.Load(stateDb, types.StakingAddress)
			arr := impawn.GetLockedAsset(saddr1)
			for addr, value := range arr {
				fmt.Println("value ", value.Value, " addr ", addr.String())
			}
		}

		sendWithdrawTransaction(number-types.MinCalcRedeemHeight(2), gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
	}
	manager := newTestPOSManager(55, executable)
	fmt.Println(" saddr1 ", manager.GetBalance(saddr1))
	//epoch  [id:1,begin:1,end:2000]   [id:2,begin:2001,end:4000]   [id:3,begin:4001,end:6000]
	//epoch  [id:2,begin:2001,end:4000]   5002
}

func TestFeeAndPK(t *testing.T) {
	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(number uint64, gen *core.BlockGen, blockchain *core.BlockChain, header *types.Header, statedb *state.StateDB) {
		sendTranction(number, gen, statedb, mAccount, saddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)

		sendDepositTransaction(number, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)

		sendUpdateFeeTransaction(number, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendUpdatePkTransaction(number, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)

		sendCancelTransaction(number-types.GetEpochFromID(2).BeginHeight, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)

		sendWithdrawTransaction(number-types.MinCalcRedeemHeight(2), gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
	}
	manager := newTestPOSManager(101, executable)
	fmt.Println(" saddr1 ", manager.GetBalance(saddr1))
	//epoch  [id:1,begin:1,end:2000]   [id:2,begin:2001,end:4000]   [id:3,begin:4001,end:6000]
	//epoch  [id:2,begin:2001,end:4000]   5002
}

///////////////////////////////////////////////////////////////////////
func TestDeposit(t *testing.T) {
	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(number uint64, gen *core.BlockGen, blockchain *core.BlockChain, header *types.Header, statedb *state.StateDB) {
		sendTranction(number, gen, statedb, mAccount, saddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)
		sendTranction(number, gen, statedb, mAccount, daddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)

		sendDepositTransaction(number, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendDelegateTransaction(number-60, gen, daddr1, saddr1, big.NewInt(1000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendCancelTransaction(number-types.GetEpochFromID(2).BeginHeight, gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendUnDelegateTransaction(number-types.GetEpochFromID(2).BeginHeight-10, gen, daddr1, saddr1, big.NewInt(1000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		if number == 130 {
			stateDb := gen.GetStateDB()
			impawn := vm.NewImpawnImpl()
			impawn.Load(stateDb, types.StakingAddress)
			arr := impawn.GetLockedAsset(saddr1)
			for addr, value := range arr {
				fmt.Println("value ", value.Value, " addr ", addr.String())
			}
			arr1 := impawn.GetLockedAsset(daddr1)
			for addr, value := range arr1 {
				fmt.Println("value D ", value.Value, " addr ", addr.String())
			}
		}
		sendWithdrawTransaction(number-types.MinCalcRedeemHeight(2), gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendWithdrawDelegateTransaction(number-types.MinCalcRedeemHeight(2)-10, gen, daddr1, saddr1, big.NewInt(1000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
	}
	manager := newTestPOSManager(55, executable)
	fmt.Println(" saddr1 ", types.ToPist(manager.GetBalance(saddr1)), " StakingAddress ", manager.GetBalance(types.StakingAddress), " ", types.ToPist(manager.GetBalance(types.StakingAddress)))
}

func TestDepositGetDeposit(t *testing.T) {
	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(number uint64, gen *core.BlockGen, blockchain *core.BlockChain, header *types.Header, statedb *state.StateDB) {
		sendTranction(number, gen, statedb, mAccount, saddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)
		sendTranction(number, gen, statedb, mAccount, daddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)

		sendDepositTransaction(number, gen, saddr1, big.NewInt(4000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-51, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)
		sendDelegateTransaction(number, gen, daddr1, saddr1, big.NewInt(4000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-61, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)

		sendCancelTransaction(number-types.GetEpochFromID(2).BeginHeight, gen, saddr1, big.NewInt(3000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-types.GetEpochFromID(2).BeginHeight-11, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)

		sendUnDelegateTransaction(number-types.GetEpochFromID(2).BeginHeight, gen, daddr1, saddr1, big.NewInt(3000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-types.GetEpochFromID(2).BeginHeight-21, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)

		if number == 130 {
			stateDb := gen.GetStateDB()
			impawn := vm.NewImpawnImpl()
			impawn.Load(stateDb, types.StakingAddress)
			arr := impawn.GetLockedAsset(saddr1)
			for addr, value := range arr {
				fmt.Println("value ", value.Value, " addr ", addr.String())
			}
			arr1 := impawn.GetLockedAsset(daddr1)
			for addr, value := range arr1 {
				fmt.Println("value D ", value.Value, " addr ", addr.String())
			}
		}
		sendWithdrawTransaction(number-types.MinCalcRedeemHeight(2), gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-types.MinCalcRedeemHeight(2)-11, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)
		sendWithdrawDelegateTransaction(number-types.MinCalcRedeemHeight(2), gen, daddr1, saddr1, big.NewInt(1000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-types.MinCalcRedeemHeight(2)-21, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)
	}
	manager := newTestPOSManager(101, executable)
	fmt.Println(" saddr1 ", types.ToPist(manager.GetBalance(saddr1)), " StakingAddress ", manager.GetBalance(types.StakingAddress), " ", types.ToPist(manager.GetBalance(types.StakingAddress)))
}

///////////////////////////////////////////////////////////////////////
func TestSendTX(t *testing.T) {
	params.MinTimeGap = big.NewInt(0)

	genesis := gspec.MustFastCommit(db)
	blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{})
	chain, _ := core.GenerateChain(gspec.Config, genesis, engine, db, 20000, func(i int, gen *core.BlockGen) {
		header := gen.GetHeader()
		statedb := gen.GetStateDB()
		if i > 60 {
			SendTX(header, true, blockchain, nil, gspec.Config, gen, statedb, nil)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		panic(err)
	}
}

///////////////////////////////////////////////////////////////////////
func TestRewardTime(t *testing.T) {
	params.MinTimeGap = big.NewInt(0)
	params.ElectionMinLimitForStaking = new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18))

	genesis := gspec.MustFastCommit(db)
	blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{})
	parentFast := genesis
	delegateNum = 50000

	for i := 0; i < 505; i++ {

		chain, _ := core.GenerateChain(gspec.Config, parentFast, engine, db, 60, func(i int, gen *core.BlockGen) {
			header := gen.GetHeader()
			stateDB := gen.GetStateDB()
			if header.Number.Uint64() > 60 {
				SendTX(header, true, blockchain, nil, gspec.Config, gen, stateDB, nil)
			}
		})
		if _, err := blockchain.InsertChain(chain); err != nil {
			panic(err)
		}
		parentFast = blockchain.CurrentBlock()
	}
}

func TestDelegateRewardNextEpochValid(t *testing.T) {
	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(number uint64, gen *core.BlockGen, blockchain *core.BlockChain, header *types.Header, statedb *state.StateDB) {
		sendTranction(number, gen, statedb, mAccount, saddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)
		sendTranction(number, gen, statedb, mAccount, daddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)

		sendDepositTransaction(number, gen, saddr1, big.NewInt(4000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-51, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)

		sendCancelTransaction(number-types.GetEpochFromID(2).BeginHeight, gen, saddr1, big.NewInt(3000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-types.GetEpochFromID(2).BeginHeight-11, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)

		sendWithdrawTransaction(number-types.MinCalcRedeemHeight(2), gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-types.MinCalcRedeemHeight(2)-11, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)

		sendDelegateTransaction(number-params.NewEpochLength, gen, daddr1, saddr1, big.NewInt(4000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-61-params.NewEpochLength, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)

		sendUnDelegateTransaction(number-types.GetEpochFromID(3).BeginHeight, gen, daddr1, saddr1, big.NewInt(3000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-types.GetEpochFromID(3).BeginHeight-21, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)

		i := number / params.NewEpochLength
		if number == 130+params.NewEpochLength*i {
			impawn := vm.NewImpawnImpl()
			impawn.Load(statedb, types.StakingAddress)
			arr := impawn.GetLockedAsset(saddr1)
			for addr, value := range arr {
				fmt.Println("value ", value.Value, " addr ", addr.String())
			}
			arr1 := impawn.GetLockedAsset(daddr1)
			for addr, value := range arr1 {
				fmt.Println("value D ", value.Value, " addr ", addr.String(), "balance", statedb.GetBalance(daddr1), "number", number)
			}
		}

		sendWithdrawDelegateTransaction(number-types.MinCalcRedeemHeight(3), gen, daddr1, saddr1, big.NewInt(1000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-types.MinCalcRedeemHeight(3)-21, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)
	}
	manager := newTestPOSManager(101, executable)
	fmt.Println(" saddr1 ", types.ToPist(manager.GetBalance(saddr1)), " StakingAddress ", manager.GetBalance(types.StakingAddress), " ", types.ToPist(manager.GetBalance(types.StakingAddress)))
}

func TestDelegateCancleInUnSelectValidator(t *testing.T) {
	StakerValidNumber := uint64(60)
	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(number uint64, gen *core.BlockGen, blockchain *core.BlockChain, header *types.Header, statedb *state.StateDB) {
		sendTranction(number, gen, statedb, mAccount, saddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)
		sendTranction(number, gen, statedb, mAccount, daddr1, big.NewInt(6000000000000000000), priKey, signer, nil, header)

		sendDepositTransaction(number, gen, saddr1, big.NewInt(4000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-51, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)

		sendCancelTransaction(number-StakerValidNumber, gen, saddr1, big.NewInt(3000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-StakerValidNumber-11, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)

		sendWithdrawTransaction(number-types.MinCalcRedeemHeight(2), gen, saddr1, big.NewInt(1000000000000000000), skey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDepositTransaction(number-types.MinCalcRedeemHeight(2)-11, gen, saddr1, skey1, signer, statedb, blockchain, abiStaking, nil)

		sendDelegateTransaction(number, gen, daddr1, saddr1, big.NewInt(4000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-61, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)

		sendUnDelegateTransaction(number-StakerValidNumber, gen, daddr1, saddr1, big.NewInt(3000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-StakerValidNumber-21, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)

		sendWithdrawDelegateTransaction(number-types.MinCalcRedeemHeight(2), gen, daddr1, saddr1, big.NewInt(1000000000000000000), dkey1, signer, statedb, blockchain, abiStaking, nil)
		sendGetDelegateTransaction(number-types.MinCalcRedeemHeight(2)-21, gen, daddr1, saddr1, dkey1, signer, statedb, blockchain, abiStaking, nil)
	}
	manager := newTestPOSManager(101, executable)
	fmt.Println(" saddr1 ", types.ToPist(manager.GetBalance(saddr1)), " StakingAddress ", manager.GetBalance(types.StakingAddress), " ", types.ToPist(manager.GetBalance(types.StakingAddress)))
}
