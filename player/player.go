package player

type ItemKind int

const (
	ITEM_SONG = iota
	ITEM_ARTIST
	ITEM_ALBUM
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
	Add(Song)
	// AddNext inserts a song at the front
	AddNext(Song)
	// Search for a term
	Search(string) []Item
}

// Item is a player item
type Item interface {
	Identifier() string
	Kind() ItemKind
}

// State stores the current player state
type State interface {
	IsPlaying() bool
	Current() Song
	Progress() float64
}

// Song is a piece of music
type Song interface {
	Item
	Name() string
	Artist() Artist
	Album() Album
	Length() int
}

// Album is a collection of songs
type Album interface {
	Item
	Name() string
	Artist() Artist
	Songs() []Song
}

// Artist is an creator or interpreter of songs
type Artist interface {
	Item
	Name() string
	Albums() []Album
}

// RemotePlayer is a MPD-based remote player
type RemotePlayer struct {
	Hostname, Port string
}

func (remote *RemotePlayer) Add(song Song) {

}
func (remote *RemotePlayer) AddNext(song Song) {

}
func (remote *RemotePlayer) Clear() {

}
func (remote *RemotePlayer) List() []Song {
	return nil
}
func (remote *RemotePlayer) Move(a, b int) {

}

func (remote *RemotePlayer) Search(term string) []Item {
	return nil
}

func (remote *RemotePlayer) Skip() {

}

func (remote *RemotePlayer) State() State {
	return nil
}

func (remote *RemotePlayer) Play() {

}
func (remote *RemotePlayer) Stop() {

}

func (remote *RemotePlayer) connect() error {
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
