package model

type Play struct {
	TrackTitle  string `json:"master_metadata_track_name"`
	TrackArtist string `json:"master_metadata_album_artist_name"`
	TrackAlbum  string `json:"master_metadata_album_name"`
	Timestamp   string `json:"ts"`
	MsPlayed    int    `json:"ms_played"`
}
