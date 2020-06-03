/**
 * App represents a tendermint app
 */
package app

import (
	"bytes"
	"github.com/tendermint/tendermint/libs/log"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type VastchainApplication struct {
	config *Config
	logger log.Logger
}

func NewVastchainApplication(config *Config, logger log.Logger) *VastchainApplication {
	logger.Debug("initialization new vastchain app")

	return &VastchainApplication{
		config: config,
	}
}

var _ abcitypes.Application = (*VastchainApplication)(nil)

func (app *VastchainApplication) Dispose() {
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

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *VastchainApplication) isValid(tx []byte) (code uint32) {
	// check format
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}

	return code
}

func (app *VastchainApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}

func (app *VastchainApplication) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *VastchainApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data

	return
}

func (VastchainApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *VastchainApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {

	return abcitypes.ResponseBeginBlock{}
}

func (VastchainApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}
