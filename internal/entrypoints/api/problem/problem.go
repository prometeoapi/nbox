package problem

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
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
	Status   int    `json:"status,omitempty"`
	Title    string `json:"title,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Type     string `json:"type,omitempty"`
	Instance string `json:"instance,omitempty"`
	Extension
}

type Extension struct {
	RequestId  string    `json:"requestId,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
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
