package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError.Error())
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}
func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}
func ErrBadRequestCustomMessage(msg string) Error {
	return NewError(http.StatusBadRequest, msg)
}
func ErrResourceNotFound(msg string) Error {
	return NewError(http.StatusNotFound, msg)
}
func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "bad request")
}
