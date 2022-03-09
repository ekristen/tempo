package storage

import (
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	prometheus_config "github.com/prometheus/prometheus/config"
	"github.com/weaveworks/common/user"
)

// generateTenantRemoteWriteConfigs creates a copy of the remote write configurations with the
// X-Scope-OrgID header present for the given tenant. If the remote write config already contains
// this header it will be overwritten.
func generateTenantRemoteWriteConfigs(originalCfgs []*prometheus_config.RemoteWriteConfig, tenant string, logger log.Logger) []*prometheus_config.RemoteWriteConfig {
	var cloneCfgs []*prometheus_config.RemoteWriteConfig

	for _, originalCfg := range originalCfgs {
		cloneCfg := &prometheus_config.RemoteWriteConfig{}
		*cloneCfg = *originalCfg

		// Copy headers so we can modify them
		cloneCfg.Headers = copyMap(cloneCfg.Headers)

		// Ensure that no variation of the X-Scope-OrgId header can be added, which might trick authentication
		for k, v := range cloneCfg.Headers {
			if strings.EqualFold(user.OrgIDHeaderName, strings.TrimSpace(k)) {
				level.Warn(logger).Log("msg", "discarding X-Scope-OrgId header", "key", k, "value", v)
				delete(cloneCfg.Headers, k)
			}
		}

		// inject the X-Scope-OrgId header for multi-tenant metrics backends
		cloneCfg.Headers[user.OrgIDHeaderName] = tenant

		cloneCfgs = append(cloneCfgs, cloneCfg)
	}

	return cloneCfgs
}

// copyMap creates a new map containing all values from the given map.
func copyMap(m map[string]string) map[string]string {
	newMap := make(map[string]string, len(m))

	for k, v := range m {
		newMap[k] = v
	}

	return newMap
}
