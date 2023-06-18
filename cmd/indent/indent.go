package main

import (
   "2a.pages.dev/rosso/xml"
   "encoding/json"
   "flag"
   "os"
)

func (f flags) indent_xml() error {
   // in
   in, err := os.Open(f.input)
   if err != nil {
      return err
   }
   defer in.Close()
   // out
   out := os.Stdout
   if f.output != "" {
      out, err = os.Create(f.output)
      if err != nil {
         return err
      }
      defer out.Close()
   }
   // Indent
   return xml.Indent(out, in, "", " ")
}

func (f flags) indent_json() error {
   file := os.Stdout
   if f.output != "" {
      var err error
      file, err = os.Create(f.output)
      if err != nil {
         return err
      }
   }
   defer file.Close()
   var value any
   {
      b, err := os.ReadFile(f.input)
      if err != nil {
         return err
      }
      if err := json.Unmarshal(b, &value); err != nil {
         return err
      }
   }
   enc := json.NewEncoder(file)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   return enc.Encode(value)
}

type flags struct {
   input string
   output string
   xml bool
}

func main() {
   var f flags
   flag.StringVar(&f.input, "f", "", "input file")
   flag.StringVar(&f.output, "o", "", "output file")
   flag.BoolVar(&f.xml, "xml", false, "use XML instead of JSON")
   flag.Parse()
   if f.input != "" {
      var err error
      if f.xml {
         err = f.indent_xml()
      } else {
         err = f.indent_json()
      }
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

