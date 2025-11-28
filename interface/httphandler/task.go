package httphandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kanaru-ssk/go-rpc-server/domain/task"
	"github.com/kanaru-ssk/go-rpc-server/interface/httpresponse"
	"github.com/kanaru-ssk/go-rpc-server/interface/response/errorresponse"
	"github.com/kanaru-ssk/go-rpc-server/interface/response/taskresponse"
	"github.com/kanaru-ssk/go-rpc-server/usecase"
)

type TaskHandler struct {
	taskUseCase *usecase.TaskUseCase
	taskMapper  *taskresponse.Mapper
	errorMapper *errorresponse.Mapper
}

func NewTaskHandler(
	taskUseCase *usecase.TaskUseCase,
	taskMapper *taskresponse.Mapper,
	errorMapper *errorresponse.Mapper,
) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
		taskMapper:  taskMapper,
		errorMapper: errorMapper,
	}
}

func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", errorresponse.ErrMethodNotAllowed)
		httpresponse.RenderJson(ctx, w, http.StatusMethodNotAllowed, h.errorMapper.MapErrorResponse(errorresponse.ErrMethodNotAllowed))
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	t, err := h.taskUseCase.Get(ctx, req.ID)
	switch {
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapGetResponse(t))

	case errors.Is(err, task.ErrInvalidID):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	case errors.Is(err, task.ErrNotFound):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusNotFound, h.errorMapper.MapErrorResponse(errorresponse.ErrNotFound))

	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

func (h *TaskHandler) HandleListV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleListV1", "err", errorresponse.ErrMethodNotAllowed)
		httpresponse.RenderJson(ctx, w, http.StatusMethodNotAllowed, h.errorMapper.MapErrorResponse(errorresponse.ErrMethodNotAllowed))
		return
	}

	t, err := h.taskUseCase.List(ctx)
	switch {
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapListResponse(t))

	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleListV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

func (h *TaskHandler) HandleCreateV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		httpresponse.RenderJson(ctx, w, http.StatusMethodNotAllowed, h.errorMapper.MapErrorResponse(errorresponse.ErrMethodNotAllowed))
		return
	}

	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleCreateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	t, err := h.taskUseCase.Create(ctx, req.Title)
	switch {
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapCreateResponse(t))

	case errors.Is(err, task.ErrInvalidTitle):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleCreateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleCreateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

func (h *TaskHandler) HandleUpdateV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", errorresponse.ErrMethodNotAllowed)
		httpresponse.RenderJson(ctx, w, http.StatusMethodNotAllowed, h.errorMapper.MapErrorResponse(errorresponse.ErrMethodNotAllowed))
		return
	}

	var req struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	t, err := h.taskUseCase.Update(ctx, req.ID, req.Title, req.Status)
	switch {
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapUpdateResponse(t))

	case errors.Is(err, task.ErrInvalidID),
		errors.Is(err, task.ErrInvalidTitle),
		errors.Is(err, task.ErrInvalidStatus):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

func (h *TaskHandler) HandleDeleteV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDeleteV1", "err", errorresponse.ErrMethodNotAllowed)
		httpresponse.RenderJson(ctx, w, http.StatusMethodNotAllowed, h.errorMapper.MapErrorResponse(errorresponse.ErrMethodNotAllowed))
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDeleteV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	err := h.taskUseCase.Delete(ctx, req.ID)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)

	case errors.Is(err, task.ErrInvalidID):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleDeleteV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}
