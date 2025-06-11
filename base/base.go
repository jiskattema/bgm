package base

// To ask the user to enter a string; the callback is called with both the question and the users' answer.
type Question struct {
	Prompt  string
	Key string
	Value string
	Callback func(key, value string)
}

// A search criteria: TAG OP VALUE , 'Title' '=' 'One'
type Constraint struct {
	Tag   string
	Op    string
	Value string
}

// For querying a remote (ie. MPD) for all Tags of songs matching the criteria.
type Query struct {
	Query_id    int
	Tag         string
	Constraints []Constraint
}

// The result of a query
type Result struct {
	Query_id int
	Result   []string
}

// the attributes of a song
type Attrs map[string]string

// a slice containing songs
type Playlist []Attrs

type Tag struct {
	Name        string
	Display     string
	Description string
	Visible     bool
}

var ServerTags = []Tag{
	{Visible: true, Display: "Artist", Name: "artist", Description: "The artist name. Its meaning is not well-defined; see 'composer' and 'performer' for more specific tags."},
	{Visible: false, Display: "Artist*", Name: "artistsort", Description: "Same as artist, but for sorting. This usually omits prefixes such as 'The'."},
	{Visible: true, Display: "Album", Name: "album", Description: "The album name."},
	{Visible: false, Display: "Album*", Name: "albumsort", Description: "Same as album, but for sorting."},
	{Visible: true, Display: "Album artist", Name: "albumartist", Description: "On multi-artist albums, this is the artist name which shall be used for the whole album. The exact meaning of this tag is not well-defined."},
	{Visible: false, Display: "Album artist*", Name: "albumartistsort", Description: "Same as albumartist, but for sorting."},
	{Visible: true, Display: "Title", Name: "title", Description: "The song title."},
	{Visible: false, Display: "Title*", Name: "titlesort", Description: "Same as title, but for sorting."},
	{Visible: true, Display: "Track", Name: "track", Description: "The decimal track number within the album."},
	{Visible: false, Display: "Name", Name: "name", Description: "A name for this song. This is not the song title. The exact meaning of this tag is not well-defined. It is often used by badly configured internet radio stations with broken tags to squeeze both the artist name and the song title in one tag."},
	{Visible: false, Display: "Genre", Name: "genre", Description: "The music genre."},
	{Visible: false, Display: "Mood", Name: "mood", Description: "The mood of the audio with a few keywords."},
	{Visible: false, Display: "Date", Name: "date", Description: "The song’s release date. This is usually a 4-digit year."},
	{Visible: false, Display: "Original date", Name: "originaldate", Description: "The song’s original release date."},
	{Visible: false, Display: "Composer", Name: "composer", Description: "The artist who composed the song."},
	{Visible: false, Display: "Composer*", Name: "composersort", Description: "Same as composer, but for sorting."},
	{Visible: false, Display: "Performer", Name: "performer", Description: "The artist who performed the song."},
	{Visible: false, Display: "Conductor", Name: "conductor", Description: "The conductor who conducted the song."},
	{Visible: false, Display: "Work", Name: "work", Description: "A work is a distinct intellectual or artistic creation, which can be expressed in the form of one or more audio recordings."},
	{Visible: false, Display: "Ensemble", Name: "ensemble", Description: "The ensemble performing this song, e.g. 'Wiener Philharmoniker'."},
	{Visible: false, Display: "Movement", Name: "movement", Description: "Name of the movement, e.g. 'Andante con moto'."},
	{Visible: false, Display: "Movement number", Name: "movementnumber", Description: "Movement number, e.g. '2' or 'II'."},
	{Visible: false, Display: "Show movement", Name: "showmovement", Description: "If this tag is set to '1' players supporting this tag will display the work, movement, and movementnumber` instead of the track title."},
	{Visible: false, Display: "Location", Name: "location", Description: "Location of the recording, e.g. 'Royal Albert Hall'."},
	{Visible: false, Display: "Grouping", Name: "grouping", Description: "Used if the sound belongs to a larger category of sounds/music (from the IDv2.4.0 TIT1 description)."},
	{Visible: false, Display: "Comment", Name: "comment", Description: "A human-readable comment about this song. The exact meaning of this tag is not well-defined."},
	{Visible: false, Display: "Disc", Name: "disc", Description: "The decimal disc number in a multi-disc album."},
	{Visible: true, Display: "Label", Name: "label", Description: "The name of the label or publisher."},
	{Visible: false, Display: "MusicBrainz artistID", Name: "musicbrainz_artistid", Description: "The artist id in the MusicBrainz database."},
	{Visible: false, Display: "MusicBrainz albumID", Name: "musicbrainz_albumid", Description: "The album id in the MusicBrainz database."},
	{Visible: false, Display: "MusicBrainz albumartistID", Name: "musicbrainz_albumartistid", Description: "The album artist id in the MusicBrainz database."},
	{Visible: false, Display: "MusicBrainz trackID", Name: "musicbrainz_trackid", Description: "The track id in the MusicBrainz database."},
	{Visible: false, Display: "MusicBrainz releasegroupID", Name: "musicbrainz_releasegroupid", Description: "The release group id in the MusicBrainz database."},
	{Visible: false, Display: "MusicBrainz releasetrackID", Name: "musicbrainz_releasetrackid", Description: "The release track id in the MusicBrainz database."},
	{Visible: false, Display: "MusicBrainz workID", Name: "musicbrainz_workid", Description: "The work id in the MusicBrainz database."},
}
