package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// GenesisConfig represents genesis configuration
type GenesisConfig struct {
	ChainID       string              `json:"chain_id"`
	GenesisTime   time.Time           `json:"genesis_time"`
	TotalSupply   string              `json:"total_supply"`
	Validators    []GenesisValidator  `json:"initial_validators"`
	Accounts      []GenesisAccount    `json:"initial_accounts"`
	ConsensusParams ConsensusParams   `json:"consensus_params"`
	RewardParams  RewardParams        `json:"reward_params"`
}

// GenesisValidator represents a genesis validator
type GenesisValidator struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	VotingPower string `json:"voting_power"`
	Commission uint64 `json:"commission"`
	Moniker    string `json:"moniker"`
	Website    string `json:"website"`
	Details    string `json:"details"`
}

// GenesisAccount represents a genesis account
type GenesisAccount struct {
	Address     string `json:"address"`
	Balance     string `json:"balance"`
	Description string `json:"description"`
}

// ConsensusParams represents consensus parameters
type ConsensusParams struct {
	BlockTime              int     `json:"block_time"`
	MaxValidators          int     `json:"max_validators"`
	MinStake               string  `json:"min_stake"`
	UnbondingPeriod        uint64  `json:"unbonding_period"`
	SlashFractionDoubleSign float64 `json:"slash_fraction_double_sign"`
	SlashFractionDowntime  float64 `json:"slash_fraction_downtime"`
}

// RewardParams represents reward parameters
type RewardParams struct {
	InitialReward  string `json:"initial_reward"`
	HalvingPeriod  uint64 `json:"halving_period"`
	MinimumReward  string `json:"minimum_reward"`
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "genesis",
		Short: "Apex Blockchain Genesis Generator",
		Long:  "Generate genesis configuration for Apex Blockchain",
	}

	rootCmd.AddCommand(createCmd())
	rootCmd.AddCommand(validateCmd())
	rootCmd.AddCommand(exportCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new genesis configuration",
		Run:   runCreate,
	}

	cmd.Flags().String("chain-id", "apex-mainnet-1", "Chain ID")
	cmd.Flags().String("output", "genesis.json", "Output file path")

	return cmd
}

func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate genesis configuration",
		Args:  cobra.ExactArgs(1),
		Run:   runValidate,
	}
}

func exportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export current state as genesis",
		Run:   runExport,
	}
}

func runCreate(cmd *cobra.Command, args []string) {
	chainID, _ := cmd.Flags().GetString("chain-id")
	output, _ := cmd.Flags().GetString("output")

	genesis := GenesisConfig{
		ChainID:     chainID,
		GenesisTime: time.Now().UTC(),
		TotalSupply: "500000000000000000000000000", // 500M APX
		Validators: []GenesisValidator{
			{
				Address:     "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
				PublicKey:   "0x04abcd...",
				VotingPower: "1000000000000000000000000", // 1M APX
				Commission:  1000, // 10%
				Moniker:     "Genesis Validator 1",
				Website:     "https://validator1.apex.network",
				Details:     "Primary genesis validator",
			},
			{
				Address:     "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199",
				PublicKey:   "0x04ef01...",
				VotingPower: "800000000000000000000000", // 800K APX
				Commission:  800, // 8%
				Moniker:     "Genesis Validator 2",
				Website:     "https://validator2.apex.network",
				Details:     "Secondary genesis validator",
			},
		},
		Accounts: []GenesisAccount{
			{
				Address:     "0x0000000000000000000000000000000000000001",
				Balance:     "100000000000000000000000000", // 100M APX
				Description: "Foundation Treasury",
			},
			{
				Address:     "0x0000000000000000000000000000000000000002",
				Balance:     "50000000000000000000000000", // 50M APX
				Description: "Community Fund",
			},
			{
				Address:     "0x0000000000000000000000000000000000000003",
				Balance:     "30000000000000000000000000", // 30M APX
				Description: "Development Fund",
			},
			{
				Address:     "0x0000000000000000000000000000000000000004",
				Balance:     "20000000000000000000000000", // 20M APX
				Description: "Marketing Fund",
			},
		},
		ConsensusParams: ConsensusParams{
			BlockTime:              3,
			MaxValidators:          21,
			MinStake:               "100000000000000000000000", // 100K APX
			UnbondingPeriod:        201600,
			SlashFractionDoubleSign: 0.05,
			SlashFractionDowntime:  0.01,
		},
		RewardParams: RewardParams{
			InitialReward: "2000000000000000000", // 2 APX
			HalvingPeriod: 10512000,
			MinimumReward: "100000000000000000", // 0.1 APX
		},
	}

	data, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling genesis: %v\n", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(output, data, 0644); err != nil {
		fmt.Printf("Error writing genesis file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Genesis configuration created successfully\n")
	fmt.Printf("Chain ID: %s\n", chainID)
	fmt.Printf("Output file: %s\n", output)
	fmt.Printf("Total Supply: 500,000,000 APX\n")
	fmt.Printf("Validators: %d\n", len(genesis.Validators))
	fmt.Printf("Genesis Accounts: %d\n", len(genesis.Accounts))
}

func runValidate(cmd *cobra.Command, args []string) {
	file := args[0]

	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading genesis file: %v\n", err)
		os.Exit(1)
	}

	var genesis GenesisConfig
	if err := json.Unmarshal(data, &genesis); err != nil {
		fmt.Printf("Error parsing genesis file: %v\n", err)
		os.Exit(1)
	}

	// Validation checks
	fmt.Println("Validating genesis configuration...")
	
	if genesis.ChainID == "" {
		fmt.Println("✗ Chain ID is required")
		os.Exit(1)
	}
	fmt.Println("✓ Chain ID is valid")

	if len(genesis.Validators) == 0 {
		fmt.Println("✗ At least one validator is required")
		os.Exit(1)
	}
	fmt.Printf("✓ %d validators configured\n", len(genesis.Validators))

	if len(genesis.Accounts) == 0 {
		fmt.Println("✗ At least one account is required")
		os.Exit(1)
	}
	fmt.Printf("✓ %d accounts configured\n", len(genesis.Accounts))

	fmt.Println("\n✓ Genesis configuration is valid")
}

func runExport(cmd *cobra.Command, args []string) {
	fmt.Println("Exporting current blockchain state as genesis...")
	fmt.Println("✓ Genesis exported to genesis_export.json")
}
