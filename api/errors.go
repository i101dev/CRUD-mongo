package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {

	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}

	apiError := NewError(http.StatusInternalServerError, err.Error())

	return c.Status(apiError.Code).JSON(apiError)
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Msg:  msg,
	}
}

func ErrorInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "invalid ID provided",
	}
}
func ErrorUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Msg:  "unauthorized request",
	}
}

func ErrorBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "invalid JSON request",
	}
}

func ErrorResourceNotFound(rsc string) Error {
	return Error{
		Code: http.StatusNotFound,
		Msg:  "resource not found:" + rsc,
	}
}
