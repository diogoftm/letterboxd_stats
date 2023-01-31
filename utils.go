package lbstats

import (
	"strconv"
	"time"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func SingleAtoi(str string) int {
	val, err := strconv.Atoi(str)

	if err != nil {
		panic(err)
	}

	return val
}

func SingleParseFloat(str string) float32 {
	val, err := strconv.ParseFloat(str, 32)

	if err != nil {
		panic(err)
	}

	return float32(val)
}

func ParseToDate(str string) time.Time {
	date, error := time.Parse("2006-01-02", str)
	if error != nil {
		panic("Wrong string format for date")
	}
	return date
}
