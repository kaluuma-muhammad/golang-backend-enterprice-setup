package ratelimiter

import "time"

type Limit struct {
	Requests int           // number of allowed requests
	Window   time.Duration // time window (e.g. 1 minute)
	Burst    int           // allowed burst (for token bucket)
}

type Config struct {
	Default         Limit
	Routes          map[string]Limit
	Public          Limit
	Protected       Limit
	Verified        Limit
	CleanupInterval time.Duration
	MaxIdleTime     time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Default: Limit{
			Requests: 100,
			Window:   time.Minute,
			Burst:    20,
		},

		Public: Limit{
			Requests: 20,
			Window:   time.Minute,
			Burst:    5,
		},

		Protected: Limit{
			Requests: 100,
			Window:   time.Minute,
			Burst:    20,
		},

		Verified: Limit{
			Requests: 300,
			Window:   time.Minute,
			Burst:    50,
		},

		Routes: map[string]Limit{
			"POST:/api/v1/auth/login": {
				Requests: 5,
				Window:   time.Minute,
				Burst:    2,
			},
			"POST:/api/v1/auth/register": {
				Requests: 3,
				Window:   time.Minute,
				Burst:    1,
			},
			"POST:/api/v1/auth/forgot-password": {
				Requests: 3,
				Window:   time.Minute,
				Burst:    1,
			},
		},

		CleanupInterval: time.Minute * 5,
		MaxIdleTime:     time.Minute * 10,
	}
}
