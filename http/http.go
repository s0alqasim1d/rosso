package http

import (
   "2a.pages.dev/rosso/strconv"
   "bytes"
   "fmt"
   "io"
   "net/http"
   "net/http/httputil"
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
const StatusFound = http.StatusFound

type Client struct {
   Log_Level int // this needs to work with flag.IntVar
   Status int
   http.Client
}

var Default_Client = Client{
   Client: http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   },
   Log_Level: 1,
   Status: http.StatusOK,
}

func (c Client) Do(req *Request) (*Response, error) {
   switch c.Log_Level {
   case 1:
      fmt.Println(req.Method, req.URL)
   case 2:
      dump, err := httputil.DumpRequest(req.Request, true)
      if err != nil {
         return nil, err
      }
      enc := strconv.Encode(dump)
      if strings.HasSuffix(enc, "\n") {
         fmt.Print(enc)
      } else {
         fmt.Println(enc)
      }
   }
   res, err := c.Client.Do(req.Request)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != c.Status {
      return nil, fmt.Errorf(res.Status)
   }
   return res, nil
}

func (c Client) Get(ref string) (*Response, error) {
   req, err := Get_Parse(ref)
   if err != nil {
      return nil, err
   }
   return c.Do(req)
}

type Cookie = http.Cookie

type Header = http.Header

type Response = http.Response
