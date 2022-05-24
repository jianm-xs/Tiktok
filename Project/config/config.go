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

var MysqlCfg = mysqlConfig{
	Url:      "127.0.0.1",
	UserName: "root",
	Password: "root",
	DBName:   "Tiktok",
	Port:     "3306",
	CharSet:  "utf8mb4",
}
