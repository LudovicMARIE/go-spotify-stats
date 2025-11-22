package stats

import (
	"sort"
	"time"

	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
)

func ComputeTopArtists(plays []model.Play) []model.ArtistStat {
	counts := map[string]int{}
	minutes := map[string]int64{}

	for _, p := range plays {
		counts[p.TrackArtist]++
		minutes[p.TrackArtist] += int64(p.MsPlayed)
	}

	var out []model.ArtistStat
	for name, c := range counts {
		out = append(out, model.ArtistStat{
			Name:        name,
			Plays:       c,
			TotalMillis: minutes[name],
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Plays > out[j].Plays
	})

	return out
}

func ComputeMonthlySeries(plays []model.Play, filterArtist string, nMonths int) ([]int, []string) {
	now := time.Now().UTC()
	buckets := map[string]int{}
	labels := []string{}
	for i := 0; i < nMonths; i++ {
		t := now.AddDate(0, -i, 0)
		key := t.Format("2006-01")
		labels = append([]string{key}, labels...)
		buckets[key] = 0
	}
	for _, p := range plays {
		if filterArtist != "" && p.TrackArtist != filterArtist {
			continue
		}
		ym := p.Timestamp.UTC().Format("2006-01")
		if _, ok := buckets[ym]; ok {
			buckets[ym]++
		}
	}
	out := make([]int, len(labels))
	for i, key := range labels {
		out[i] = buckets[key]
	}
	return out, labels
}
