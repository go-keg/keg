package config

import "strings"

type Log struct {
	Dir          string
	Level        string
	MaxAge       int
	RotationTime int
}

type Database struct {
	Driver          string
	Dsn             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime string // time.Duration
	ConnMaxIdleTime string // time.Duration
}

type Email struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type Redis struct {
	Addr     string
	DB       string
	Password string
	Prefix   string
}

type ElasticSearch struct {
	Hosts    string
	UserName string
	Password string
	Sniff    string
}

type Kafka struct {
	Addrs    string `yaml:"addrs"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (r Kafka) GetAddr() []string {
	return strings.Split(r.Addrs, ",")
}

type KafkaConsumerGroup struct {
	GroupID string   `yaml:"groupId"`
	Topics  []string `yaml:"topics"`
}
