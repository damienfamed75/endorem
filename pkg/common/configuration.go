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
	Game struct {
		Volume int `json:"volume"`
		Screen struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"screen"`
	} `json:"game"`

	Player struct {
		AddedHealth     int   `json:"addedHealth"`
		JumpHeight      int32 `json:"jumpHeight"`
		MoveSpeed       int32 `json:"moveSpeed"`
		AttackTimer     int   `json:"attackTimer"`     // milliseconds
		InvincibleTimer int   `json:"invincibleTimer"` // milliseconds
	} `json:"player"`

	Enemy struct {
		AddedHealth         int     `json:"addedHealth"`
		MoveSpeedMultiplier float32 `json:"moveSpeedMultiplier"`
		AttackTimer         int     `json:"attackTimer"` // milliseconds
	} `json:"enemy"`
}

var (
	// GlobalConfig is used amongst many packages to read what is
	// said in the game's configuration file to do.
	GlobalConfig *Config
)

// LoadConfig loads the configuration file in from the designated path
// and then unmarshals it into a Config structure.
func LoadConfig() error {
	// Instantiate the global configuration variable to non-nil.
	GlobalConfig = &Config{}

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
	if err := json.Unmarshal(rawCfg, GlobalConfig); err != nil {
		return fmt.Errorf("unmarshal config bytes: %w", err)
	}

	return nil
}

func (c *Config) ScreenWidth() int {
	return c.Game.Screen.Width
}

func (c *Config) ScreenHeight() int {
	return c.Game.Screen.Height
}
