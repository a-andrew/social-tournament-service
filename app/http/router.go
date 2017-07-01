package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"fmt"
)

type Router struct{
	r *httprouter.Router
}

type handlerFunc func(*Context) (int, interface{})

func NewRouter() *Router{
	return &Router{
		r: httprouter.New(),
	}
}

func (r *Router) Get(rout string, handler handlerFunc, requiredUrlParams ...string){
	r.r.GET(rout, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){
		urlQuery := req.URL.Query()
		for _, p := range requiredUrlParams{
			if v := urlQuery.Get(p); v == ""{
				r.sendResponse(w, http.StatusBadRequest, ResError{
					Msg: fmt.Sprintf("`%s` is required", p),
				})
				
				return
			}
		}
		
		statusCode, res := handler(&Context{
			params: ps,
			values: urlQuery,
		})
		
		r.sendResponse(w, statusCode, res)
	})
}

func (r *Router) Post(rout string, handler handlerFunc){
	r.r.POST(rout, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){
		statusCode, res := handler(&Context{
			params: ps,
			body: req.Body,
		})
		
		r.sendResponse(w, statusCode, res)
	})
}

func (r *Router) sendResponse(w http.ResponseWriter, code int, res interface{}){
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

func (r *Router) Serve() error{
	return http.ListenAndServe(":8080", r.r)
}