package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func processCSVData(config CsvConfiguration, data [][]string) {
	for r, row := range data {
		for c, columnValue := range row {
			//			if c+1 == 17 {
			for _, test := range config.Tests {
				if test.Column == c+1 {
					err := processTest(test, columnValue)
					if err != nil {
						fmt.Println("Row " + strconv.Itoa(r+3) + ", Col " + strconv.Itoa(c+1) + ":\n" + err.Error())
					}
				}
			}
			//			}
		}
	}
}

func processRecord(data []string, r int, testsByColumn map[int][]CsvTest) {
	for c, columnValue := range data {
		if tests, ok := testsByColumn[c+1]; ok { // Column index is 1-based
			for _, test := range tests {
				if err := processTest(test, columnValue); err != nil {
					fmt.Println("Row " + strconv.Itoa(r+3) + ", Col " + strconv.Itoa(c+1) + ":\n" + err.Error())
				}
			}
		}
	}
}

func processTest(test CsvTest, data string) error {
	switch test.Test {
	case "BeginsWith":
		return BeginsWith(test, data)
	case "EndsWith":
		return EndsWith(test, data)
	case "HasOneOf":
		return HasOneOf(test, data)
	case "IsDate":
		return IsDate(test, data)
	case "IsDateTime":
		return IsDateTime(test, data)
	case "IsLength":
		return IsLength(test, data)
	case "IsNumber":
		return IsNumber(test, data)
	case "IsNumberDecimal":
		return IsNumberDecimal(test, data)
	case "IsTime":
		return IsTime(test, data)
	case "IsTrimmed":
		return IsTrimmed(test, data)
	case "MaxLength":
		return MaxLength(test, data)
	case "MaxValue":
		return MaxValue(test, data)
	case "MinLength":
		return MinLength(test, data)
	case "MinValue":
		return MinValue(test, data)
	case "NotEmpty":
		return NotEmpty(test, data)
	case "NotNumber":
		return NotNumber(test, data)
	}

	return nil
}

func BeginsWith(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	values := strings.Split(test.Values, ",")
	for _, value := range values {
		if strings.HasPrefix(data, value) {
			return nil
		}
	}

	return errors.New("data does not begin with provided value(s)")
}

func EndsWith(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	values := strings.Split(test.Values, ",")
	for _, value := range values {
		if strings.HasSuffix(data, value) {
			return nil
		}
	}

	return errors.New("data does not end with provided value(s)")
}

func HasOneOf(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	values := strings.Split(test.Values, ",")
	for _, value := range values {
		if value == data {
			return nil
		}
	}

	return errors.New("data not equal to provided values")
}

// Maybe do this by pointer for performance increase?
func mapDateStringToGoString(data string) string {
	data = strings.Replace(data, "yyyy", "2006", -1)
	data = strings.Replace(data, "MM", "01", -1)
	data = strings.Replace(data, "dd", "02", -1)

	return data
}

func mapTimeStringToGoString(data string) string {
	data = strings.Replace(data, "hh", "15", -1)
	data = strings.Replace(data, "mm", "04", -1)
	data = strings.Replace(data, "ss", "05", -1)

	return data
}

func mapPMStringToGoString(data string) string {
	data = strings.Replace(data, "PM", "PM", -1)
	data = strings.Replace(data, "AM", "PM", -1)

	return data
}

func IsDate(test CsvTest, data string) error {
	if len(data) != 0 {
		test.DateTimeValue = mapDateStringToGoString(test.DateTimeValue)

		_, err := time.Parse(test.DateTimeValue, data)
		if err != nil {
			return errors.New("IsDate: data is not a date in the format of " + test.DateTimeValue)
		}
	}

	return nil
}

func IsDateTime(test CsvTest, data string) error {
	if len(data) != 0 {
		test.DateTimeValue = mapDateStringToGoString(test.DateTimeValue)
		test.DateTimeValue = mapTimeStringToGoString(test.DateTimeValue)
		test.DateTimeValue = mapPMStringToGoString(test.DateTimeValue)

		if strings.Contains(data, "PM") && strings.Contains(data, "15") {
			data = strings.Replace(test.DateTimeValue, "PM", "", -1)
		}

		_, err := time.Parse(test.DateTimeValue, data)
		if err != nil {
			return errors.New("IsDate: data " + data + " is not a date in the format of " + test.DateTimeValue)
		}
	}

	return nil
}

func IsTime(test CsvTest, data string) error {
	if len(data) != 0 {
		test.DateTimeValue = mapTimeStringToGoString(test.DateTimeValue)

		_, err := time.Parse(test.DateTimeValue, data)
		if err != nil {
			return errors.New("IsTime: time is not a time in the format of " + test.DateTimeValue)
		}
	}

	return nil
}

func IsLength(test CsvTest, data string) error {
	length := len(data)

	if length != 0 && length != test.Length {
		return errors.New("IsLength: length " + strconv.Itoa(length) +
			" data not equal to configured length " + strconv.Itoa(test.Length))
	}

	return nil
}

func IsNumber(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	_, err := strconv.Atoi(data)
	if err != nil {
		return errors.New("IsNumber: non number value found ***" + data + "***")
	}

	return nil
}

func IsNumberDecimal(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	if _, err := strconv.ParseFloat(data, 64); err != nil {
		return errors.New("IsNumberDecimal: non decimal number value found")
	}

	return nil
}

func IsTrimmed(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	if unicode.IsSpace(rune(data[0])) || unicode.IsSpace(rune(data[len(data)-1])) {
		return errors.New("data is not trimmed")
	}

	return nil
}

func MaxLength(test CsvTest, data string) error {
	length := len(data)

	if length != 0 && length > test.Length {
		return errors.New("MaxLength: length of " + strconv.Itoa(length) +
			" longer than maximum length of " + strconv.Itoa(test.Length))
	}

	return nil
}

func MaxValue(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	v, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return errors.New("MaxValue: data is not a number " + data)
	}

	if v > test.Value {
		return errors.New("MaxValue: " + strconv.FormatInt(v, 10) +
			" higher than maximum value of " + strconv.FormatInt(test.Value, 10))
	}

	return nil
}

func MinLength(test CsvTest, data string) error {
	length := len(data)

	if length != 0 && length < test.Length {
		return errors.New("MinLength: length of " + strconv.Itoa(length) +
			" shorter than minimum length of " + strconv.Itoa(test.Length))
	}

	return nil
}

func MinValue(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	v, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return errors.New("MinValue: data is not a number " + data)
	}

	if v < test.Value {
		return errors.New("MinValue: " + strconv.FormatInt(v, 10) +
			" lower than minimum value of " + strconv.FormatInt(test.Value, 10))
	}

	return nil
}

func NotEmpty(test CsvTest, data string) error {
	if len(data) == 0 {
		return errors.New("NotEmpty: blank value found where no blank allowed")
	}

	return nil
}

func NotNumber(test CsvTest, data string) error {
	if len(data) == 0 {
		return nil
	}

	if _, err := strconv.Atoi(data); err == nil {
		return errors.New("NotNumber: number value found")
	}

	return nil
}
