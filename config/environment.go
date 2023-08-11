package config

import (
	"os"
	"strings"
)

// Default configuration
var Environment = map[string]string{
	"app_port": "8999",

	// Database
	"db_host":     "pgsql-server",
	"db_port":     "5432",
	"db_user":     "postgres",
	"db_password": "secr3tPWD",
	"db_name":     "web_db",

	// Database Pool
	"db_pool_max_open_conns": "25",
	"db_pool_max_idle_conns": "25",
	"db_pool_max_lifetime":   "5",

	// Redis
	"redis_addr":     "localhost:6379",
	"redis_password": "my_master_password",
	"redis_db":       "0",

	// others
	"cache_article_exp_time": "1", // article cache expired in minutes
}

// Get perform to get configuration
// Default taken from Environment Map with lowercase key
// if uppercase key from env is not exist
func Get(key string) string {
	cfg := os.Getenv(strings.ToUpper(key))
	if cfg == "" {
		if val, ok := Environment[strings.ToLower(key)]; ok {
			cfg = val // Default config taken from Environment
		}
	}

	return cfg
}
