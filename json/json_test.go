package json

import (
   "fmt"
   "testing"
)

const dirty = `hello world {"year":12,"month":31}`

func Test_Cut(t *testing.T) {
   text, sep := []byte(dirty), []byte(" world ")
   var value date
   err := Cut(text, sep, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}

type date struct {
   Year int
   Month int
}

func Test_Before(t *testing.T) {
   text, sep := []byte(dirty), []byte(`{"year"`)
   var value date
   err := Cut_Before(text, sep, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}
