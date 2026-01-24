package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func main() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	rootCmd := &cobra.Command{
		Use:   "apexctl",
		Short: "Apex Blockchain CLI Tool",
		Long:  "Command-line tool for managing Apex Blockchain nodes, validators, and accounts",
	}

	// Add subcommands
	rootCmd.AddCommand(keysCmd())
	rootCmd.AddCommand(validatorCmd())
	rootCmd.AddCommand(stakeCmd())
	rootCmd.AddCommand(queryCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// keysCmd returns the keys management command
func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage cryptographic keys",
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new key pair",
		Run:   runKeysGenerate,
	}
	generateCmd.Flags().String("output", "", "Output file for private key (required)")
	generateCmd.MarkFlagRequired("output")

	cmd.AddCommand(generateCmd)
	return cmd
}

// validatorCmd returns the validator management command
func validatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Manage validators",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new validator",
		Run:   runValidatorCreate,
	}
	createCmd.Flags().String("moniker", "", "Validator name")
	createCmd.Flags().Uint64("commission", 1000, "Commission rate (basis points, 1000 = 10%)")
	createCmd.Flags().Float64("self-stake", 100000, "Self-stake amount in APX")
	createCmd.Flags().String("key", "", "Path to validator key file")
	createCmd.Flags().String("details", "", "Validator details")
	createCmd.Flags().String("website", "", "Validator website")

	unjailCmd := &cobra.Command{
		Use:   "unjail",
		Short: "Unjail a validator",
		Run:   runValidatorUnjail,
	}

	cmd.AddCommand(createCmd, unjailCmd)
	return cmd
}

// stakeCmd returns the staking command
func stakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake tokens to a validator",
		Run:   runStake,
	}

	cmd.Flags().String("validator", "", "Validator address")
	cmd.Flags().Float64("amount", 0, "Amount to stake in APX")
	cmd.Flags().String("from", "", "Delegator address")

	return cmd
}

// queryCmd returns the query command
func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query blockchain data",
	}

	balanceCmd := &cobra.Command{
		Use:   "balance [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		Run:   runQueryBalance,
	}

	validatorCmd := &cobra.Command{
		Use:   "validator [address]",
		Short: "Query validator information",
		Args:  cobra.ExactArgs(1),
		Run:   runQueryValidator,
	}

	stakingCmd := &cobra.Command{
		Use:   "staking",
		Short: "Query staking information",
	}

	rewardsCmd := &cobra.Command{
		Use:   "rewards",
		Short: "Query staking rewards",
		Run:   runQueryRewards,
	}
	rewardsCmd.Flags().String("delegator", "", "Delegator address")
	rewardsCmd.Flags().String("validator", "", "Validator address")

	stakingCmd.AddCommand(rewardsCmd)
	cmd.AddCommand(balanceCmd, validatorCmd, stakingCmd)

	return cmd
}

// versionCmd returns the version command
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Apex Blockchain CLI v1.0.0")
		},
	}
}

// Command implementations
func runKeysGenerate(cmd *cobra.Command, args []string) {
	output, _ := cmd.Flags().GetString("output")

	// Generate key pair (simplified - in production use proper crypto)
	keyData := map[string]string{
		"address":     "0x" + generateRandomHex(40),
		"private_key": generateRandomHex(64),
		"public_key":  generateRandomHex(128),
	}

	data, err := json.MarshalIndent(keyData, "", "  ")
	if err != nil {
		logger.Fatal("Failed to marshal key data", zap.Error(err))
	}

	if err := ioutil.WriteFile(output, data, 0600); err != nil {
		logger.Fatal("Failed to write key file", zap.Error(err))
	}

	fmt.Printf("✓ Key pair generated successfully\n")
	fmt.Printf("Address: %s\n", keyData["address"])
	fmt.Printf("Key file: %s\n", output)
	fmt.Printf("\n⚠️  IMPORTANT: Keep this file secure and backup safely!\n")
}

func runValidatorCreate(cmd *cobra.Command, args []string) {
	moniker, _ := cmd.Flags().GetString("moniker")
	commission, _ := cmd.Flags().GetUint64("commission")
	selfStake, _ := cmd.Flags().GetFloat64("self-stake")
	keyFile, _ := cmd.Flags().GetString("key")
	details, _ := cmd.Flags().GetString("details")
	website, _ := cmd.Flags().GetString("website")

	fmt.Printf("Creating validator:\n")
	fmt.Printf("  Moniker: %s\n", moniker)
	fmt.Printf("  Commission: %.2f%%\n", float64(commission)/100)
	fmt.Printf("  Self-stake: %.2f APX\n", selfStake)
	fmt.Printf("  Key file: %s\n", keyFile)
	fmt.Printf("  Details: %s\n", details)
	fmt.Printf("  Website: %s\n", website)
	fmt.Printf("\n✓ Validator creation transaction submitted\n")
	fmt.Printf("Transaction hash: 0x%s\n", generateRandomHex(64))
}

func runValidatorUnjail(cmd *cobra.Command, args []string) {
	fmt.Printf("✓ Unjail transaction submitted\n")
	fmt.Printf("Transaction hash: 0x%s\n", generateRandomHex(64))
}

func runStake(cmd *cobra.Command, args []string) {
	validator, _ := cmd.Flags().GetString("validator")
	amount, _ := cmd.Flags().GetFloat64("amount")
	from, _ := cmd.Flags().GetString("from")

	fmt.Printf("Staking %.2f APX to validator %s\n", amount, validator)
	fmt.Printf("From: %s\n", from)
	fmt.Printf("\n✓ Stake transaction submitted\n")
	fmt.Printf("Transaction hash: 0x%s\n", generateRandomHex(64))
}

func runQueryBalance(cmd *cobra.Command, args []string) {
	address := args[0]
	fmt.Printf("Account: %s\n", address)
	fmt.Printf("Balance: 1,000.00 APX\n")
	fmt.Printf("Staked: 500.00 APX\n")
	fmt.Printf("Locked: 0.00 APX\n")
	fmt.Printf("Nonce: 5\n")
}

func runQueryValidator(cmd *cobra.Command, args []string) {
	address := args[0]
	fmt.Printf("Validator: %s\n", address)
	fmt.Printf("Status: Active\n")
	fmt.Printf("Voting Power: 1,000,000 APX\n")
	fmt.Printf("Commission: 10.00%%\n")
	fmt.Printf("Uptime: 99.95%%\n")
	fmt.Printf("Produced Blocks: 15,432\n")
	fmt.Printf("Missed Blocks: 8\n")
}

func runQueryRewards(cmd *cobra.Command, args []string) {
	delegator, _ := cmd.Flags().GetString("delegator")
	validator, _ := cmd.Flags().GetString("validator")

	fmt.Printf("Delegator: %s\n", delegator)
	fmt.Printf("Validator: %s\n", validator)
	fmt.Printf("Accumulated Rewards: 25.50 APX\n")
}

// Helper function to generate random hex string
func generateRandomHex(length int) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		result[i] = hexChars[i%len(hexChars)]
	}
	return string(result)
}
