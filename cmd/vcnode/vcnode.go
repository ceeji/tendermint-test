package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	cfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	vcApp "vastchain.ltd/vastchain/modules/app"
)

var dataDir string

func initCLI() {
	// get default dataDir
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	defaultDataDir := filepath.Join(currentDir, "data")

	// read CLI options
	pDataDir := flag.String("dataDir", defaultDataDir, "the path to the root directory of data (defaults to 'data' directory in current directory)")
	flag.Parse()
	dataDir = *pDataDir

	// create data directory if it does not exist
	_ = os.MkdirAll(dataDir, os.ModePerm)
	_ = os.MkdirAll(filepath.Join(dataDir, "store"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(dataDir, "config"), os.ModePerm)
}

func run() error {
	runtime.GOMAXPROCS(128)

	initCLI()

	node, err := newVcNode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return err
	}

	err = node.Start()
	if err != nil {
		return err
	}
	defer func() {
		node.Stop()
		node.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}

func main() {
	err := run()

	if err != nil {
		os.Exit(2)
	}
}

// newVcNode creates a new vastchain node (VcNode)
func newVcNode() (*nm.Node, error) {
	// read config
	appConfig, err := vcApp.ReadInConfig(dataDir)
	if err != nil {
		panic(err)
	}
	tendermintConfig, err := appConfig.GetTendermintConfig()
	if err != nil {
		panic(err)
	}

	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger, err = tmflags.ParseLogLevel(tendermintConfig.LogLevel, logger, cfg.DefaultLogLevel())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}

	// read private validator
	pv := privval.LoadFilePV(
		tendermintConfig.PrivValidatorKeyFile(),
		tendermintConfig.PrivValidatorStateFile(),
	)

	// read node key
	nodeKey, err := p2p.LoadOrGenNodeKey(tendermintConfig.NodeKeyFile())
	if err != nil {
		return nil, errors.Wrap(err, "failed to load node's key")
	}

	// create app
	app := vcApp.NewVastchainApplication(appConfig, logger)
	defer app.Dispose()

	// create node
	node, err := nm.NewNode(
		tendermintConfig,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		nm.DefaultGenesisDocProviderFunc(tendermintConfig),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(tendermintConfig.Instrumentation),
		logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new Vastchain node")
	}

	return node, nil
}
