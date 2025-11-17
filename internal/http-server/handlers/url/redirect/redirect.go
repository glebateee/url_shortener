package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	resp "urlshortener/internal/lib/api/response"
	"urlshortener/internal/lib/logger/sl"
	"urlshortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)

		alias := chi.URLParam(request, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(writer, request, resp.Error("invalid request"))
			return
		}

		resUrl, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", alias)
				render.JSON(writer, request, resp.Error(storage.ErrURLNotFound.Error()))
				return
			}
			log.Error("failed to get url", sl.Err(err))
			render.JSON(writer, request, resp.Error("internal error"))
			return
		}
		http.Redirect(writer, request, resUrl, http.StatusFound)
	}
}
