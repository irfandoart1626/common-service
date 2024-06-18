package util

import (
	"cicd-gitlab-ee.telkomsel.co.id/phincon-go/common-service/log"
	"crypto/rand"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var vips *viper.Viper = &viper.Viper{}

func GetEnv(key string) string {
	if value := vips.GetString(key); value != "" {
		return value
	}
	panic(fmt.Errorf("config %s not found", key))
}

func init() {
	configname := os.Getenv("GO_PROFILE")

	if strings.EqualFold(configname, "") {
		configname = "default"
	}

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}

	// Search for the "resources" folder recursively starting from the current directory
	resourcesPath := findResourcesFolder(currentDir)

	viper.AddConfigPath(resourcesPath)

	viper.SetConfigType("properties")
	viper.SetConfigName(configname)

	vips = viper.GetViper()

	vips.WatchConfig()
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error %w", err))
	}

	vips = viper.GetViper()

	InitConfig(vips)
}

func findResourcesFolder(dir string) string {
	for dir != "" {

		resourcesPath := filepath.Join(dir, "config")

		_, err := os.Stat(resourcesPath)
		if err == nil {
			return resourcesPath
		}

		// this will substring the last path
		parentDir := filepath.Dir(dir)

		// if root directory already reaced
		if parentDir == dir {
			break
		}
		dir = parentDir
	}
	return "" // "resources" folder not found
}

func InitConfig(vipers *viper.Viper) {
	zerolog.SetGlobalLevel(log.GetLevel(GetEnv("LOG_LEVEL")))
	zerolog.TimestampFieldName = ""
	devDebugMode := GetEnv("DEV_DEBUG_MODE")
	log.SetupLogger(devDebugMode == "true")
}

func ValidateMSISDN(msisdn string) bool {
	// convert to number
	if _, err := strconv.Atoi(msisdn); err != nil {
		return false
	}

	// check msisdn format and length
	if !strings.HasPrefix(msisdn, "62") || len(msisdn) < 11 ||
		len(msisdn) > 13 {
		return false
	}

	return true
}

// Generate NOIS Transaction ID
func GenerateTransactionID(trxID string, msisdn string, internalCode string, apiID string) string {
	if trxID != "" { // not empty string && not nil
		return trxID
	}

	varInternalCode := "0"
	varMSISDNSuffix := ""
	varAppID := "N001"
	varTimestamp := strings.Replace(time.Now().Format("060102150405.000"), ".", "", -1)

	if (internalCode) != "" {
		varInternalCode = internalCode[len(internalCode)-1:]
	}

	if msisdn != "" { // not empty string && not nil
		lgt := len(msisdn)
		if lgt > 4 {
			varMSISDNSuffix = msisdn[lgt-5:]
		} else {
			for i := 0; i < 5-lgt; i++ {
				varMSISDNSuffix += "0"
			}
			varMSISDNSuffix += msisdn
		}
	} else {
		varMSISDNSuffix = generateRandomNumber(5)
	}

	if (apiID) != "" {
		varAppID = apiID
	}

	return varAppID + varTimestamp + varMSISDNSuffix + varInternalCode
}

func generateRandomNumber(maxDigits uint32) string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%0*d", maxDigits, bi)
}
