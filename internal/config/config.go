package config

type Config struct {
	ServerHost   string `json:"server_host"`
	ServerPort   string `json:"server_port"`
	DBHost       string `json:"db_host"`
	DBPort       string `json:"db_port"`
	DBUsername   string `json:"db_username"`
	DBPassword   string `json:"db_password"`
	DatabaseName string `json:"database_name"`
}
