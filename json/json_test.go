package json

import (
   "fmt"
   "testing"
)

const text = `hello world {"year":12,"month":31}`

type date struct {
   Year int
   Month int
}

func Test_Cut(t *testing.T) {
   data, sep := []byte(text), []byte(" world ")
   var v date
   err := Cut(data, sep, &v)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", v)
}

func Test_Before(t *testing.T) {
   data, sep := []byte(text), []byte(`{"year"`)
   var v date
   err := Cut_Before(data, sep, &v)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", v)
}
