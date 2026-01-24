package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/pkg/api/jsonrpc"
	"github.com/apex/pkg/consensus"
	"github.com/apex/pkg/core"
	"github.com/apex/pkg/storage"
	"github.com/apex/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile string
	logger  *zap.Logger
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "apex",
		Short: "Apex Blockchain Node",
		Long:  "Apex Blockchain - A high-performance Layer-1 blockchain with DPoS consensus",
		Run:   runNode,
	}
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	
	// Initialize logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err != nil {
		logger.Warn("Failed to read config file", zap.Error(err))
	}
}

func runNode(cmd *cobra.Command, args []string) {
	logger.Info("Starting Apex Blockchain Node")
	
	// Initialize database
	dbPath := viper.GetString("storage.path")
	if dbPath == "" {
		dbPath = "./data/apex.db"
	}
	
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		logger.Fatal("Failed to open database", zap.Error(err))
	}
	defer db.Close()
	
	// Initialize state DB
	stateDB := storage.NewStateDB(db)
	
	// Initialize block store
	blockStore := storage.NewBlockStore(db)
	
	// Initialize DPoS consensus
	dpos := consensus.NewDPoS()
	
	// Initialize blockchain
	blockchain := core.NewBlockchain(stateDB, blockStore, dpos)
	
	// Check if genesis exists
	height := blockchain.GetHeight()
	if height == 0 {
		logger.Info("Initializing genesis block")
		if err := initGenesis(blockchain); err != nil {
			logger.Fatal("Failed to initialize genesis", zap.Error(err))
		}
	}
	
	logger.Info("Blockchain initialized", zap.Uint64("height", blockchain.GetHeight()))
	
	// Start JSON-RPC server
	rpcPort := viper.GetInt("rpc.port")
	if rpcPort == 0 {
		rpcPort = 8545
	}
	
	rpcServer := jsonrpc.NewServer(blockchain, rpcPort, logger)
	go func() {
		if err := rpcServer.Start(); err != nil {
			logger.Error("RPC server error", zap.Error(err))
		}
	}()
	
	logger.Info("Apex node running", zap.Int("rpc_port", rpcPort))
	
	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	logger.Info("Shutting down Apex node")
}

func initGenesis(blockchain *core.Blockchain) error {
	// Create genesis validators
	genesisValidators := make([]*types.Validator, 0)
	
	// You would load these from genesis.json in production
	// For now, create a default validator
	
	// Create genesis accounts with initial distribution
	totalSupply := types.ToWei(float64(types.TotalSupply))
	genesisAccounts := []*types.Account{
		{
			Address: types.HexToAddress("0x0000000000000000000000000000000000000001"),
			Balance: totalSupply,
			Nonce:   0,
			Staked:  types.ToWei(0),
			Locked:  types.ToWei(0),
		},
	}
	
	return blockchain.InitGenesis(genesisValidators, genesisAccounts)
}
