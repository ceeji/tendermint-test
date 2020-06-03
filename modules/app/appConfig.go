package app

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	"os"
	"path/filepath"
	"time"
	"vastchain.ltd/vastchain/utils"
)

// AppConfig represents the configuration of VastChainApp.
type Config struct {
	// dataDir saves all files including configuration and data
	dataDir string
}

func ReadInConfig(dataDir string) (*Config, error) {
	configFile := filepath.Join(dataDir, "config", "config.toml")
	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		// use default config
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "viper failed to unmarshal config")
	}

	config.dataDir = dataDir
	return &config, nil
}

func (conf *Config) GetTendermintConfig() (*cfg.Config, error) {
	// read config
	config := cfg.DefaultConfig()
	config.RootDir = conf.dataDir

	config.PrivValidatorKey = filepath.Join(conf.dataDir, "config", "priv_validator_key.json")
	config.PrivValidatorState = filepath.Join(conf.dataDir, "config", "priv_validator_state.json")
	config.NodeKey = filepath.Join(conf.dataDir, "config", "node_key.json")
	config.Consensus.SetWalFile(filepath.Join(conf.dataDir, "store", "cs.wal", "wal"))
	config.Consensus.CreateEmptyBlocksInterval = time.Second * 60
	config.P2P.AddrBookStrict = false
	config.P2P.AddrBook = filepath.Join(conf.dataDir, "store", "addrbook.json")
	config.Genesis = filepath.Join(conf.dataDir, "config", "vastchain-genesis.json")
	config.RPC.MaxBodyBytes = int64(10240000) // 10 MB

	// generating priv validator key if it does not exist
	var filePV *privval.FilePV
	if _, err := os.Stat(config.PrivValidatorKey); os.IsNotExist(err) {
		filePV = privval.GenFilePV(config.PrivValidatorKey, config.PrivValidatorState)
		filePV.Save()
	} else {
		filePV = privval.LoadFilePV(config.PrivValidatorKey, config.PrivValidatorState)
	}

	// generate genesis file if it does not exist
	if _, err := os.Stat(config.Genesis); os.IsNotExist(err) {
		err = generateGenesisFile(config.Genesis, []crypto.PubKey{filePV.GetPubKey()})
		if err != nil {
			return nil, err
		}
	}

	if err := config.ValidateBasic(); err != nil {
		panic(errors.Wrap(err, "config is invalid"))
	}

	return config, nil
}

// generateGenesisFile generates a genesis file if it does not exist or do nothing.
func generateGenesisFile(genFile string, validators []crypto.PubKey) error {
	// genesis file
	if _, err := os.Stat(genFile); err == nil {
		// if file existed, do nothing and do not throw an error
		return nil
	} else {
		genDoc := types.GenesisDoc{
			ChainID:         fmt.Sprintf("vastchain-test-%v", utils.RandomStr(6)),
			GenesisTime:     time.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
		}

		genDoc.Validators = []types.GenesisValidator{}
		for _, validator := range validators {
			genDoc.Validators = append(genDoc.Validators, types.GenesisValidator{
				Address: validator.Address(),
				PubKey:  validator,
				Power:   10,
			})
		}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
	}

	return nil
}
