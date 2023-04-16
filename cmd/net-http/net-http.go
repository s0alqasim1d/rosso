package main

import (
   "2a.pages.dev/rosso/strconv"
   "bufio"
   "bytes"
   "embed"
   "fmt"
   "io"
   "net/http"
   "net/http/httputil"
   "net/textproto"
   "net/url"
   "os"
   "strings"
   "text/template"
   "unicode/utf8"
)

// go.dev/ref/spec#String_literals
func can_backquote(s string) bool {
   for i := range s {
      b := s[i]
      if b == '\r' {
         return false
      }
      if b == '`' {
         return false
      }
      if strconv.Binary_Data(b) {
         return false
      }
   }
   return utf8.ValidString(s)
}

func write_request(req *http.Request, dst io.Writer) error {
   var v values
   if req.Body != nil && req.Method != "GET" {
      body, err := io.ReadAll(req.Body)
      if err != nil {
         return err
      }
      text := string(body)
      if can_backquote(text) {
         v.Raw_Req_Body = "`" + text + "`"
      } else {
         v.Raw_Req_Body = strconv.Quote(text)
      }
      v.Req_Body = "io.NopCloser(req_body)"
      req.Body = io.NopCloser(bytes.NewReader(body))
   } else {
      v.Raw_Req_Body = `""`
      v.Req_Body = "nil"
   }
   v.Query = req.URL.Query()
   v.Request = req
   temp, err := template.ParseFS(content, "_template.go")
   if err != nil {
      return err
   }
   return temp.Execute(dst, v)
}

func write(req *http.Request, dst io.Writer) error {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if dst != nil {
      dump, err := httputil.DumpResponse(res, false)
      if err != nil {
         return err
      }
      os.Stdout.Write(dump)
      if _, err := io.Copy(dst, res.Body); err != nil {
         return err
      }
   } else {
      dump, err := httputil.DumpResponse(res, true)
      if err != nil {
         return err
      }
      enc := strconv.Encode(dump)
      if strings.HasSuffix(enc, "\n") {
         fmt.Print(enc)
      } else {
         fmt.Println(enc)
      }
   }
   return nil
}

func read_request(r *bufio.Reader) (*http.Request, error) {
   var req http.Request
   text := textproto.NewReader(r)
   // .Method
   raw_method_path, err := text.ReadLine()
   if err != nil {
      return nil, err
   }
   method_path := strings.Fields(raw_method_path)
   req.Method = method_path[0]
   // .URL
   ref, err := url.Parse(method_path[1])
   if err != nil {
      return nil, err
   }
   req.URL = ref
   // .URL.Host
   head, err := text.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   if req.URL.Host == "" {
      req.URL.Host = head.Get("Host")
   }
   // .Header
   req.Header = http.Header(head)
   // .Body
   buf := new(bytes.Buffer)
   length, err := text.R.WriteTo(buf)
   if err != nil {
      return nil, err
   }
   if length >= 1 {
      req.Body = io.NopCloser(buf)
   }
   // .ContentLength
   req.ContentLength = length
   return &req, nil
}

//go:embed _template.go
var content embed.FS

type values struct {
   *http.Request
   Query url.Values
   Req_Body string
   Raw_Req_Body string
}

type flags struct {
   golang bool
   https bool
   name string
   output string
}
