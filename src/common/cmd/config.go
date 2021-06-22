package cmd

type DatabaseConfig struct {
	DatabaseDriver    string `envconfig:"database_driver"`
	DatabaseName      string `envconfig:"database_name"`
	DatabaseAddress   string `envconfig:"database_address"`
	DatabaseUser      string `envconfig:"database_user"`
	DatabasePassword  string `envconfig:"database_password"`
	DatabaseArguments string `envconfig:"database_arguments"`
}

type HTTPConfig struct {
	ServerPort int `envconfig:"http_port"`
}

type GRPCConfig struct {
	ServerPort int `envconfig:"grpc_port"`
}

type AMQPConfig struct {
	ServerUrl    string `envconfig:"amqp_server_url"`
	QueueName    string `envconfig:"amqp_queue_name"`
	ExchangeType string `envconfig:"amqp_exchange_type"`
}
