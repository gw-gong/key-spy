package hotcfg

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/gw-gong/gwkit-go/log"
	"github.com/spf13/viper"
)

// Calculate configuration hash for change detection
func CalculateConfigHash(v *viper.Viper) string {
	// Get all configuration settings
	settings := v.AllSettings()

	// Convert configuration to JSON string
	configBytes, err := json.Marshal(settings)
	if err != nil {
		log.Error("Failed to serialize configuration: %v", log.Err(err))
		return ""
	}

	// Calculate MD5 hash
	hash := md5.Sum(configBytes)
	return fmt.Sprintf("%x", hash)
}
