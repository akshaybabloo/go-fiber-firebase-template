package problem

import (
	"net/http"

	"github.com/akshaybabloo/go-fiber-template/model"
)

// MIMEApplicationProblemJSON is the media type for RFC 9457 problem details.
const MIMEApplicationProblemJSON = "application/problem+json"

// typeBlank is used when the problem has no semantics beyond the HTTP status
// code, as recommended by RFC 9457 §4.2.1.
const typeBlank = "about:blank"

// Problems builds RFC 9457 problem detail responses.
type Problems interface {
	Problem(detail, instance string, status int, title, errorType string) *model.ProblemDetails
	InternalServerErrorProblem(title, detail, instance string) *model.ProblemDetails
	UnauthorizedProblem(title, detail, instance string) *model.ProblemDetails
	PageNotFoundProblem(instance string) *model.ProblemDetails
}

func NewProblems() Problems {
	return &problems{}
}

type problems struct{}

func (p *problems) Problem(detail, instance string, status int, title, errorType string) *model.ProblemDetails {
	if errorType == "" {
		errorType = typeBlank
	}
	return &model.ProblemDetails{
		Type:     errorType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

func (p *problems) InternalServerErrorProblem(title, detail, instance string) *model.ProblemDetails {
	return &model.ProblemDetails{
		Type:     typeBlank,
		Title:    title,
		Status:   http.StatusInternalServerError,
		Detail:   detail,
		Instance: instance,
	}
}

func (p *problems) UnauthorizedProblem(title, detail, instance string) *model.ProblemDetails {
	return &model.ProblemDetails{
		Type:     typeBlank,
		Title:    title,
		Status:   http.StatusUnauthorized,
		Detail:   detail,
		Instance: instance,
	}
}

func (p *problems) PageNotFoundProblem(instance string) *model.ProblemDetails {
	return &model.ProblemDetails{
		Type:     typeBlank,
		Title:    "page not found",
		Status:   http.StatusNotFound,
		Detail:   "the page you are looking for does not exist",
		Instance: instance,
	}
}
