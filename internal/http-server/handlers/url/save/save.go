package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "urlshortener/internal/lib/logger/api/response"
	"urlshortener/internal/lib/logger/sl"
	"urlshortener/internal/lib/random"
	"urlshortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type UrlSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)

		var reqData Request

		err := render.DecodeJSON(request.Body, &reqData)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(writer, request, resp.Error("failed to decode request"))
			return
		}
		slog.Info("request body decoded", slog.Any("request", reqData))

		if err := validator.New().Struct(reqData); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(writer, request, resp.ValidationError(validationErrs))
			return
		}

		alias := reqData.Alias
		if alias == "" {
			// TODO : handle collisions
			alias = random.NewrandomString(aliasLength)
		}
		_, err = urlSaver.SaveURL(reqData.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", reqData.URL))
				render.JSON(writer, request, resp.Error("url already exists"))
				return
			}
			log.Error("failed to add url", sl.Err(err))
			render.JSON(writer, request, resp.Error("failed to add url"))
			return
		}
	}
}
