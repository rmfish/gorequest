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

type PagedRequest struct {
	method  string
	opt     OptBuilder
	url     UrlBuilder
	params  ParamsBuilder
	referer RefererBuilder
	cookie  CookieBuilder
	body    BodyBuilder
}

func New(method string, opt OptBuilder, url UrlBuilder, params ParamsBuilder, referer RefererBuilder, cookie CookieBuilder, body BodyBuilder) *PagedRequest {
	return &PagedRequest{method, opt, url, params, referer, cookie, body}
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
	if req.url == nil {
		log.Error("Do get request failed. Missing url. ")
		return nil, errors.New("Do get request failed. Missing url. ")
	}

	url := req.url(page)
	log.Debug(fmt.Sprintf("Do %s request.", req.method), zap.String("Url", url))
	opt := &grequests.RequestOptions{}
	if req.opt != nil {
		opt = req.opt(page)
	}
	if req.cookie != nil {
		opt.Cookies = req.cookie(page)
	}
	if req.referer != nil {
		if opt.Headers == nil {
			opt.Headers = make(map[string]string)
		}
		opt.Headers["Referer"] = req.referer(page)
	}
	if req.body != nil {
		opt.JSON = req.body(page)
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
