package util

import (
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)


func StringToPgInterval(s string) pgtype.Interval {
	s = strings.ToUpper(s)
	var microseconds int64 = 0
	var days int32 = 0
	var months int32 = 0

	const multiplier = 1000000
	v := strings.Split(s, " ")
	n, err := strconv.Atoi(v[0])
	if err != nil {
		return pgtype.Interval{}
	}

	switch {
	case strings.Contains(s, "SEC"):
		microseconds += int64(n) * multiplier
	case strings.Contains(s, "MIN"):
		microseconds += int64(n) * 60 * multiplier
	case strings.Contains(s, "HOUR"):
		microseconds += int64(n) * 60 * 60 * multiplier
	case strings.Contains(s, "DAY"):
		days += int32(n)
	case strings.Contains(s, "WEEK"):
		days += int32(n) * 7
	case strings.Contains(s, "MONTH"):
		months += int32(n)
	case strings.Contains(s, "YEAR"):
		months += int32(n) * 12
	}

	return pgtype.Interval{Microseconds: microseconds, Days: days, Months: months, Valid: true}
}