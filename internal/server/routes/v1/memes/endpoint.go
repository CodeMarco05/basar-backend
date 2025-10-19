package memes

import (
	"encoding/json"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMemePage(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	ctx := r.Context()

	pageStr := chi.URLParam(r, "page")
	if pageStr == "" {
		log.Error().Msgf("Getting meme failed: page parameter is missing")
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Error().Msgf("Getting meme failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}

	// page size for paging
	pageSize := 10

	skip := (page - 1) * pageSize

	//get total item count
	collection := config.MongoDB.Collection("memes")
	totalItems, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Error().Msgf("Getting meme failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// calculate total Pages
	totalPages := int(totalItems) / pageSize

	// if the page has overhang it creates a new page
	if int(totalPages)%pageSize != 0 {
		totalPages++
	}

	// query db
	opts := options.Find().
		SetSort(bson.D{{Key: "id", Value: 1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	// execute the query
	res, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Error().Msgf("Getting meme failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer res.Close(ctx)

	// decode results
	var items []models.Meme
	if err := res.All(ctx, &items); err != nil {
		log.Error().Msgf("Getting meme failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if items == nil {
		items = []models.Meme{}
	}

	resp := models.MemeResponse{
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
		MemeList:   items,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error().Msgf("Getting meme failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
