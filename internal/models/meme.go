package models

type Meme struct {
	Id        string   `json:"id"`
	ImageUrl  string   `json:"imageUrl" bson:"image_url"`
	AltText   string   `json:"altText" bson:"alt_text"`
	ImageB64  string   `json:"imageB64" bson:"image_base64"`
	Title     string   `json:"title"`
	Permalink string   `json:"permalink"`
	Likes     string   `json:"likes"`
	TimeStamp string   `json:"timeStamp"`
	User      string   `json:"user"`
	Tags      []string `json:"tags"`
	Comments  []string `json:"comments"`
}

type MemeResponse struct {
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`

	MemeList []Meme `json:"memeList"`
}
