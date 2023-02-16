package weatherio

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateConvertor(t *testing.T) {
	input := "2023-01-31"
	expected := time.Date(2023, time.January, 31, 0, 0, 0, 0, time.Local)

	actual, err := ParseDate(input)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestDateConvertor_Fail1(t *testing.T) {
	input := "2023-01-312"

	_, err := ParseDate(input)
	assert.Error(t, err)
}

func TestDateConvertor_Fail2(t *testing.T) {
	input := "2023-13-31"

	_, err := ParseDate(input)
	assert.Error(t, err)
}

func TestGetMin(t *testing.T) {
	f1 := ResponseDataItem{
		MinTemp:  -1,
		MaxTemp:  -2,
		LowTemp:  -3,
		HighTemp: 4,
	}

	f2 := ResponseDataItem{
		MinTemp:  2,
		MaxTemp:  -8,
		LowTemp:  -3,
		HighTemp: -1,
	}

	expected := -8.0
	actual := GetMinTemperature(f1, f2)

	assert.Equal(t, expected, actual)
}
