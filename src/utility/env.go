package utility

import "os"

func GetEnvPath(path string) string {
	envPath := "./"
	_, res := os.LookupEnv("HEROKU")
	if res {
		envPath = "/app/src/"
	}
	return envPath + path
}

func GetEnv(envName string) string {
	env, res := os.LookupEnv(envName)
	if !res {
		env = ""
	}
	return env
}

func GetPortString() string {
	port, res := os.LookupEnv("PORT")
	if !res {
		port = "8080"
	}
	return ":" + port
}
