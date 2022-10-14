package configs

type Config struct {
	AppHost            string `env:"HOST" envDefault:"127.0.0.1"`
	AppPort            string `env:"PORT" envDefault:"9090"`
	DatabaseDsn        string `env:"DATABASE_DSN" envDefault:"host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"`
	KafkaNetwork       string `env:"KAFKA_NETWORK" envDefault:"tcp"`
	KafkaHost          string `env:"KAFKA_HOST" envDefault:"localhost"`
	KafkaPort          string `env:"KAFKA_PORT" envDefault:"9092"`
	KafkaTopic         string `env:"KAFKA_TOPIC" envDefault:"company-api"`
	KafkaGroupId       string `env:"KAFKA_GROUP" envDefault:"compamy_api_group"`
	KafkaWriteDeadline int    `env:"KAFKA_WRITE_DEADLINE" envDefault:"8"`
	AccessTokenTTL     string `env:"ACCESS_TOKEN_TTL" envDefault:"120s"`
}
