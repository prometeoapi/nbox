package response

import (
	"encoding/json"
	"errors"
	"nbox/internal/entrypoints/api/problem"
	"net/http"
)

type Empty struct {
}

type Envelop struct {
	Body    *interface{}
	Problem *problem.ProblemDetail
}

type EnvelopFunc func(*Envelop)

func EnvelopWithErr(err problem.ErrOptions) EnvelopFunc {
	return func(envelop *Envelop) {
		envelop.Body = nil
		envelop.Problem = problem.NewProblem(err)
	}
}

func EnvelopWithBody(body interface{}) EnvelopFunc {
	return func(envelop *Envelop) {
		if body == nil {
			body = Empty{}
		}
		envelop.Body = &body
		envelop.Problem = nil
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, r, errors.New("404 page not found"), http.StatusNotFound)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, r, errors.New("method not allowed"), http.StatusMethodNotAllowed)
}

func Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	Json(w, r, EnvelopWithErr(
		problem.ErrOptions{
			Status:  code,
			Err:     err,
			Kind:    "Err",
			Request: r,
		}),
	)
}

func Success(w http.ResponseWriter, r *http.Request, body interface{}) {
	Json(w, r, EnvelopWithBody(body))
}

func Json(w http.ResponseWriter, r *http.Request, optFns ...EnvelopFunc) {
	w.Header().Set("Content-Type", "application/json")
	var out []byte
	envelop := &Envelop{}
	//var err error

	for _, optFn := range optFns {
		optFn(envelop)
	}

	if envelop.Body == nil {
		out, _ = json.Marshal(envelop.Problem)
		w.WriteHeader(envelop.Problem.Status)
		_, _ = w.Write(out)
		return
	}

	out, _ = json.Marshal(envelop.Body)
	_, _ = w.Write(out)
	return
}
