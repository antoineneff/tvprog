package formatter

import (
	"fmt"
	"math"
	"strings"
	"time"
	"tvprog/pkg/filter"
	"unicode/utf8"
)

func maxChannelLength(dayPrograms *filter.DayPrograms) int {
	maxLen := utf8.RuneCountInString("Chaine")
	for _, program := range dayPrograms.Programs {
		if utf8.RuneCountInString(program.Channel) > maxLen {
			maxLen = utf8.RuneCountInString(program.Channel)
		}
	}
	return maxLen + 2
}

func maxTitleLength(dayPrograms *filter.DayPrograms) int {
	maxLen := utf8.RuneCountInString("Titre")
	for _, program := range dayPrograms.Programs {
		title := program.Title
		if utf8.RuneCountInString(title) > 55 {
			runes := []rune(title)
			title = string(runes[:52]) + "..."
		}
		titleLen := utf8.RuneCountInString(title)
		if titleLen > maxLen {
			maxLen = titleLen
		}
	}
	return maxLen + 2
}

func FormatTable(dayPrograms *filter.DayPrograms) string {
	maxChannelLen := maxChannelLength(dayPrograms)
	titleLength := maxTitleLength(dayPrograms)
	timetableLength := len("00:00 - 00:00") + 2
	lineLength := maxChannelLen + titleLength + timetableLength + 4

	now := time.Now()
	paris, _ := time.LoadLocation("Europe/Paris")
	parisTime := now.In(paris)
	programTitle := fmt.Sprintf("PROGRAMME TV DU %s", parisTime.Format("02/01/2006"))

	spacesBeforeTitle := int(math.Ceil(float64(lineLength-len(programTitle)-2)/2)) - 1
	spacesAfterTitle := spacesBeforeTitle
	if (lineLength-len(programTitle))%2 != 0 {
		spacesAfterTitle = (lineLength-len(programTitle)-2)/2 - 1
	}

	var result strings.Builder

	result.WriteString(strings.Repeat(" ", spacesBeforeTitle))
	result.WriteString("┌" + strings.Repeat("─", len(programTitle)+2) + "┐")
	result.WriteString(strings.Repeat(" ", spacesAfterTitle) + "\n")

	result.WriteString(strings.Repeat(" ", spacesBeforeTitle))
	result.WriteString("│ " + programTitle + " │")
	result.WriteString(strings.Repeat(" ", spacesAfterTitle) + "\n")

	result.WriteString("┌" + strings.Repeat("─", maxChannelLen) + "┬")
	result.WriteString(strings.Repeat("─", spacesBeforeTitle-maxChannelLen-2) + "┴")
	result.WriteString(strings.Repeat("─", len(programTitle)+2) + "┴")
	result.WriteString(strings.Repeat("─", titleLength-len(programTitle)-(spacesBeforeTitle-maxChannelLen)-2) + "┬")
	result.WriteString(strings.Repeat("─", timetableLength) + "┐\n")

	result.WriteString("│ Chaine" + strings.Repeat(" ", maxChannelLen-utf8.RuneCountInString("chaine")-1) + "│")
	result.WriteString(" Titre" + strings.Repeat(" ", titleLength-utf8.RuneCountInString("titre")-1) + "│")
	result.WriteString(" Horaires" + strings.Repeat(" ", timetableLength-utf8.RuneCountInString("horaires")-1) + "│\n")

	result.WriteString("├" + strings.Repeat("─", maxChannelLen) + "┼")
	result.WriteString(strings.Repeat("─", titleLength) + "┼")
	result.WriteString(strings.Repeat("─", timetableLength) + "┤\n")

	for _, channel := range dayPrograms.ChannelOrder {
		program := dayPrograms.Programs[channel]
		title := program.Title
		if utf8.RuneCountInString(title) > 55 {
			runes := []rune(title)
			title = string(runes[:52]) + "..."
		}

		result.WriteString("│ " + channel + strings.Repeat(" ", maxChannelLen-utf8.RuneCountInString(channel)-2) + " │")
		result.WriteString(" " + title + strings.Repeat(" ", titleLength-utf8.RuneCountInString(title)-2) + " │")
		result.WriteString(" " + program.Start + " - " + program.End + " │\n")
	}

	result.WriteString("└" + strings.Repeat("─", maxChannelLen) + "┴")
	result.WriteString(strings.Repeat("─", titleLength) + "┴")
	result.WriteString(strings.Repeat("─", timetableLength) + "┘\n")

	return result.String()
}
