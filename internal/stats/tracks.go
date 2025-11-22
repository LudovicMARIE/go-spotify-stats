package stats

import (
	"sort"

	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
)

func ComputeTopTracks(plays []model.Play, filterArtist string) []model.TrackStat {
	counts := map[model.TrackKey]int{}
	minutes := map[model.TrackKey]int64{}

	for _, p := range plays {
		if filterArtist != "" && p.TrackArtist != filterArtist {
			continue
		}
		k := model.TrackKey{Title: p.TrackTitle, Artist: p.TrackArtist}
		counts[k]++
		minutes[k] += int64(p.MsPlayed)
	}

	var out []model.TrackStat
	for k, c := range counts {
		out = append(out, model.TrackStat{
			Title:  k.Title,
			Artist: k.Artist,
			Plays:  c,
			Millis: minutes[k],
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Plays > out[j].Plays
	})

	return out
}
