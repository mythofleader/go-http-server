// Package errors provides error classes for HTTP status codes.
package errors

import "net/http"

type HTTPError interface {
	error
	StatusCode() int // 상태 코드를 반환
}

type BadRequestHttpError struct {
	Message string
}

func (e *BadRequestHttpError) Error() string {
	return e.Message
}

func (e *BadRequestHttpError) StatusCode() int {
	return http.StatusBadRequest
}

func NewBadRequestHttpError(err error) *BadRequestHttpError {
	return &BadRequestHttpError{
		Message: err.Error(),
	}
}

type UnauthorizedHttpError struct {
	Message string
}

func (e *UnauthorizedHttpError) Error() string {
	return e.Message
}

func (e *UnauthorizedHttpError) StatusCode() int {
	return http.StatusUnauthorized
}

func NewUnauthorizedHttpError(err error) *UnauthorizedHttpError {
	return &UnauthorizedHttpError{
		Message: err.Error(),
	}
}

type ForbiddenHttpError struct {
	Message string
}

func (e *ForbiddenHttpError) Error() string {
	return e.Message
}

func (e *ForbiddenHttpError) StatusCode() int {
	return http.StatusForbidden
}

func NewForbiddenHttpError(err error) *ForbiddenHttpError {
	return &ForbiddenHttpError{
		Message: err.Error(),
	}
}

type NotFoundHttpError struct {
	Message string
}

func (e *NotFoundHttpError) Error() string {
	return e.Message
}

func (e *NotFoundHttpError) StatusCode() int {
	return http.StatusNotFound
}

func NewNotFoundHttpError(err error) *NotFoundHttpError {
	return &NotFoundHttpError{
		Message: err.Error(),
	}
}

type InternalServerHttpError struct {
	Message string
}

func (e *InternalServerHttpError) Error() string {
	return e.Message
}

func (e *InternalServerHttpError) StatusCode() int {
	return http.StatusInternalServerError
}

func NewInternalServerHttpError(err error) *InternalServerHttpError {
	return &InternalServerHttpError{
		Message: err.Error(),
	}
}

type ServiceUnavailableHttpError struct {
	Message string
}

func (e *ServiceUnavailableHttpError) Error() string {
	return e.Message
}

func (e *ServiceUnavailableHttpError) StatusCode() int {
	return http.StatusServiceUnavailable
}

func NewServiceUnavailableHttpError(err error) *ServiceUnavailableHttpError {
	return &ServiceUnavailableHttpError{
		Message: err.Error(),
	}
}

type MethodNotAllowedHttpError struct {
	Message string
}

func (e *MethodNotAllowedHttpError) Error() string {
	return e.Message
}

func (e *MethodNotAllowedHttpError) StatusCode() int {
	return http.StatusMethodNotAllowed
}

func NewMethodNotAllowedHttpError(err error) *MethodNotAllowedHttpError {
	return &MethodNotAllowedHttpError{
		Message: err.Error(),
	}
}
