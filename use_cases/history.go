package use_cases

import (
	"fmt"
	"time"
	"timeMate_bot/storage"
)

func GetHistory(db *storage.Store, startDate time.Time, endDate time.Time) ([]string, error) {
	events, err := db.GetEventsByDate(startDate, endDate)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(events)+1)
	res = append(res, fmt.Sprintf("     %s\n\n", startDate.Format("02.01.2006")))
	for i, v := range events {
		dateForPrint, err := time.Parse("2006-01-02 15:04:05", v.Date)
		if err != nil {
			return nil, err
		}
		res = append(res, fmt.Sprintf("%4d.   %-10s %6s %s\n",
			i+1, v.Tag, dateForPrint.Format("15:04"), v.Comment))
	}

	if len(events) == 0 {
		res = append(res, fmt.Sprintf("нет событий на эту дату"))
	}

	return res, nil

}
