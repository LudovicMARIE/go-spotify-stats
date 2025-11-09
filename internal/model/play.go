package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Play struct {
	TrackTitle  string    `json:"master_metadata_track_name"`
	TrackArtist string    `json:"master_metadata_album_artist_name"`
	TrackAlbum  string    `json:"master_metadata_album_name"`
	Timestamp   Timestamp `json:"ts"`
	MsPlayed    int       `json:"ms_played"`
}

type Timestamp struct {
	time.Time
}

// Only accept RFC3339 timestamps like "2023-07-28T13:30:23Z".
var customLayout = time.RFC3339

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "null" || s == `""` {
		t.Time = time.Time{}
		return nil
	}

	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		str := s[1 : len(s)-1]
		parsed, err := time.ParseInLocation(customLayout, str, time.UTC)
		if err != nil {
			return fmt.Errorf("parsing timestamp %q: %w", str, err)
		}
		t.Time = parsed
		return nil
	}

	return fmt.Errorf("invalid timestamp JSON token: %s", s)
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`null`), nil
	}
	return json.Marshal(t.Time.Format(customLayout))
}
