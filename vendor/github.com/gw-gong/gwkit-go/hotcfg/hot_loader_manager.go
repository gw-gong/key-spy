package hotcfg

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gw-gong/gwkit-go/util"
)

type HotLoaderManager interface {
	RegisterHotLoader(hotLoader HotLoader) error
	Watch() error
}

func NewHotLoaderManager() HotLoaderManager {
	return &hotLoaderManager{
		hotLoaders: make([]HotLoader, 0),
	}
}

type hotLoaderManager struct {
	mux        sync.Mutex
	isWatching int32

	hotLoaders []HotLoader
}

func (hlm *hotLoaderManager) RegisterHotLoader(hotLoader HotLoader) error {
	hlm.hotLoaders = append(hlm.hotLoaders, hotLoader)
	return nil
}

// After all the registrations are completed, start the watch.
func (hlm *hotLoaderManager) Watch() error {
	if atomic.LoadInt32(&hlm.isWatching) == 1 {
		return fmt.Errorf("already watching, don't call Watch again")
	}
	atomic.StoreInt32(&hlm.isWatching, 1)

	var errors []error

	for _, hotLoader := range hlm.hotLoaders {
		if localConfig := hotLoader.AsLocalConfig(); localConfig != nil {
			localConfig.WatchLocalConfig(hotLoader.LoadConfig)
		} else if consulConfig := hotLoader.AsConsulConfig(); consulConfig != nil {
			go util.WithRecover(func() {
				func(consulConfig ConsulConfig, hotLoader HotLoader) {
					hlm.watchConsulConfig(consulConfig, hotLoader.LoadConfig)
				}(consulConfig, hotLoader)
			})
		} else {
			errors = append(errors, fmt.Errorf("hot loader config struct error: %v", hotLoader.GetBaseConfig()))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("hot loader config struct error: %v", errors)
	}
	return nil
}

func (hlm *hotLoaderManager) watchConsulConfig(consulConfig ConsulConfig, loadConfig func()) {
	ticker := time.NewTicker(time.Duration(consulConfig.GetConsulReloadTime()) * time.Second)
	defer ticker.Stop()

	lastConfigHash := consulConfig.CalculateConsulConfigHash()

	for range ticker.C {
		if err := consulConfig.ReadConsulConfig(); err != nil {
			continue
		}

		currentConfigHash := consulConfig.CalculateConsulConfigHash()
		if currentConfigHash != lastConfigHash {
			lastConfigHash = currentConfigHash
			hlm.mux.Lock()
			loadConfig()
			hlm.mux.Unlock()
		}
	}
}
