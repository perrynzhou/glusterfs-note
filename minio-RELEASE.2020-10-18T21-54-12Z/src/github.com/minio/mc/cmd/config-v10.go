/*
 * MinIO Client (C) 2020 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"sync"

	"github.com/minio/mc/pkg/probe"
	"github.com/minio/minio/pkg/quick"
)

const (
	defaultAccessKey = "YOUR-ACCESS-KEY-HERE"
	defaultSecretKey = "YOUR-SECRET-KEY-HERE"
)

var (
	// set once during first load.
	cacheCfgV10 *configV10
	// All access to mc config file should be synchronized.
	cfgMutex = &sync.RWMutex{}
)

// aliasConfig configuration of an alias.
type aliasConfigV10 struct {
	URL          string `json:"url"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	SessionToken string `json:"sessionToken,omitempty"`
	API          string `json:"api"`
	Path         string `json:"path"`
}

// configV10 config version.
type configV10 struct {
	Version string                    `json:"version"`
	Aliases map[string]aliasConfigV10 `json:"aliases"`
}

// newConfigV10 - new config version.
func newConfigV10() *configV10 {
	cfg := new(configV10)
	cfg.Version = globalMCConfigVersion
	cfg.Aliases = make(map[string]aliasConfigV10)
	return cfg
}

// SetAlias sets host config if not empty.
func (c *configV10) setAlias(alias string, cfg aliasConfigV10) {
	if _, ok := c.Aliases[alias]; !ok {
		c.Aliases[alias] = cfg
	}
}

// load default values for missing entries.
func (c *configV10) loadDefaults() {
	// MinIO server running locally.
	c.setAlias("local", aliasConfigV10{
		URL:       "http://localhost:9000",
		AccessKey: "",
		SecretKey: "",
		API:       "S3v4",
		Path:      "auto",
	})

	// Amazon S3 cloud storage service.
	c.setAlias("s3", aliasConfigV10{
		URL:       "https://s3.amazonaws.com",
		AccessKey: defaultAccessKey,
		SecretKey: defaultSecretKey,
		API:       "S3v4",
		Path:      "dns",
	})

	// Google cloud storage service.
	c.setAlias("gcs", aliasConfigV10{
		URL:       "https://storage.googleapis.com",
		AccessKey: defaultAccessKey,
		SecretKey: defaultSecretKey,
		API:       "S3v2",
		Path:      "dns",
	})

	// MinIO anonymous server for demo.
	c.setAlias("play", aliasConfigV10{
		URL:       "https://play.min.io",
		AccessKey: "Q3AM3UQ867SPQQA43P2F",
		SecretKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		API:       "S3v4",
		Path:      "auto",
	})
}

// loadConfigV10 - loads a new config.
func loadConfigV10() (*configV10, *probe.Error) {
	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	// If already cached, return the cached value.
	if cacheCfgV10 != nil {
		return cacheCfgV10, nil
	}

	if !isMcConfigExists() {
		return nil, errInvalidArgument().Trace()
	}

	// Initialize a new config loader.
	qc, e := quick.NewConfig(newConfigV10(), nil)
	if e != nil {
		return nil, probe.NewError(e)
	}

	// Load config at configPath, fails if config is not
	// accessible, malformed or version missing.
	if e = qc.Load(mustGetMcConfigPath()); e != nil {
		return nil, probe.NewError(e)
	}

	cfgV10 := qc.Data().(*configV10)

	// Cache config.
	cacheCfgV10 = cfgV10

	// Success.
	return cfgV10, nil
}

// saveConfigV10 - saves an updated config.
func saveConfigV10(cfgV10 *configV10) *probe.Error {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	qs, e := quick.NewConfig(cfgV10, nil)
	if e != nil {
		return probe.NewError(e)
	}

	// update the cache.
	cacheCfgV10 = cfgV10

	e = qs.Save(mustGetMcConfigPath())
	if e != nil {
		return probe.NewError(e).Trace(mustGetMcConfigPath())
	}
	return nil
}
