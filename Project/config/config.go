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
	Password: "Aa509224765_",
	DBName:   "Tiktok",
	Port:     "3306",
	CharSet:  "utf8mb4",
}

type snowflakeConfig struct {
	WorkerId     int64 //工作机器ID
	DateCenterId int64 //数据中心ID
}

var SnowFlakeCfg = snowflakeConfig{
	WorkerId:     0,
	DateCenterId: 0,
}
