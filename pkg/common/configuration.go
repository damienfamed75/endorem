package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config stores all the data from the configuration file.
// This structure should match the "structure" of the file itself and
// how the variables are laid out in it.
type Config struct {
	Player struct {
		AddedHealth int   `json:"addedHealth"`
		JumpHeight  int32 `json:"jumpHeight"`
		MoveSpeed   int32 `json:"moveSpeed"`
	} `json:"person"`
	Enemy struct {
		AddedHealth         int     `json:"addedHealth"`
		MoveSpeedMultiplier float32 `json:"moveSpeedMultiplier"`
	} `json:"enemy"`
}

var (
	// GlobalConfiguration is used amongst many packages to read what is
	// said in the game's configuration file to do.
	GlobalConfiguration *Config
)

// LoadConfig loads the configuration file in from the designated path
// and then unmarshals it into a Config structure.
func LoadConfig() error {
	// Instantiate the global configuration variable to non-nil.
	GlobalConfiguration = &Config{}

	// Open the config file from the designated path.
	cfgFile, err := os.Open("config/game.json")
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	defer cfgFile.Close()

	// Read all the bytes from the config file.
	rawCfg, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	// Unmarshal the raw bytes to the config structure.
	if err := json.Unmarshal(rawCfg, GlobalConfiguration); err != nil {
		return fmt.Errorf("unmarshal config bytes: %w", err)
	}

	return nil
}
