/**
 * App represents a tendermint app
 */
package app

import (
	"bytes"
	"github.com/tendermint/tendermint/libs/log"
	"path/filepath"
	"vastchain.ltd/vastchain/modules/kvstore"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type VastchainApplication struct {
	config       *Config
	logger       log.Logger
	db           *kvstore.DB
	currentBatch *kvstore.Tx
}

func NewVastchainApplication(config *Config, logger log.Logger) *VastchainApplication {
	logger.Debug("initialization new vastchain app")
	db, err := kvstore.NewDB(filepath.Join(config.dataDir, "store"))
	if err != nil {
		panic(err)
	}

	return &VastchainApplication{
		config: config,
		db:     db,
	}
}

var _ abcitypes.Application = (*VastchainApplication)(nil)

func (app *VastchainApplication) Dispose() {
	app.db.Dispose()
}

func (VastchainApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (VastchainApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *VastchainApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	parts := bytes.Split(req.Tx, []byte("="))
	key, value := parts[0], parts[1]

	err := app.currentBatch.Set("global", key, value)
	if err != nil {
		panic(err)
	}

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *VastchainApplication) isValid(tx []byte) (code uint32) {
	// check format
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}

	key, value := parts[0], parts[1]

	// check if the same key=value already exists
	err := app.db.View(func(txn *kvstore.Tx) error {
		item, err := txn.Get("global", key)
		if err != nil && err != kvstore.ErrKeyNotFound {
			return err
		}
		if err == nil {
			if bytes.Equal(item, value) {
				code = 2
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return code
}

func (app *VastchainApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

func (app *VastchainApplication) Commit() abcitypes.ResponseCommit {
	app.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *VastchainApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *kvstore.Tx) error {
		item, err := txn.Get("global", reqQuery.Data)
		if err != nil && err != kvstore.ErrKeyNotFound {
			return err
		}
		if err == kvstore.ErrKeyNotFound {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
			resQuery.Value = item
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func (VastchainApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *VastchainApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (VastchainApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}
