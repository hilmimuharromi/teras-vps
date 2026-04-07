package config

import (
	"os"
	"strconv"
)

// Config holds all configuration values
type Config struct {
	// Server
	Port string

	// Frontend
	FrontendURL string

	// Database
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	// Redis
	RedisHost string
	RedisPort string
	RedisPass string

	// JWT
	JWTSecret     string
	JWTExpiration int // in hours

	// Proxmox
	ProxmoxHost     string
	ProxmoxUser     string
	ProxmoxPassword string
	ProxmoxNode     string

	// Business Logic
	MaxVMsPerUser     int
	SuspendDays       int // days after due date to suspend
	DeleteDays        int // days after due date to delete
	MaxBackupsPerVM   int

	// Currency
	Currency string
}

// Load configuration from environment variables
func Load() *Config {
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION", "24"))
	maxVMsPerUser, _ := strconv.Atoi(getEnv("MAX_VMS_PER_USER", "5"))
	suspendDays, _ := strconv.Atoi(getEnv("SUSPEND_DAYS", "7"))
	deleteDays, _ := strconv.Atoi(getEnv("DELETE_DAYS", "14"))
	maxBackupsPerVM, _ := strconv.Atoi(getEnv("MAX_BACKUPS_PER_VM", "3"))

	return &Config{
		Port:       getEnv("PORT", "3000"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:4321"),

		// Database
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "teras_vps"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", "teras_vps"),

		// Redis
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
		RedisPass: getEnv("REDIS_PASSWORD", ""),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiration: jwtExpiration,

		// Proxmox
		ProxmoxHost:     getEnv("PROXMOX_HOST", ""),
		ProxmoxUser:     getEnv("PROXMOX_USER", "root@pam"),
		ProxmoxPassword: getEnv("PROXMOX_PASSWORD", ""),
		ProxmoxNode:     getEnv("PROXMOX_NODE", "proxmox"),

		// Business Logic
		MaxVMsPerUser:   maxVMsPerUser,
		SuspendDays:     suspendDays,
		DeleteDays:      deleteDays,
		MaxBackupsPerVM: maxBackupsPerVM,

		// Currency
		Currency: getEnv("CURRENCY", "IDR"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
