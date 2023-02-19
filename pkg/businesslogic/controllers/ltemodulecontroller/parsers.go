package ltemodulecontroller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

func ParseAccBalance(rawBalance string) (float64, error) {
	re := regexp.MustCompile(`Balans ([\d]+\.[\d]+) grn,`)
	parsed := re.FindStringSubmatch(rawBalance)

	if len(parsed) < 2 || len(parsed[1]) < 3 {
		log.Println("NOTICE: failed to parse balance output:", rawBalance)
		return 0.0, fmt.Errorf("failed to parse balance: %s", rawBalance)
	}

	res, err := strconv.ParseFloat(parsed[1], 64)
	if err != nil {
		log.Println("ERROR: failed to parse balance output: ", rawBalance)
		return 0.0, fmt.Errorf("failed to parse balance output: %s", rawBalance)
	}

	return res, nil
}

func ParseInetBalance(rawBalance string) (float64, error) {
	re := regexp.MustCompile(` ([\d]+\.[\d]+) GB`)
	parsed := re.FindStringSubmatch(rawBalance)

	if len(parsed) < 2 || len(parsed[1]) < 3 {
		log.Println("NOTICE: failed to parse inet balance output:", rawBalance)
		return 0.0, fmt.Errorf("failed to parse inet balance: %s", rawBalance)
	}

	res, err := strconv.ParseFloat(parsed[1], 64)
	if err != nil {
		log.Println("ERROR: failed to parse inet balance output: ", rawBalance)
		return 0.0, fmt.Errorf("failed to parse inet balance output: %s", rawBalance)
	}

	return res, nil
}

func ParseTariff(rawTariff string) (string, error) {
	re := regexp.MustCompile(`Vash taryf:(.+)[\r\n]`)
	parsed := re.FindStringSubmatch(rawTariff)

	if len(parsed) < 2 || len(parsed[1]) < 3 {
		log.Println("NOTICE: failed to parse tariff output:", rawTariff)
		return "", fmt.Errorf("failed to parse tariff output: %s", rawTariff)
	}

	return parsed[1], nil
}

func ParsePhoneNumber(rawPhoneNumber string) (string, error) {
	re := regexp.MustCompile(`380[\d]{9}`)
	parsed := re.FindString(rawPhoneNumber)

	if len(parsed) != 12 {
		log.Println("NOTICE: failed to parse phone number output:", rawPhoneNumber)
		return "", fmt.Errorf("failed to parse phone number output: %s", rawPhoneNumber)
	}

	return parsed, nil
}
