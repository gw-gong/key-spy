package setting

type Env string

const (
	ENV_TEST    Env = "test" // test or dev
	ENV_STAGING Env = "staging"
	ENV_LIVE    Env = "live" // live or prod
)

var globalEnv = ENV_TEST

func SetEnv(env Env) {
	globalEnv = env
}

func GetEnv() Env {
	return globalEnv
}
