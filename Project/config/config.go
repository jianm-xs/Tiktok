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
	Url:      "81.70.17.190",
	UserName: "root",
	Password: "root",
	DBName:   "gorm_test",
	Port:     "3306",
	CharSet:  "utf8mb4",
}
