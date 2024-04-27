package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName string = "novel_auth"

func LoadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		panic(err)
	}
}
