package remote_manager

import (
	"time"

	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_manager/plugin_errors"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/log"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/routine"
)

func (r *RemotePluginRuntime) InitEnvironment() error {
	return nil
}

func (r *RemotePluginRuntime) Stopped() bool {
	return !r.alive
}

func (r *RemotePluginRuntime) Stop() {
	r.alive = false
	if r.conn == nil {
		return
	}
	r.conn.Close()
}

func (r *RemotePluginRuntime) Type() plugin_entities.PluginRuntimeType {
	return plugin_entities.PLUGIN_RUNTIME_TYPE_REMOTE
}

func (r *RemotePluginRuntime) StartPlugin() error {
	var exitError error

	identity, err := r.Identity()
	if err != nil {
		return err
	}

	// handle heartbeat
	routine.Submit(map[string]string{
		"module":    "remote_manager",
		"function":  "StartPlugin",
		"plugin_id": identity.String(),
	}, func() {
		r.lastActiveAt = time.Now()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if time.Since(r.lastActiveAt) > 60*time.Second {
					// kill this connection if it's not active for a long time
					r.conn.Close()
					exitError = plugin_errors.ErrPluginNotActive
					return
				}
			case <-r.shutdownChan:
				return
			}
		}
	})

	r.response.Async(func(data []byte) {
		plugin_entities.ParsePluginUniversalEvent(
			data,
			func(session_id string, data []byte) {
				r.callbacksLock.RLock()
				listeners := r.callbacks[session_id][:]
				r.callbacksLock.RUnlock()

				// handle session event
				for _, listener := range listeners {
					listener(data)
				}
			},
			func() {
				r.lastActiveAt = time.Now()
			},
			func(err string) {
				log.Error("plugin %s: %s", r.Configuration().Identity(), err)
			},
			func(message string) {
				log.Info("plugin %s: %s", r.Configuration().Identity(), message)
			},
		)
	})

	return exitError
}

func (r *RemotePluginRuntime) Wait() (<-chan bool, error) {
	return r.shutdownChan, nil
}

func (r *RemotePluginRuntime) Checksum() (string, error) {
	return r.checksum, nil
}