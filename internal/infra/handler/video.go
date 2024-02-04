package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/domain/usecase"
	"github.com/protomem/gotube/pkg/request"
	"github.com/protomem/gotube/pkg/response"
	"github.com/protomem/gotube/pkg/validation"
)

const (
	_defaultLimit  = 10
	_defaultOffset = 0
)

type Video struct {
	*Base

	userAcc  port.UserAccessor
	videoAcc port.VideoAccessor
	videoMut port.VideoMutator
}

func NewVideo(userAcc port.UserAccessor, videoAcc port.VideoAccessor, videoMut port.VideoMutator) *Video {
	return &Video{
		Base: NewBase(),

		userAcc:  userAcc,
		videoAcc: videoAcc,
		videoMut: videoMut,
	}
}

func (h *Video) HandleFind(w http.ResponseWriter, r *http.Request) {
	findOpts, findOptsOk := h.getFindOptionsFromRequest(r)
	if !findOptsOk {
		h.BadRequest(w, r, errors.New("invalid limit or offset"))
		return
	}

	searchQuery, searchQueryOk := h.getSearchQueryFromRequest(r)
	authorNickname, authorNicknameOk := h.getAuthorNicknameFromRequest(r)

	sortBy := h.defaultSortByFromRequest(r, "latest")

	requester, isAuth := ctxstore.User(r.Context())

	var (
		err    error
		videos []entity.Video
	)

	var author entity.User
	if authorNicknameOk {
		author, err = h.userAcc.ByNickname(r.Context(), authorNickname)
		if err != nil {
			if entity.IsError(err, entity.ErrUserNotFound) {
				h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrUserNotFound.Error(), nil)
				return
			}

			h.ServerError(w, r, err)
			return
		}
	}

	if searchQueryOk {
		if authorNicknameOk {
			videos, err = h.videoAcc.AllByLikeTitleAndByAuthorIDAndWherePublic(r.Context(), searchQuery, author.ID, findOpts)
		} else {
			videos, err = h.videoAcc.AllByLikeTitleAndWherePublic(r.Context(), searchQuery, findOpts)
		}
	} else if authorNicknameOk {
		if r.URL.Query().Has("private") {
			if !isAuth || requester.ID != author.ID {
				h.ErrorMessage(w, r, http.StatusForbidden, "access denied", nil)
				return
			}

			videos, err = h.videoAcc.AllByAuthorID(r.Context(), author.ID, findOpts)
		} else {
			videos, err = h.videoAcc.AllByAuthorIDAndWherePublic(r.Context(), author.ID, findOpts)
		}
	} else {
		switch sortBy {
		case "latest", "news", "created_at":
			videos, err = h.videoAcc.AllWherePublic(r.Context(), findOpts)
		case "popular", "trends", "views":
			videos, err = h.videoAcc.AllWherePublicAndSortByViews(r.Context(), findOpts)
		}
	}

	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, response.Data{"videos": videos})
}

func (h *Video) HandleGet(w http.ResponseWriter, r *http.Request) {
	videoID, ok := h.getVideoIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}

	video, err := h.videoAcc.ByID(r.Context(), videoID)
	if err != nil {
		if entity.IsError(err, entity.ErrVideoNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrVideoNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	requester, isAuth := ctxstore.User(r.Context())
	if !video.Public && (!isAuth || requester.ID != video.Author.ID) {
		h.ErrorMessage(w, r, http.StatusForbidden, "access denied", nil)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, response.Data{"video": video})
}

func (h *Video) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input port.CreateVideoInput
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	requester := ctxstore.MustUser(r.Context())
	input.AuthorID = requester.ID

	deps := usecase.CreateVideoDeps{
		Accessor: h.videoAcc,
		Mutator:  h.videoMut,
	}
	video, err := usecase.CreateVideo(deps).Invoke(r.Context(), input)
	if err != nil {
		var v *validation.Validator
		if errors.As(err, &v) {
			h.FailedValidation(w, r, v)
			return
		}

		if entity.IsError(err, entity.ErrVideoAlreadyExists) {
			h.ErrorMessage(w, r, http.StatusConflict, entity.ErrVideoAlreadyExists.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusCreated, response.Data{"video": video})
}

func (h *Video) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	videoID, ok := h.getVideoIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}

	var input port.UpdateVideoData
	if err := request.DecodeJSONStrict(w, r, &input); err != nil {
		h.BadRequest(w, r, err)
		return
	}

	deps := usecase.UpdateVideoDeps{
		Accessor: h.videoAcc,
		Mutator:  h.videoMut,
	}
	video, err := usecase.UpdateVideo(deps).Invoke(r.Context(), port.UpdateVideoInput{
		ID:   videoID,
		Data: input,
	})
	if err != nil {
		if entity.IsError(err, entity.ErrVideoNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrVideoNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	h.MustSendJSON(w, r, http.StatusOK, response.Data{"video": video})
}

func (h *Video) HandleDelete(w http.ResponseWriter, r *http.Request) {
	videoID, ok := h.getVideoIDFromRequest(r)
	if !ok {
		h.BadRequest(w, r, errors.New("missing or invalid video id"))
		return
	}

	if _, err := h.videoAcc.ByID(r.Context(), videoID); err != nil {
		if entity.IsError(err, entity.ErrVideoNotFound) {
			h.ErrorMessage(w, r, http.StatusNotFound, entity.ErrVideoNotFound.Error(), nil)
			return
		}

		h.ServerError(w, r, err)
		return
	}

	if err := h.videoMut.Delete(r.Context(), videoID); err != nil {
		h.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Video) getVideoIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	videoIDRaw := chi.URLParam(r, "videoId")

	videoID, err := uuid.Parse(videoIDRaw)
	if err != nil {
		return uuid.UUID{}, false
	}

	return videoID, true
}

func (h *Video) getFindOptionsFromRequest(r *http.Request) (port.FindOptions, bool) {
	limit, err := strconv.ParseUint(h.DefaultQueryValue(r, "limit", strconv.FormatUint(_defaultLimit, 10)), 10, 64)
	if err != nil {
		return port.FindOptions{}, false
	}

	offset, err := strconv.ParseUint(h.DefaultQueryValue(r, "offset", strconv.FormatUint(_defaultOffset, 10)), 10, 64)
	if err != nil {
		return port.FindOptions{}, false
	}

	return port.FindOptions{Limit: limit, Offset: offset}, true
}

func (h *Video) getSearchQueryFromRequest(r *http.Request) (string, bool) {
	return h.GetQueryValue(r, "q")
}

func (h *Video) getAuthorNicknameFromRequest(r *http.Request) (string, bool) {
	return h.GetQueryValue(r, "author")
}

func (h *Video) defaultSortByFromRequest(r *http.Request, defaultValue string) string {
	return h.DefaultQueryValue(r, "sortBy", defaultValue)
}
