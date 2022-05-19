package config

type mysqlConfig struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
	Port     string `yaml:"port"`
	CharSet  string `yaml:"charset"`
}

var MysqlCfg = mysqlConfig{
	Url:      "127.0.0.1",
	UserName: "root",
	Password: "root",
	DBName:   "Tiktok",
	Port:     "3306",
	CharSet:  "utf8mb4",
}
