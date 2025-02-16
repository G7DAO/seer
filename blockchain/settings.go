package blockchain

import (
	"fmt"
	"os"
)

var (
	BlockchainURLs map[string]string
)

// CheckVariablesForBlockchains checks required environment variables but only for the specified chain.
func CheckVariablesForBlockchains(chain string) error {

	// 3. Based on the chain, require only the relevant environment variables
	switch chain {
	case "ethereum":
		if err := checkEnv("MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "sepolia":
		if err := checkEnv("MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "polygon":
		if err := checkEnv("MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "arbitrum_one":
		if err := checkEnv("MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "arbitrum_sepolia":
		if err := checkEnv("MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "game7_orbit_arbitrum_sepolia":
		if err := checkEnv("MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "game7_testnet":
		if err := checkEnv("MOONSTREAM_NODE_GAME7_TESTNET_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "game7":
		if err := checkEnv("MOONSTREAM_NODE_GAME7_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "xai":
		if err := checkEnv("MOONSTREAM_NODE_XAI_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "xai_sepolia":
		if err := checkEnv("MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "mantle":
		if err := checkEnv("MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "mantle_sepolia":
		if err := checkEnv("MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "imx_zkevm":
		if err := checkEnv("MOONSTREAM_NODE_IMX_ZKEVM_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "imx_zkevm_sepolia":
		if err := checkEnv("MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "b3":
		if err := checkEnv("MOONSTREAM_NODE_B3_A_EXTERNAL_URI"); err != nil {
			return err
		}
	case "b3_sepolia":
		if err := checkEnv("MOONSTREAM_NODE_B3_SEPOLIA_A_EXTERNAL_URI"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown chain specified: %s", chain)
	}

	// 4. Populate BlockchainURLs with only the relevant entry for `chain`
	BlockchainURLs = make(map[string]string)

	// Retrieve all env vars even if we only check for one. This is minimal,
	// but prevents major changes to how BlockchainURLs are built:

	MOONSTREAM_NODE_ETHEREUM := os.Getenv("MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI")
	MOONSTREAM_NODE_SEPOLIA := os.Getenv("MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI")
	MOONSTREAM_NODE_POLYGON := os.Getenv("MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI")
	MOONSTREAM_NODE_ARBITRUM_ONE := os.Getenv("MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI")
	MOONSTREAM_NODE_ARBITRUM_SEPOLIA := os.Getenv("MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI")
	MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA := os.Getenv("MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI")
	MOONSTREAM_NODE_GAME7_TESTNET := os.Getenv("MOONSTREAM_NODE_GAME7_TESTNET_A_EXTERNAL_URI")
	MOONSTREAM_NODE_GAME7 := os.Getenv("MOONSTREAM_NODE_GAME7_A_EXTERNAL_URI")
	MOONSTREAM_NODE_XAI := os.Getenv("MOONSTREAM_NODE_XAI_A_EXTERNAL_URI")
	MOONSTREAM_NODE_XAI_SEPOLIA := os.Getenv("MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI")
	MOONSTREAM_NODE_MANTLE := os.Getenv("MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI")
	MOONSTREAM_NODE_MANTLE_SEPOLIA := os.Getenv("MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI")
	MOONSTREAM_NODE_IMX_ZKEVM := os.Getenv("MOONSTREAM_NODE_IMX_ZKEVM_A_EXTERNAL_URI")
	MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA := os.Getenv("MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA_A_EXTERNAL_URI")
	MOONSTREAM_NODE_B3 := os.Getenv("MOONSTREAM_NODE_B3_A_EXTERNAL_URI")
	MOONSTREAM_NODE_B3_SEPOLIA := os.Getenv("MOONSTREAM_NODE_B3_SEPOLIA_A_EXTERNAL_URI")

	BlockchainURLs["ethereum"] = MOONSTREAM_NODE_ETHEREUM
	BlockchainURLs["sepolia"] = MOONSTREAM_NODE_SEPOLIA
	BlockchainURLs["polygon"] = MOONSTREAM_NODE_POLYGON
	BlockchainURLs["arbitrum_one"] = MOONSTREAM_NODE_ARBITRUM_ONE
	BlockchainURLs["arbitrum_sepolia"] = MOONSTREAM_NODE_ARBITRUM_SEPOLIA
	BlockchainURLs["game7_orbit_arbitrum_sepolia"] = MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA
	BlockchainURLs["game7_testnet"] = MOONSTREAM_NODE_GAME7_TESTNET
	BlockchainURLs["game7"] = MOONSTREAM_NODE_GAME7
	BlockchainURLs["xai"] = MOONSTREAM_NODE_XAI
	BlockchainURLs["xai_sepolia"] = MOONSTREAM_NODE_XAI_SEPOLIA
	BlockchainURLs["mantle"] = MOONSTREAM_NODE_MANTLE
	BlockchainURLs["mantle_sepolia"] = MOONSTREAM_NODE_MANTLE_SEPOLIA
	BlockchainURLs["imx_zkevm"] = MOONSTREAM_NODE_IMX_ZKEVM
	BlockchainURLs["imx_zkevm_sepolia"] = MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA
	BlockchainURLs["b3"] = MOONSTREAM_NODE_B3
	BlockchainURLs["b3_sepolia"] = MOONSTREAM_NODE_B3_SEPOLIA

	return nil
}

// checkEnv is a tiny helper verifying that an environment variable is non-empty
func checkEnv(envVar string) error {
	if os.Getenv(envVar) == "" {
		return fmt.Errorf("%s environment variable is required", envVar)
	}
	return nil
}
