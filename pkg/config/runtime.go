package config

import (
	"fmt"
)

// RuntimeConfig is the set of configurable options that are read when the program starts.
type RuntimeConfig struct {
	Host               string `env:"HOST" envDefault:"0.0.0.0"`
	Port               int    `env:"PORT" envDefault:"8080"`
	Debug              bool   `env:"DEBUG" envDefault:"false"`
	LogFormat          string `env:"LOG_FORMAT" envDefault:"json"`
	DatabaseFile       string `env:"DATABASE_FILE" envDefault:"pagebin.data"`
	BlobBackend        string `env:"BLOB_BACKEND" envDefault:"localfs"`
	BlobLocalFSRootDir string `env:"BLOB_LOCAL_FS_ROOT_DIR" envDefault:"pagebin-content"`
	ContentCacheSize   int    `env:"CONTENT_CACHE_SIZE" envDefault:"100"`
}

// ServerAddr returns the concatenated hostname with port.
func (cfg RuntimeConfig) ServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
