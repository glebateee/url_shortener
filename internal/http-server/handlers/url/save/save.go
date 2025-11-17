package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "urlshortener/internal/lib/api/response"
	"urlshortener/internal/lib/logger/sl"
	"urlshortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"urlshortener/internal/lib/random"

	"github.com/go-playground/validator/v10"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(targetUrl, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)

		var req Request

		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(writer, request, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(writer, request, resp.ValidationError(validationErrs))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewrandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(writer, request, resp.Error(storage.ErrURLExists.Error()))
				return
			}
			log.Error("failed to add url", sl.Err(err))
			render.JSON(writer, request, resp.Error("failed to add url"))
			return
		}
		log.Info("url added", slog.Int64("id", id))
		sendOK(writer, request, alias)
	}
}

func sendOK(writer http.ResponseWriter, request *http.Request, alias string) {
	render.JSON(writer, request, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
