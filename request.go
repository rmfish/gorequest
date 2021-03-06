package gorequest

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/rmfish/logger"
	"go.uber.org/zap"
	"net/http"
)

var log = logger.GetLogger()

type RefererBuilder func(page interface{}) string
type CookieBuilder func(page interface{}) []*http.Cookie
type BodyBuilder func(page interface{}) interface{}
type ParamsBuilder func(page interface{}) interface{}
type UrlBuilder func(page interface{}) string
type OptBuilder func(page interface{}) *grequests.RequestOptions

type Request struct {
	method  string
	opt     OptBuilder
	Url     UrlBuilder
	Params  ParamsBuilder
	Referer RefererBuilder
	Cookie  CookieBuilder
	Body    BodyBuilder
}

type PagedRequest struct {
	*Request
}

func NewPagedRequest(method string, opt OptBuilder, url UrlBuilder, params ParamsBuilder, referer RefererBuilder, cookie CookieBuilder, body BodyBuilder) *PagedRequest {
	return &PagedRequest{&Request{method, opt, url, params, referer, cookie, body}}
}

func (req *PagedRequest) DoRequestReturnResult(page interface{}, result interface{}) error {
	resp, err := req.DoRequest(page)
	if err == nil {
		resp.JSON(result)
		return nil
	}
	return err
}

func (req *PagedRequest) DoRequest(page interface{}) (*grequests.Response, error) {
	if req.Url == nil {
		log.Error("Do get request failed. Missing url. ")
		return nil, errors.New("Do get request failed. Missing url. ")
	}

	url := req.Url(page)
	log.Debug(fmt.Sprintf("Do %s request.", req.method), zap.String("Url", url))
	opt := &grequests.RequestOptions{}
	if req.opt != nil {
		opt = req.opt(page)
	}
	if req.Cookie != nil {
		opt.Cookies = req.Cookie(page)
	}
	if req.Referer != nil {
		if opt.Headers == nil {
			opt.Headers = make(map[string]string)
		}
		opt.Headers["Referer"] = req.Referer(page)
	}
	if req.Body != nil {
		opt.JSON = req.Body(page)
	}
	switch req.method {
	case "GET":
		return grequests.Get(url, opt)
	case "POST":
		return grequests.Post(url, opt)
	case "PUT":
		return grequests.Put(url, opt)
	case "DELETE":
		return grequests.Delete(url, opt)
	case "PATCH":
		return grequests.Patch(url, opt)
	case "HEAD":
		return grequests.Head(url, opt)
	default:
		return grequests.Get(url, opt)
	}
}
