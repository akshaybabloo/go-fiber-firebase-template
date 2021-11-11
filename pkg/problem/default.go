package problem

import (
	"net/http"

	"github.com/akshaybabloo/go-fiber-template/model"
)

type Problems interface {
	Problem(detail, instance string, status int, title, errorType string) *model.ProblemDetails
	InternalServerErrorProblem(title, detail, instance string) *model.ProblemDetails
	PageNotFoundProblem(instance string) *model.ProblemDetails
}

func NewProblems() Problems {
	return &problemsStubs{}
}

type problemsStubs struct {
	Problems
}

func (p *problemsStubs) Problem(detail, instance string, status int, title, errorType string) *model.ProblemDetails {
	return &model.ProblemDetails{
		Type:     errorType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

func (p *problemsStubs) InternalServerErrorProblem(title, detail, instance string) *model.ProblemDetails {
	return &model.ProblemDetails{
		Title:    title,
		Instance: instance,
		Type:     "https://tools.ietf.org/html/rfc7231#section-6.6.1",
		Status:   http.StatusInternalServerError,
		Detail:   detail,
	}
}

func (p *problemsStubs) PageNotFoundProblem(instance string) *model.ProblemDetails {
	return &model.ProblemDetails{
		Title:    "page not found",
		Instance: instance,
		Type:     "https://tools.ietf.org/html/rfc7231#section-6.5.4",
		Status:   http.StatusNotFound,
		Detail:   "the page you are looking for does not exists",
	}
}
