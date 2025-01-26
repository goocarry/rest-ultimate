package user

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/goocarry/rest-ultimate/internal/storage"

	"github.com/goocarry/rest-ultimate/internal/lib/api/response"
	"github.com/goocarry/rest-ultimate/internal/lib/logger/sl"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	TgUserId string `json:"tg_user_id" validate:"required"`
}

type Response struct {
	response.Response
}

func New(log *slog.Logger, userSaver storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"

		l := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			l.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			l.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		l.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			l.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		user := storage.User{
			TgUserId: req.TgUserId,
		}
		id, err := userSaver.User().RegisterUser(user)
		if err != nil {
			l.Error("failed to save user", sl.Err(err))

			render.JSON(w, r, response.Error("failed to save user"))

			return
		}

		l.Info("user saved", slog.Int64("id", id))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}
