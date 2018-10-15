package envload

import (
	"errors"
	"fmt"
	"os"
	"plotbot-server/logging"
	"strconv"

	"github.com/fatih/color"
	"github.com/subosito/gotenv"
)

//VirtualPatientConfig is used to describe the configuration state of the Go server
type VirtualPatientConfig struct {
	Production bool
	BindPort   string
	BindIP     string
}

func loadBoolEnv(varName string) bool {
	if os.Getenv(varName) == "" {
		color.Yellow(fmt.Sprintf("Missing %s value in .env file, automatically setting to false.\nSet a boolean value for %s in your .env file to disable this warning.", varName, varName))
		return false
	}
	boolEnv, err := strconv.ParseBool(os.Getenv(varName))
	if err != nil {
		color.Yellow(fmt.Sprintf("%s value must be a valid bool (true or false)\n Automatically setting to false.", varName))
		return false
	}
	return boolEnv

}

func loadStringEnv(varName string) string {
	if os.Getenv(varName) == "" {
		logging.Fatal(fmt.Sprintf("Missing %s value in .env file.", varName), errors.New("No such string var defined"))
	}
	return os.Getenv(varName)
}

//LoadEnv uses gotenv to pull in the env files specified in the args, and returns a config struct
func LoadEnv(filenames ...string) (VirtualPatientConfig, error) {
	//Load Env
	err := gotenv.Load(filenames...)
	if err != nil {
		logging.Error("Error loading .env file: ", err)
		return VirtualPatientConfig{}, err
	}
	//Setup global env variables
	config := VirtualPatientConfig{
		Production: loadBoolEnv("PRODUCTION"),
		BindPort:   loadStringEnv("BIND_PORT"),
		BindIP:     loadStringEnv("BIND_IP"),
		/*BindURL:    loadStringEnv("BIND_URL"),
		JWTSecret:  loadStringEnv("JWT_SECRET"),
		MongoURI:   loadStringEnv("MONGO_URI"),*/
	}
	return config, nil
}
