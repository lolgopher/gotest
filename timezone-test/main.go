package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func toUnixtime(datetime string) {

	t, _ := time.Parse(time.RFC3339, datetime)
	fmt.Printf("%s, %d\n", datetime, t.Unix())
}

func toRfc3339(unixtime int64) {
	t := time.Unix(unixtime, 0)
	fmt.Printf("%d, %v\n", unixtime, t)
}

func main() {

	toUnixtime("2023-09-30T00:00:00+09:00")
	toUnixtime("2023-09-30T00:00:00+00:00")
	toRfc3339(1695999600)

	testTime := time.Unix(1696487950, 0)
	loc := time.FixedZone("UTC-A", 9*60*60)
	fmt.Println(testTime)
	fmt.Println(testTime.In(loc))

	t, _ := time.Parse(time.RFC3339, "2023-10-05T06:39:10+09:00")
	fmt.Println(t.Unix())
	fmt.Println("---")

	zone := "+09:00"
	metering, rollup, err := getRollupDate(time.Now(), zone)
	location, _ := getLocation(zone)

	fmt.Println("now:\t\t", time.Now())
	fmt.Println("now(KST):\t", time.Now().In(location))
	fmt.Println("metering:\t", metering)
	fmt.Println("rollup:\t\t", rollup)

	fmt.Println("rollup(KST):\t", rollup.In(location))
	fmt.Println("err:\t\t", err)

}

func getRollupDate(t time.Time, zone string) (time.Time, time.Time, error) {
	// timezone 반영
	location, err := getLocation(zone)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "fail to get location")
	}
	metering := t.In(location)

	// RFC3339 형식으로 rollup 계산
	year, month, day := metering.Date()
	timeStr := fmt.Sprintf("%04d-%02d-%02dT00:00:00%s", year, month, day, zone)
	rollup, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "fail to parse RFC3339")
	}

	// rollup는 UTC로 반환
	return metering, rollup.UTC(), nil
}

func getLocation(zone string) (*time.Location, error) {
	var sign rune
	var hours, minutes int
	parseNum, err := fmt.Sscanf(zone, "%c%02d:%02d", &sign, &hours, &minutes)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse zone")
	}
	if parseNum != 3 {
		return nil, errors.New("invalid offset format")
	}

	offset := hours*60*60 + minutes*60
	if sign == '-' {
		offset *= -1
	} else if sign != '+' {
		return nil, errors.New("invalid sign format")
	}

	return time.FixedZone("UTC", offset), nil
}
