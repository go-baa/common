package now

import "time"

// BeginningOfHour truncate time hour
func BeginningOfHour(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

// BeginningOfDay truncate time day
func BeginningOfDay(t time.Time) time.Time {
	d := time.Duration(-t.Hour()) * time.Hour
	return BeginningOfHour(t).Add(d)
}

// EndOfDay get the last time of day
func EndOfDay(t time.Time) time.Time {
	return BeginningOfDay(t).Add(24*time.Hour - time.Nanosecond)
}

// TimeQuantum 获取时间区间，例：近3天
func TimeQuantum(timeQuantum string) (start time.Time, end time.Time) {
	now := time.Now()
	beginningOfNow := BeginningOfDay(now)
	endOfNow := EndOfDay(now)
	endOfYesterday := beginningOfNow.Add(-time.Nanosecond)

	switch timeQuantum {
	case "today": // 今天
		start = beginningOfNow
		end = endOfNow
	case "week": // 近一周
		start = beginningOfNow.AddDate(0, 0, -7)
		end = endOfYesterday
	case "month": // 近一月
		start = beginningOfNow.AddDate(0, 0, -30)
		end = endOfYesterday
	case "quarter": // 近一季
		start = beginningOfNow.AddDate(0, 0, -90)
		end = endOfYesterday
	case "year": // 近一年
		start = beginningOfNow.AddDate(-1, 0, 0)
		end = endOfYesterday
	}

	return
}

const (
	// AnalyticsTimeUnitHour hour
	AnalyticsTimeUnitHour = "hour"
	// AnalyticsTimeUnitDay day
	AnalyticsTimeUnitDay = "day"
	// AnalyticsTimeUnitMonth month
	AnalyticsTimeUnitMonth = "month"
	// AnalyticsTimeUnitYear year
	AnalyticsTimeUnitYear = "year"
)

// GetAnalyticsTimeUnit 时间划分工具
func GetAnalyticsTimeUnit(start time.Time, end time.Time) string {
	day := 24 * time.Hour
	sub := end.Sub(start)
	if sub <= day {
		return AnalyticsTimeUnitHour
	}
	if sub > day && sub <= 90*day {
		return AnalyticsTimeUnitDay
	}
	if sub > 90*day && start.Sub(end.AddDate(-1, 0, 0)) >= 0 {
		return AnalyticsTimeUnitMonth
	}
	if start.Sub(end.AddDate(-1, 0, 0)) < 0 {
		return AnalyticsTimeUnitYear
	}

	return AnalyticsTimeUnitDay
}

// FormatAnalyticsTime 格式化
func FormatAnalyticsTime(timeUnit string, start time.Time, end time.Time, rows []time.Time) ([]string, map[string]time.Time) {
	start = BeginningOfDay(start)
	end = EndOfDay(end)

	var mapTimeLayout string
	var jsonTimeLayout string
	switch timeUnit {
	case AnalyticsTimeUnitHour:
		jsonTimeLayout = "15:00"
		mapTimeLayout = "2006-01-02 15"
	case AnalyticsTimeUnitDay:
		jsonTimeLayout = "01-02"
		mapTimeLayout = "2006-01-02"
	case AnalyticsTimeUnitMonth:
		jsonTimeLayout = "2006-01"
		mapTimeLayout = "2006-01"
	case AnalyticsTimeUnitYear:
	}
	kvs := make(map[string]time.Time)
	for _, v := range rows {
		key := v.Format(mapTimeLayout)
		kvs[key] = v
	}
	timerStrs := make(map[string]time.Time)
	key := make([]string, 0)
	for end.Sub(start) >= 0 {
		str := start.Format(jsonTimeLayout)
		key = append(key, str)
		timerStrs[str] = kvs[start.Format(mapTimeLayout)]

		switch timeUnit {
		case AnalyticsTimeUnitHour:
			start = start.Add(time.Hour)
		case AnalyticsTimeUnitDay:
			start = start.AddDate(0, 0, 1)
		case AnalyticsTimeUnitMonth:
			start = start.AddDate(0, 1, 0)
		case AnalyticsTimeUnitYear:
		}
	}

	return key, timerStrs
}
