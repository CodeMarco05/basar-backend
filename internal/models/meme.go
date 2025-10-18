package models

type Meme struct {
	Id        string   `json:"id"`
	ImageUrl  string   `json:"imageUrl"`
	AltText   string   `json:"altText"`
	ImageB64  string   `json:"imageB64"`
	Title     string   `json:"title"`
	Permalink string   `json:"permalink"`
	Likes     int      `json:"likes"`
	TimeStamp string   `json:"timeStamp"`
	User      string   `json:"user"`
	Tags      []string `json:"tags"`
	Comments  []string `json:"comments"`
}
