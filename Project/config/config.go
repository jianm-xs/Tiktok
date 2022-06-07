// package config 项目相关配置
// 创建人：吴润泽
// 创建时间：2022-05-23

package config

type mysqlConfig struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
	Port     string `yaml:"port"`
	CharSet  string `yaml:"charset"`
}

type RedisConfig struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

var MysqlCfg = mysqlConfig{
	Url:      "127.0.0.1",
	UserName: "root",
	Password: "root",
	DBName:   "Tiktok",
	Port:     "3306",
	CharSet:  "utf8mb4",
}

var RedisCfg = RedisConfig{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

type snowflakeConfig struct {
	WorkerID     int64 //工作机器ID
	DateCenterID int64 //数据中心ID
}

// 本机器的雪花 ID 配置

var SnowFlakeCfg = snowflakeConfig{
	WorkerID:     0,
	DateCenterID: 0,
}

// ServerHost 服务器地址
var ServerHost = "192.168.31.213"

// ServerPort 服务端口
var ServerPort = ":1080"
