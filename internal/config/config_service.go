package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env string

	DatabaseDSN string

	RedisAddr            string
	RedisPassword        string
	RedisDatabase        int
	RedisDialTimeout     time.Duration
	RedisReadTimeout     time.Duration
	RedisWriteTimeout    time.Duration
	RedisMinRetryBackoff time.Duration
	RedisMaxRetryBackoff time.Duration
	RedisMaxRetries      int

	ServerAddress         string
	ServerReadTimeout     time.Duration
	ServerWriteTimeout    time.Duration
	ServerGracefulTimeout time.Duration
	ServerMaxHeaderBytes  int
	ServerChildProcesses  int

	CORSEnabled          bool
	CORSAllowOrigin      string
	CORSAllowMethods     string
	CORSAllowHeaders     string
	CORSAllowCredentials bool
	CORSExposeHeaders    string
	CORSMaxAge           time.Duration

	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string

	JWTSecret            string
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration
}

var (
	cfg  *Config
	once sync.Once
)

func Init(env string) *Config {
	once.Do(func() {
		if env == "" {
			env = "dev"
		}
		cfg = &Config{Env: env}
		envFile := filepath.Join(".env." + env)

		if err := godotenv.Load(envFile); err != nil {
			slog.Warn("Failed to load .env file, fallback to system env", "file", envFile, "err", err)
		}

		cfg.DatabaseDSN = getEnv("DATABASE_DSN", "postgres://postgres:password@postgres:5432/brickbang_dev?sslmode=disable")

		cfg.RedisAddr = getEnv("REDIS_ADDR", "127.0.0.1:6379")
		cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
		cfg.RedisDatabase = parseInt(getEnv("REDIS_DATABASE", "0"))
		cfg.RedisDialTimeout = parseDuration(getEnv("REDIS_DIAL_TIMEOUT", "5s"))
		cfg.RedisReadTimeout = parseDuration(getEnv("REDIS_READ_TIMEOUT", "3s"))
		cfg.RedisWriteTimeout = parseDuration(getEnv("REDIS_WRITE_TIMEOUT", "3s"))
		cfg.RedisMinRetryBackoff = parseDuration(getEnv("REDIS_MIN_RETRY_BACKOFF", "100ms"))
		cfg.RedisMaxRetryBackoff = parseDuration(getEnv("REDIS_MAX_RETRY_BACKOFF", "2s"))
		cfg.RedisMaxRetries = parseInt(getEnv("REDIS_MAX_RETRIES", "10"))

		cfg.ServerAddress = getEnv("SERVER_ADDRESS", "0.0.0.0:8080")
		cfg.ServerReadTimeout = parseDuration(getEnv("SERVER_READ_TIMEOUT", "3s"))
		cfg.ServerWriteTimeout = parseDuration(getEnv("SERVER_WRITE_TIMEOUT", "3s"))
		cfg.ServerGracefulTimeout = parseDuration(getEnv("SERVER_GRACEFUL_TIMEOUT", "5s"))
		cfg.ServerMaxHeaderBytes = parseInt(getEnv("SERVER_MAX_HEADER_BYTES", "1048576"))
		cfg.ServerChildProcesses = parseInt(getEnv("SERVER_CHILD_PROCESSES", strconv.Itoa(runtime.NumCPU())))

		cfg.CORSEnabled = parseBool(getEnv("CORS_ENABLED", "true"))
		cfg.CORSAllowOrigin = getEnv("CORS_ALLOW_ORIGIN", "*")
		cfg.CORSAllowMethods = getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		cfg.CORSAllowHeaders = getEnv("CORS_ALLOW_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token")
		cfg.CORSAllowCredentials = parseBool(getEnv("CORS_ALLOW_CREDENTIALS", "false"))
		cfg.CORSExposeHeaders = getEnv("CORS_EXPOSE_HEADERS", "*")
		cfg.CORSMaxAge = parseDuration(getEnv("CORS_MAX_AGE", "10m"))

		cfg.TLSEnabled = parseBool(getEnv("TLS_ENABLED", "false"))
		cfg.TLSCertFile = getEnv("TLS_CERT_FILE", "./cert/server.crt")
		cfg.TLSKeyFile = getEnv("TLS_KEY_FILE", "./cert/server.key")

		cfg.JWTSecret = getEnv("JWT_SECRET", "devsecret")
		cfg.JWTAccessExpiration = parseDuration(getEnv("JWT_ACCESS_EXPIRATION", "15m"))
		cfg.JWTRefreshExpiration = parseDuration(getEnv("JWT_REFRESH_EXPIRATION", "24h"))
	})

	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Init("") // fallback: dev
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseBool(s string) bool {
	return s == "true" || s == "1"
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
