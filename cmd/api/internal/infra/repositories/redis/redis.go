package redis

import "github.com/redis/rueidis"

func NewClient(ips []string, password string) (rueidis.Client, error) {
	return rueidis.NewClient(rueidis.ClientOption{
		InitAddress: ips,
		Password:    password,
	})
}
