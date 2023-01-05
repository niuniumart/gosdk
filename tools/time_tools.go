package tools

import "time"

const (
	timeTemplate1 = "2006-01-02 15:04:05" //常规类型
	timeTemplate2 = "2006-01-02"          //其他类型
	timeTemplate3 = "2006/01/02 15:04:05" //其他类型
	timeTemplate4 = "20060102 15:04:05"   //其他类型
)

// GetTimeFromStr with format "2006-01-02 15:04:05"
func GetTimeFromStr(timeStr string) (time.Time, error) {
	stamp, err := time.ParseInLocation(timeTemplate1, timeStr, time.Local)
	return stamp, err
}

// GetTimeFromDayStr with format "2006-01-02 15:04:05"
func GetTimeFromDayStr(timeStr string) (time.Time, error) {
	stamp, err := time.ParseInLocation(timeTemplate2, timeStr, time.Local)
	return stamp, err
}

// GetTimeFromSpritStr with format "2006-01-02 15:04:05"
func GetTimeFromSpritStr(timeStr string) (time.Time, error) {
	stamp, err := time.ParseInLocation(timeTemplate3, timeStr, time.Local)
	return stamp, err
}

// GetTimeFromNoSprit with format "20060102 15:04:05"
func GetTimeFromNoSprit(timeStr string) (time.Time, error) {
	stamp, err := time.ParseInLocation(timeTemplate4, timeStr, time.Local)
	return stamp, err
}

// GetStrFromTime with format "2006-01-02 15:04:05"
func GetStrFromTime(timeTime time.Time) string {
	return timeTime.Format(timeTemplate1)
}

// GetDayStrFromTime with format "2006-01-02 15:04:05"
func GetDayStrFromTime(timeTime time.Time) string {
	return timeTime.Format(timeTemplate2)
}

// GetSplitStrFromTime with format "2006-01-02 15:04:05"
func GetSplitStrFromTime(timeTime time.Time) string {
	return timeTime.Format(timeTemplate3)
}

// GetNoSplitFromTime with format "20060102 15:04:05"
func GetNoSplitFromTime(timeTime time.Time) string {
	return timeTime.Format(timeTemplate4)
}
