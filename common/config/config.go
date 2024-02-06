package config

type Config struct {
	Server      Server `yaml:"server"`
	Postgres    Postgres
	Redis       Redis
	LogFile     LogFile     `yaml:"logfile"`
	UrlGenerate UrlGenerate `yaml:"url"`
}

type LogFile struct {
	Path string `yaml:"path"`
}

type Server struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	GrpcPort string `yaml:"grpc_port"`
}

type UrlGenerate struct {
	Alphabet string `yaml:"alphabet"`
	Length   int    `yaml:"length"`
}

type Postgres struct {
	Name     string
	User     string
	Port     int
	Password string
	Host     string
}

type Redis struct {
	Address string
	DB      int
}
