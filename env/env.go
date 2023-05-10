package env

type AppEnv int8

const (
	EnvDev AppEnv = iota + 1
	EnvTest
	EnvPre
	EnvProd
)

var EnvMap = map[AppEnv]string{
	EnvDev:  "dev",
	EnvTest: "test",
	EnvPre:  "pre",
	EnvProd: "dp",
}

var EnvProdMap = map[AppEnv]string{
	EnvDev:  "dev",
	EnvTest: "test",
	EnvPre:  "pre",
	EnvProd: "prod",
}
