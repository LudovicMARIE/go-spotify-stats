package ui

import (
	"fmt"
	"strings"

	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
	"github.com/LudovicMARIE/go-spotify-stats/internal/stats"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunTUI(allPlays *[]model.Play) error {
	app := tview.NewApplication()

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	artistsList := tview.NewList().ShowSecondaryText(false)
	tracksTable := tview.NewTable().SetFixed(1, 0).SetSelectable(true, false)
	details := tview.NewTextView().SetDynamicColors(true)
	footer := tview.NewTextView().SetTextAlign(tview.AlignLeft)

	leftCol := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Top Artists").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(artistsList, 0, 1, true)

	centerCol := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Top Tracks").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(tracksTable, 0, 1, false)

	rightCol := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Details").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(details, 0, 1, false)

	main := tview.NewFlex().
		AddItem(leftCol, 30, 0, true).
		AddItem(centerCol, 0, 2, false).
		AddItem(rightCol, 40, 0, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(main, 0, 1, true).
		AddItem(footer, 1, 0, false)

	var populateTracks func(filterArtist string)
	var renderDetailsForArtist func(artist string)
	var renderAsciiBars func(values []int, labels []string, width int) string
	var formatMillis func(ms int64) string

	renderHeader := func() {
		totalPlays := len(*allPlays)
		header.SetText(fmt.Sprintf(" go-spotify-stats  —  Plays: %d  —  q:quit  r:refresh  Enter:show details - Navigate with arrows", totalPlays))
	}

	renderFooter := func(msg string) {
		footer.SetText(msg)
	}

	populateArtists := func() {
		artistsList.Clear()
		top := stats.ComputeTopArtists(*allPlays)
		limit := 50
		if len(top) < limit {
			limit = len(top)
		}
		for i := 0; i < limit; i++ {
			a := top[i]
			label := fmt.Sprintf("%d. %s (%d)", i+1, a.Name, a.Plays)
			name := a.Name
			artistsList.AddItem(label, "", 0, func() {
				populateTracks(name)
				renderDetailsForArtist(name)
			})
		}
		artistsList.AddItem("[All artists]", "", 0, func() {
			populateTracks("")
			details.SetText("")
		})
	}

	populateTracks = func(filterArtist string) {
		tracksTable.Clear()
		tracksTable.SetCell(0, 0, tview.NewTableCell("Rank").SetSelectable(false).SetAttributes(tcell.AttrBold))
		tracksTable.SetCell(0, 1, tview.NewTableCell("Title").SetSelectable(false).SetAttributes(tcell.AttrBold))
		tracksTable.SetCell(0, 2, tview.NewTableCell("Artist").SetSelectable(false).SetAttributes(tcell.AttrBold))
		tracksTable.SetCell(0, 3, tview.NewTableCell("Plays").SetSelectable(false).SetAttributes(tcell.AttrBold))
		top := stats.ComputeTopTracks(*allPlays, filterArtist)
		limit := 100
		if len(top) < limit {
			limit = len(top)
		}
		for i := 0; i < limit; i++ {
			tk := top[i]
			tracksTable.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", i+1)))
			tracksTable.SetCell(i+1, 1, tview.NewTableCell(tk.Title))
			tracksTable.SetCell(i+1, 2, tview.NewTableCell(tk.Artist))
			tracksTable.SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf("%d", tk.Plays)))
		}
	}

	renderDetailsForArtist = func(artist string) {
		var b strings.Builder
		top := stats.ComputeTopTracks(*allPlays, artist)
		totalPlays := 0
		totalMillis := int64(0)
		for _, t := range top {
			totalPlays += t.Plays
			totalMillis += t.Millis
		}
		b.WriteString(fmt.Sprintf("[yellow]Artist:[white] %s\n", artist))
		b.WriteString(fmt.Sprintf("[yellow]Plays (top tracks):[white] %d\n", totalPlays))
		b.WriteString(fmt.Sprintf("[yellow]Total listening time (top tracks):[white] %s\n\n", formatMillis(totalMillis)))

		b.WriteString("[green]Top tracks:\n")
		limit := 5
		if len(top) < limit {
			limit = len(top)
		}
		for i := 0; i < limit; i++ {
			tk := top[i]
			b.WriteString(fmt.Sprintf("  %d. %s (%d)\n", i+1, tk.Title, tk.Plays))
		}
		series, labels := stats.ComputeMonthlySeries(*allPlays, artist, 12)
		b.WriteString("\n[yellow]Monthly plays (last 12 months):\n")
		b.WriteString(renderAsciiBars(series, labels, 20))

		details.SetText(b.String())
	}

	// format milliseconds to human readable
	formatMillis = func(ms int64) string {
		seconds := ms / 1000
		minutes := seconds / 60
		hours := minutes / 60
		if hours > 0 {
			return fmt.Sprintf("%dh%02dm", hours, minutes%60)
		}
		if minutes > 0 {
			return fmt.Sprintf("%dm%02ds", minutes%60, seconds%60)
		}
		return fmt.Sprintf("%ds", seconds)
	}

	renderAsciiBars = func(values []int, labels []string, width int) string {
		if len(values) == 0 {
			return ""
		}
		max := 0
		for _, v := range values {
			if v > max {
				max = v
			}
		}
		if max == 0 {
			var sb strings.Builder
			for i, lbl := range labels {
				sb.WriteString(fmt.Sprintf("%s | \n", lbl))
				if i >= 11 {
					break
				}
			}
			return sb.String()
		}
		var sb strings.Builder
		step := float64(width) / float64(max)
		for i, v := range values {
			barLen := int(float64(v) * step)
			if barLen < 0 {
				barLen = 0
			}
			bar := strings.Repeat("█", barLen)
			sb.WriteString(fmt.Sprintf("%s |%s %d\n", labels[i], bar, v))
			if i >= 11 {
				break
			}
		}
		return sb.String()
	}

	renderHeader()
	renderFooter("Ready.")
	populateArtists()
	populateTracks("")
	details.SetText("Select an artist to see details.")

	focusOrder := []tview.Primitive{artistsList, tracksTable, details}
	currentFocus := 0
	focusNext := func() {
		currentFocus = (currentFocus + 1) % len(focusOrder)
		app.SetFocus(focusOrder[currentFocus])
	}
	focusPrev := func() {
		currentFocus = (currentFocus - 1 + len(focusOrder)) % len(focusOrder)
		app.SetFocus(focusOrder[currentFocus])
	}

	artistsList.SetChangedFunc(func(index int, mainText string, secondary string, shortcut rune) {
		if index < 0 {
			return
		}
		name := extractNameFromLabel(mainText)
		populateTracks(name)
		if name != "" {
			renderDetailsForArtist(name)
		} else {
			details.SetText("")
		}
	})

	tracksTable.SetSelectedFunc(func(row, column int) {
		if row <= 0 {
			return
		}
		title := tracksTable.GetCell(row, 1).Text
		artist := tracksTable.GetCell(row, 2).Text
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("[yellow]Track:[white] %s\n", title))
		sb.WriteString(fmt.Sprintf("[yellow]Artist:[white] %s\n", artist))
		ptop := stats.ComputeTopTracks(*allPlays, "")
		found := false
		for _, t := range ptop {
			if t.Title == title && t.Artist == artist {
				sb.WriteString(fmt.Sprintf("[yellow]Plays:[white] %d\n", t.Plays))
				sb.WriteString(fmt.Sprintf("[yellow]Total time:[white] %s\n", formatMillis(t.Millis)))
				found = true
				break
			}
		}
		if !found {
			sb.WriteString("[yellow]Plays:[white] 0\n")
		}
		details.SetText(sb.String())
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			focusNext()
			return nil
		case tcell.KeyBacktab: // Shift + Tab
			focusPrev()
			return nil
		case tcell.KeyRight:
			focusNext()
			return nil
		case tcell.KeyLeft:
			focusPrev()
			return nil
		case tcell.KeyEscape:
			currentFocus = 0
			app.SetFocus(artistsList)
			return nil
		}

		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
			return nil
		case 'r', 'R':
			populateArtists()
			populateTracks("")
			renderFooter("Refreshed.")
			return nil
		case 'h', 'H':
			focusPrev()
			return nil
		case 'l', 'L':
			focusNext()
			return nil
		}
		return event
	})

	if err := app.SetRoot(layout, true).SetFocus(artistsList).Run(); err != nil {
		return err
	}
	return nil
}

func extractNameFromLabel(label string) string {
	parts := strings.SplitN(label, ". ", 2)
	if len(parts) == 2 {
		label = parts[1]
	}
	if idx := strings.LastIndex(label, " ("); idx != -1 {
		return label[:idx]
	}
	if strings.HasPrefix(label, "[All") {
		return ""
	}
	return label
}
