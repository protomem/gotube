package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Comment struct {
	logger logging.Logger
	serv   service.Comment
}

func NewComment(logger logging.Logger, serv service.Comment) *Comment {
	return &Comment{
		logger: logger.With("handler", "comment"),
		serv:   serv,
	}
}

func (h *Comment) List() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		videoIDRaw, ok := mux.Vars(r)["videoId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing video id")
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid video id").WithInternal(err)
		}

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

		comments, err := h.serv.FindByVideo(r.Context(), videoID, findOpts)
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"comments": comments})
	}, h.errorHandler("handler.Comment.List"))
}

func (h *Comment) Create() http.HandlerFunc {
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
			Comment string `json:"comment"`
		}

		if err := httplib.DecodeJSON(r, &request); err != nil {
			return err
		}

		author := ctxstore.MustUser(r.Context())

		comment, err := h.serv.Create(r.Context(), service.CreateCommentDTO{
			Message:  request.Comment,
			VideoID:  videoID,
			AuthorID: author.ID,
		})
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusCreated, httplib.JSON{"comment": comment})
	}, h.errorHandler("handler.Comment.List"))
}

func (h *Comment) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		commentIDRaw, ok := mux.Vars(r)["commentId"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing comment id")
		}

		commentID, err := uuid.Parse(commentIDRaw)
		if err != nil {
			return httplib.NewAPIError(http.StatusBadRequest, "invalid comment id").WithInternal(err)
		}

		if err := h.serv.Delete(r.Context(), commentID); err != nil {
			return err
		}

		return httplib.NoContent(w)
	}, h.errorHandler("handler.Comment.List"))
}

func (h *Comment) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)
		httplib.DefaultErrorHandler(w, r, err)
	}
}
