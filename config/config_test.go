package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input string = `{
    "business_logic":{
        "consumer":{
            "enabled": true,
            "api_keys":"consumer_key",
            "address":"api.thingspeak.com",
            "schema":"https",
            "uri":"update"
        }
    },
    "weather_provider":{
        "enabled":true
    },
    "forecast_provider":{
        "enabled":true
    } 
}`

func TestConfig(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), "testCfgFile.json")
	assert.NoError(t, err)
	fname := file.Name()
	defer func() {
		os.Remove(fname)
	}()

	_, err = file.WriteString(input)
	assert.NoError(t, err)
	file.Close()

	cfg, err := ReadFromFile(fname)
	assert.NoError(t, err)

	assert.Equal(t, cfg.BusinessLogic.Consumer.Enabled, true)
	assert.Equal(t, cfg.BusinessLogic.Consumer.APIKeys, "consumer_key")
	assert.Equal(t, cfg.BusinessLogic.Consumer.Address, "api.thingspeak.com")
	assert.Equal(t, cfg.BusinessLogic.Consumer.Schema, "https")
	assert.Equal(t, cfg.BusinessLogic.Consumer.URI, "update")

	assert.Equal(t, cfg.WeatherProvider.Enabled, true)

	assert.Equal(t, cfg.ForecastProvider.Enabled, true)
}
