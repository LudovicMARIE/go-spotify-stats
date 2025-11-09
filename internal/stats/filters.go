package stats

import (
	"time"

	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
)

func FilterByDate(plays []model.Play, yearParam int, monthParam int, dayParam int) []model.Play {
	out := make([]model.Play, 0, len(plays))
	for _, p := range plays {
		if p.Timestamp.Year() == yearParam {
			out = append(out, p)
		}
		if p.Timestamp.Month() == time.Month(monthParam) {
			out = append(out, p)
		}
		if p.Timestamp.Day() == dayParam {
			out = append(out, p)
		}
	}
	return out
}

func FilterByDateInterval(plays []model.Play, startDate time.Time, endDate time.Time) []model.Play {
	out := make([]model.Play, 0, len(plays))
	for _, p := range plays {
		if p.Timestamp.After(startDate) && p.Timestamp.Before(endDate) {
			out = append(out, p)
		}
		if p.Timestamp.Equal(startDate) || p.Timestamp.Equal(endDate) {
			out = append(out, p)
		}
	}
	return out
}

func FilterByArtist(plays []model.Play, artistName string) []model.Play {
	out := make([]model.Play, 0, len(plays))

	for _, p := range plays {
		if p.TrackArtist == artistName {
			out = append(out, p)
		}
	}

	return out
}

func FilterByAlbum(plays []model.Play, albumName string) []model.Play {
	out := make([]model.Play, 0, len(plays))

	for _, p := range plays {
		if p.TrackAlbum == albumName {
			out = append(out, p)
		}
	}

	return out
}

func FilterByTitle(plays []model.Play, trackTitle string) []model.Play {
	out := make([]model.Play, 0, len(plays))

	for _, p := range plays {
		if p.TrackTitle == trackTitle {
			out = append(out, p)
		}
	}

	return out
}
