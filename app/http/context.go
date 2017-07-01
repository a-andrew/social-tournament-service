package http

import (
	"github.com/julienschmidt/httprouter"
	"net/url"
	"io"
	"encoding/json"
)

type Context struct{
	params httprouter.Params
	values url.Values
	body io.ReadCloser
}

func (c *Context) GetRoutParams() httprouter.Params{
	return c.params
}

func (c *Context) GetUrlParams() url.Values{
	return c.values
}

func (c *Context) ParseBody(v interface{}){
	defer c.body.Close()
	json.NewDecoder(c.body).Decode(v)
}
