package problem

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type ErrOptions struct {
	Status  int
	Err     error
	Kind    string
	Request *http.Request
}

type OptionsFunc func(*ProblemDetail)

/*
ProblemDetail
https://tools.ietf.org/html/rfc7807
https://datatracker.ietf.org/doc/rfc9457/

HTTP/1.1 403 Forbidden
Content-Type: application/problem+json
Content-Language: en

	{
	   "type": "https://example.com/probs/out-of-credit",
	   "title": "You do not have enough credit.",
	   "detail": "Your current balance is 30, but that costs 50.",
	   "instance": "/account/12345/msgs/abc",
	   "balance": 30,
	   "accounts": ["/account/12345", "/account/67890"]
	}
*/
type ProblemDetail struct {
	Status   int    `json:"status,omitempty" example:"401"`
	Title    string `json:"title,omitempty" example:"Unauthorized"`
	Detail   string `json:"detail,omitempty" example:"invalid credentials"`
	Type     string `json:"type,omitempty" example:"Err"`
	Instance string `json:"instance,omitempty" example:"/api/example"`
	Extension
}

type Extension struct {
	RequestId  string    `json:"requestId,omitempty" example:"123"`
	Timestamp  time.Time `json:"timestamp" example:"2024-12-11T20:23:55.248212-03:00"`
	StackTrace string    `json:"stackTrace,omitempty"`
}

// Error implements the error interface
func (p ProblemDetail) Error() string {
	return p.Title
}

func NewProblem(opt ErrOptions) *ProblemDetail {
	problem := &ProblemDetail{
		Extension: Extension{
			Timestamp: time.Now(),
		},
	}

	if opt.Request != nil {
		requestIdCtx := opt.Request.Context().Value(middleware.RequestIDKey)
		requestId := ""
		if requestIdCtx != nil {
			requestId = requestIdCtx.(string)
		}
		problem.Instance = opt.Request.RequestURI
		problem.Extension.RequestId = requestId
	}

	problem.Status = opt.Status
	problem.Detail = opt.Err.Error()
	problem.Title = opt.Kind

	return problem
}
