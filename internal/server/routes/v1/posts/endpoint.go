package posts

import (
	"net/http"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all users"))
}
