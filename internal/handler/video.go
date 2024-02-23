package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

const (
	_defaultLimit  = 10
	_defaultOffset = 0
)

type Video struct {
	logger logging.Logger
	serv   service.Video
}

func NewVideo(logger logging.Logger, serv service.Video) *Video {
	return &Video{
		logger: logger.With("handler", "video"),
		serv:   serv,
	}
}

func (h *Video) List() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		var limit uint64 = _defaultLimit
		if r.URL.Query().Has("limit") {
			value, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				return httplib.NewAPIError(http.StatusBadRequest, "invalid limit").WithInternal(err)
			}
			limit = value
		}

		var offset uint64 = _defaultOffset
		if r.URL.Query().Has("offset") {
			value, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				return httplib.NewAPIError(http.StatusBadRequest, "invalid offset").WithInternal(err)
			}
			offset = value
		}

		findOpts := service.FindOptions{Limit: limit, Offset: offset}

		var sortBy string = "latest"
		if r.URL.Query().Has("sortBy") {
			sortBy = r.URL.Query().Get("sortBy")
		}

		authorNickname, authorNicknameOk := r.URL.Query().Get("author"), r.URL.Query().Has("author")

		var (
			err    error
			videos []model.Video
		)

		if authorNicknameOk {
			videos, err = h.serv.FindByAuthor(r.Context(), authorNickname, findOpts)
		} else {
			switch sortBy {
			case "latest", "news", "createdAt":
				videos, err = h.serv.FindLatest(r.Context(), findOpts)
			case "popular", "trends", "views":
				videos, err = h.serv.FindPopular(r.Context(), findOpts)
			default:
				return httplib.NewAPIError(http.StatusBadRequest, "invalid sortBy")
			}
		}

		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"videos": videos})
	}, h.errorHandler("handler.Video.List"))
}

func (h *Video) Get() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

		video, err := h.serv.Get(r.Context(), videoID)
		if err != nil {
			return err
		}

		author, isAuth := ctxstore.User(r.Context())
		if video.Public || (!video.Public && isAuth && author.ID == video.Author.ID) {
			return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"video": video})
		}

		return model.ErrVideoNotFound
	}, h.errorHandler("handler.Video.Get"))
}

func (h *Video) Creaate() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		var request struct {
			Title         string
			Description   *string
			ThumbnailPath string
			VideoPath     string
			Public        *bool
		}

		if err := httplib.DecodeJSON(r, &request); err != nil {
			return err
		}

		author := ctxstore.MustUser(r.Context())

		video, err := h.serv.Create(r.Context(), service.CreateVideoDTO{
			Title:         request.Title,
			Description:   request.Description,
			ThumbnailPath: request.ThumbnailPath,
			VideoPath:     request.VideoPath,
			AuthorID:      author.ID,
			Public:        request.Public,
		})
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusCreated, httplib.JSON{"video": video})
	}, h.errorHandler("handler.Video.Create"))
}

func (h *Video) Update() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

		var request struct {
			Title         *string `json:"title"`
			Description   *string `json:"description"`
			ThumbnailPath *string `json:"thumbnailPath"`
			VideoPath     *string `json:"videoPath"`
			Public        *bool   `json:"isPublic"`
		}

		if err := httplib.DecodeJSON(r, &request); err != nil {
			return err
		}

		video, err := h.serv.Update(r.Context(), videoID, service.UpdateVideoDTO(request))
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"video": video})
	}, h.errorHandler("handler.Video.Update"))
}

func (h *Video) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

		if err := h.serv.Delete(r.Context(), videoID); err != nil {
			return err
		}

		return httplib.NoContent(w)
	}, h.errorHandler("handler.Video.Delete"))
}

func (h *Video) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)

		if errors.Is(err, model.ErrVideoNotFound) {
			err = httplib.NewAPIError(http.StatusNotFound, model.ErrVideoNotFound.Error())
		}
		if errors.Is(err, model.ErrVideoExists) {
			err = httplib.NewAPIError(http.StatusConflict, model.ErrVideoExists.Error())
		}
		if errors.Is(err, model.ErrUserNotFound) {
			err = httplib.NewAPIError(http.StatusNotFound, model.ErrUserNotFound.Error())
		}

		httplib.DefaultErrorHandler(w, r, err)
	}
}
