package player

import (
	"strconv"

	"github.com/fhs/gompd/mpd"
)

// ItemKind stores the item type
type ItemKind int

const (
	// ItemSong is a song item
	ItemSong = iota
	// ItemArtist is a artist item
	ItemArtist
	// ItemAlbum is album item
	ItemAlbum
	ItemUndefined
)

// Player is an abstract audio player
type Player interface {
	// Play current song
	Play()
	// Stop current song
	Stop()
	// Get player state
	State() State
	// List songs in playlist
	List() []Song
	// Clear the playlist
	Clear()
	// Skip current song
	Skip()
	// Move song at place x to place y
	Move(int, int)
	// Add inserts a song at the end
	Add(Item)
	// AddNext inserts a song at the front
	AddNext(Item)
	// Search for a term
	Search(string) []Item
}

// Item is a player item
type Item interface {
	// Identifier is the item's internal URI.
	Identifier() string
	// Kind contains the item's type
	Kind() ItemKind
}

// State stores the current player state
type State interface {
	// IsPlaying is true if a song is currently played
	IsPlaying() bool
	// Current returns the currently played song
	Current() Song
	// Progress stores the song's progress between 0 and 1
	Progress() float64
	Position() int
}

// Song is a piece of music
type Song interface {
	Item
	// The songs name
	Name() string
	// The songs artist
	Artist() Artist
	// The songs album
	Album() Album
	// The songs length in seconds
	Length() float64
}

// Album is a collection of songs
type Album interface {
	Item
	// The albums name
	Name() string
	// The artists information
	Artist() Artist
	// The list of songs in the album
	Songs() []Song
}

// Artist is an creator or interpreter of songs
type Artist interface {
	Item
	// The artist name
	Name() string
	// The artists albums
	Albums() []Album
	// The artists top tracks
	TopTracks() []Song
}

// RemotePlayer is a MPD-based remote player
type RemotePlayer struct {
	Hostname, Port string
	Conn           *mpd.Client
}

type remoteItem struct {
	source    *mpd.Client
	uri, name string
}

func (item *remoteItem) Identifier() string {
	return item.uri
}

func (item *remoteItem) Kind() ItemKind {
	return ItemUndefined
}

type remoteState struct {
	random, consume, playing bool
	current, next            int
	elapsed                  float64
	song                     *remoteSong
	source                   *mpd.Client
}

func (state *remoteState) IsPlaying() bool {
	return state.playing
}

func (state *remoteState) Current() Song {
	return state.song
}

func (state *remoteState) Progress() float64 {
	return state.elapsed / state.song.length
}

func (state *remoteState) Position() int {
	return state.current
}

func parseState(attrs mpd.Attrs, song *remoteSong, client *mpd.Client) *remoteState {
	elapsed, _ := strconv.ParseFloat(attrs["elapsed"], 64)
	current, _ := strconv.Atoi(attrs["song"])
	next, _ := strconv.Atoi(attrs["nextsong"])
	return &remoteState{
		elapsed: elapsed,
		playing: attrs["state"] == "play",
		current: current,
		next:    next,
		random:  attrs["random"] != "0",
		consume: attrs["consume"] != "0",
		source:  client,
		song:    song,
	}
}

type remoteSong struct {
	remoteItem
	albumURI, albumName string
	artistName          string
	length              float64
}

func (song *remoteSong) Kind() ItemKind {
	return ItemSong
}

func (song *remoteSong) Name() string {
	return song.name
}

func (song *remoteSong) Artist() Artist {
	return &remoteArtist{
		remoteItem: remoteItem{
			uri:    song.artistName,
			name:   song.artistName,
			source: song.source,
		},
	}
}

func (song *remoteSong) Length() float64 {
	return song.length
}

func (song *remoteSong) Album() Album {
	return &remoteAlbum{
		remoteItem: remoteItem{
			uri:    song.albumURI,
			name:   song.albumURI,
			source: song.source,
		},
	}
}

func parseSong(attrs mpd.Attrs, client *mpd.Client) *remoteSong {
	len, _ := strconv.ParseFloat(attrs["Time"], 64)
	return &remoteSong{
		remoteItem: remoteItem{
			uri:    attrs["file"],
			name:   attrs["Title"],
			source: client,
		},
		artistName: attrs["Artist"],
		albumName:  attrs["Album"],
		albumURI:   attrs["X-AlbumUri"],
		length:     len,
	}
}

type remoteAlbum struct {
	remoteItem
}

func (album *remoteAlbum) Kind() ItemKind {
	return ItemAlbum
}

func (album *remoteAlbum) Artist() Artist {
	// TODO
	return nil
}

func (album *remoteAlbum) Name() string {
	return album.name
}

func (album *remoteAlbum) Songs() []Song {
	// TODO
	return []Song{}
}

type remoteArtist struct {
	remoteItem
}

func (artist *remoteArtist) Kind() ItemKind {
	return ItemArtist
}

func (artist *remoteArtist) Name() string {
	return artist.name
}

func (artist *remoteArtist) Albums() []Album {
	// TODO
	return []Album{}
}

func (artist *remoteArtist) TopTracks() []Song {
	// TODO
	return []Song{}
}

func (remote *RemotePlayer) Add(item Item) {
	remote.Conn.Add(item.Identifier())
}
func (remote *RemotePlayer) AddNext(item Item) {
	remote.Conn.AddID(item.Identifier(), 1)
}
func (remote *RemotePlayer) Clear() {
	remote.Conn.Clear()
}
func (remote *RemotePlayer) List() []Song {
	return nil
}
func (remote *RemotePlayer) Move(a, b int) {
	remote.Conn.MoveID(a, b)
}

func (remote *RemotePlayer) Search(term string) []Item {
	return nil
}

func (remote *RemotePlayer) Skip() {
	remote.Conn.Next()
}

func (remote *RemotePlayer) State() State {
	current, _ := remote.Conn.CurrentSong()
	status, _ := remote.Conn.Status()
	state := parseState(status, parseSong(current, remote.Conn), remote.Conn)
	return state
}

func (remote *RemotePlayer) Play() {
	remote.Conn.Play(-1)
}
func (remote *RemotePlayer) Stop() {
	remote.Conn.Stop()
}

func (remote *RemotePlayer) ItemByURI(uri string) Item {
	return &remoteItem{
		uri:  uri,
		name: uri,
	}
}

func (remote *RemotePlayer) connect() error {
	conn, err := mpd.Dial("tcp", remote.Hostname+":"+remote.Port)
	if err != nil {
		return err
	}
	remote.Conn = conn
	return nil
}

// Connect to a remote player
func Connect(host, port string) (*RemotePlayer, error) {
	player := &RemotePlayer{Hostname: host, Port: port}
	if err := player.connect(); err != nil {
		return nil, err
	}
	return player, nil
}
