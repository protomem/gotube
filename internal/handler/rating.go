package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Rating struct {
	logger logging.Logger
	serv   service.Rating
}

func NewRating(logger logging.Logger, serv service.Rating) *Rating {
	return &Rating{
		logger: logger.With("handler", "rating"),
		serv:   serv,
	}
}

func (h *Rating) Count() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

		likes, dislikes, err := h.serv.Count(r.Context(), videoID)
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{
			"likes":    likes,
			"dislikes": dislikes,
		})
	}, h.errorHandler("handler.Rating.Count"))
}

func (h *Rating) Like() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

		user := ctxstore.MustUser(r.Context())

		if err := h.serv.Like(r.Context(), service.RatingDTO{UserID: user.ID, VideoID: videoID}); err != nil {
			return err
		}

		return httplib.NoContent(w)
	}, h.errorHandler("handler.Rating.Like"))
}

func (h *Rating) Dislike() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

		user := ctxstore.MustUser(r.Context())

		if err := h.serv.Dislike(r.Context(), service.RatingDTO{UserID: user.ID, VideoID: videoID}); err != nil {
			return err
		}

		return httplib.NoContent(w)
	}, h.errorHandler("handler.Rating.Dislike"))
}

func (h *Rating) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

		user := ctxstore.MustUser(r.Context())

		if err := h.serv.Delete(r.Context(), service.RatingDTO{UserID: user.ID, VideoID: videoID}); err != nil {
			return err
		}

		return httplib.NoContent(w)
	}, h.errorHandler("handler.Rating.Delete"))
}

func (h *Rating) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)
		httplib.DefaultErrorHandler(w, r, err)
	}
}
