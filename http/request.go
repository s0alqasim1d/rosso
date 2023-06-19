package http

import (
   "bytes"
   "io"
   "net/http"
   "net/url"
   "strings"
)

type Request struct {
   *http.Request
}

func Get(ref *url.URL) *Request {
   return New_Request(http.MethodGet, ref)
}

func Get_Parse(ref string) (*Request, error) {
   href, err := url.Parse(ref)
   if err != nil {
      return nil, err
   }
   return New_Request(http.MethodGet, href), nil
}

func New_Request(method string, ref *url.URL) *Request {
   req := http.Request{
      Header: make(http.Header),
      Method: method,
      ProtoMajor: 1,
      ProtoMinor: 1,
      URL: ref,
   }
   return &Request{&req}
}

func Patch(ref *url.URL) *Request {
   return New_Request(http.MethodPatch, ref)
}

func Post(ref *url.URL) *Request {
   return New_Request(http.MethodPost, ref)
}

func Post_Parse(ref string) (*Request, error) {
   href, err := url.Parse(ref)
   if err != nil {
      return nil, err
   }
   return New_Request(http.MethodPost, href), nil
}

func (r Request) Body_Bytes(b []byte) {
   body := bytes.NewReader(b)
   r.Body = io.NopCloser(body)
   r.ContentLength = body.Size()
}

func (r Request) Body_String(s string) {
   body := strings.NewReader(s)
   r.Body = io.NopCloser(body)
   r.ContentLength = body.Size()
}
