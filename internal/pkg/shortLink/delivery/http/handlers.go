package http

import (
	"errors"
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/models"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"net/http"
	"net/url"
)

type Handler struct {
	useCase shortLink.Usecase
	log     logger.Logger
}

func NewHandler(useCase shortLink.Usecase, log logger.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		log:     log,
	}
}

func handleError(ctx *gin.Context, log logger.Logger, err error, statusCode int) {
	var errorResponse models.ErrorResponse

	if statusCode >= 500 {
		errorResponse = models.ErrorResponse{Error: "Internal Server Error"}
	} else {
		errorResponse = models.ErrorResponse{Error: err.Error()}
	}

	errorJSON, jsonErr := easyjson.Marshal(errorResponse)
	if jsonErr != nil {
		log.Errorf("Error during JSON marshaling: %v", jsonErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	logFunc := log.Info
	if statusCode >= 500 {
		logFunc = log.Error
	}

	logFunc("Error: %v", err)
	ctx.Data(statusCode, "application/json", errorJSON)
}

func (h *Handler) CreateLink(ctx *gin.Context) {
	var longUrl models.Request
	if err := easyjson.UnmarshalFromReader(ctx.Request.Body, &longUrl); err != nil {
		handleError(ctx, h.log, err, http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(longUrl.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		handleError(ctx, h.log, errors.New("bad url"), http.StatusBadRequest)
		return
	}

	link, err := h.useCase.CreateLink(parsedURL.Host)
	if err != nil {
		handleError(ctx, h.log, errors.New("bad url"), http.StatusBadRequest)
		return
	}

	data, err := easyjson.Marshal(models.Response{URL: link})
	if err != nil {
		handleError(ctx, h.log, err, http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, "application/json", data)
}

func (h *Handler) GetLink(ctx *gin.Context) {
	shortURL := ctx.Param("url")

	link, err := h.useCase.GetUrl(shortURL)

	var noSuchLinkErr *models.NoSuchLink
	if err != nil {
		if errors.As(err, &noSuchLinkErr) {
			handleError(ctx, h.log, errors.New("such link does not exist or it has expired"), http.StatusNotFound)
			return
		} else {
			handleError(ctx, h.log, err, http.StatusInternalServerError)
			return
		}
	}

	data, err := easyjson.Marshal(models.Response{URL: link})

	if err != nil {
		handleError(ctx, h.log, err, http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, "application/json", data)
}
