package filter

import (
	"fmt"
	"time"
	"tvprog/pkg/parser"
)

const (
	MinimumProgramLength = 35 // List programs that are at least 35 mins long
)

type FilteredProgram struct {
	Title   string `json:"title"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Channel string `json:"channel"`
}

type ProgramsByDate map[string]*DayPrograms

type DayPrograms struct {
	Programs     map[string]FilteredProgram
	ChannelOrder []string
}

func parseTimeString(timeStr string) (time.Time, error) {
	layout := "20060102150405 -0700"
	return time.Parse(layout, timeStr)
}

func FilterPrograms(tv *parser.TV) (ProgramsByDate, error) {
	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return nil, fmt.Errorf("failed to load Paris timezone: %w", err)
	}

	channelMap := make(map[string]string)
	for _, channel := range tv.Channels {
		channelMap[channel.ID] = channel.DisplayName
	}

	result := make(ProgramsByDate)

	for _, program := range tv.Programmes {
		start, err := parseTimeString(program.Start)
		if err != nil {
			continue
		}

		end, err := parseTimeString(program.Stop)
		if err != nil {
			continue
		}

		start = start.In(paris)
		end = end.In(paris)

		duration := end.Sub(start)
		if duration.Minutes() < float64(MinimumProgramLength) {
			continue
		}

		hour := start.Hour()
		minute := start.Minute()
		if !((hour == 20 && minute > 45) || (hour == 21 && minute < 20)) {
			continue
		}

		channelName, exists := channelMap[program.Channel]
		if !exists {
			continue
		}

		dateStr := start.Format("2006-01-02")

		filteredProgram := FilteredProgram{
			Title:   program.Title,
			Start:   start.Format("15:04"),
			End:     end.Format("15:04"),
			Channel: channelName,
		}

		if result[dateStr] == nil {
			result[dateStr] = &DayPrograms{
				Programs:     make(map[string]FilteredProgram),
				ChannelOrder: []string{},
			}
		}

		// Add channel to order if not already present
		dayPrograms := result[dateStr]
		if _, exists := dayPrograms.Programs[channelName]; !exists {
			dayPrograms.ChannelOrder = append(dayPrograms.ChannelOrder, channelName)
		}

		dayPrograms.Programs[channelName] = filteredProgram
	}

	return result, nil
}
