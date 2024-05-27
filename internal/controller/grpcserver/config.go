package grpcserver

type Config struct {
	BootstrapAddress  string
	SecretJWTKey      string
	JWTDurationMinute int
}
