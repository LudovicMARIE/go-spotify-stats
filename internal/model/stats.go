package model

type ArtistStat struct {
	Name        string
	Plays       int
	TotalMillis int64
}

// Used as a map key to ensure unique tracks
// (Because two different artists might have a song named "Intro")
type TrackKey struct {
	Title  string
	Artist string
}

type TrackStat struct {
	Title  string
	Artist string
	Plays  int
	Millis int64
}
