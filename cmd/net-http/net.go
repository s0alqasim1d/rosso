package main

import (
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
   "strconv"
   "strings"
   "text/template"
)

func write(req *http.Request, file io.Writer) error {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if file != nil {
      dump, err := httputil.DumpResponse(res, false)
      if err != nil {
         return err
      }
      os.Stdout.Write(dump)
      if _, err := io.Copy(file, res.Body); err != nil {
         return err
      }
   } else {
      dump, err := httputil.DumpResponse(res, true)
      if err != nil {
         return err
      }
      str := string(dump)
      if !strconv.CanBackquote(str) {
         str = strconv.Quote(str)
      }
      if strings.HasSuffix(str, "\n") {
         fmt.Print(str)
      } else {
         fmt.Println(str)
      }
   }
   return nil
}

func write_request(req *http.Request, dst io.Writer) error {
   var v values
   if req.Body != nil && req.Method != "GET" {
      body, err := io.ReadAll(req.Body)
      if err != nil {
         return err
      }
      req.Body = io.NopCloser(bytes.NewReader(body))
      v.Req_Body = "io.NopCloser(req_body)"
      str := string(body)
      if strconv.CanBackquote(str) {
         v.Raw_Req_Body = "`" + str + "`"
      } else {
         v.Raw_Req_Body = strconv.Quote(str)
      }
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
