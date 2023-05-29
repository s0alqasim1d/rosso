package hls

import (
   "bytes"
   "fmt"
   "os"
   "sort"
   "strings"
   "testing"
)

func Test_Stream(t *testing.T) {
   for key, value := range master_tests {
      file, err := os.Open(key)
      if err != nil {
         t.Fatal(err)
      }
      master, err := New_Scanner(file).Master()
      if err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      items := Filter(master.Stream, value.stream)
      index := Index_Func(items, func(s Stream) bool {
         return s.Bandwidth >= 1
      })
      fmt.Println(key)
      for i, item := range items {
         if i == index {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      fmt.Println()
   }
}

var master_tests = map[string]filters{
   "m3u8/nbc-master.m3u8": {nbc_media, nil},
   "m3u8/roku-master.m3u8": {nil, nil},
   "m3u8/cbc-master.m3u8": {cbc_media, cbc_stream},
}

func Test_Media(t *testing.T) {
   for key, value := range master_tests {
      file, err := os.Open(key)
      if err != nil {
         t.Fatal(err)
      }
      master, err := New_Scanner(file).Master()
      if err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      master.Media = Filter(master.Media, value.media)
      target := Index_Func(master.Media, func(m Media) bool {
         return m.Name == "English"
      })
      fmt.Println(key)
      for i, media := range master.Media {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(media)
      }
      fmt.Println()
   }
}

func Test_Info(t *testing.T) {
   for test := range master_tests {
      data, err := os.ReadFile(test + ".txt")
      if err != nil {
         t.Fatal(err)
      }
      sort.SliceStable(data, func(int, int) bool {
         return true
      })
      master, err := New_Scanner(bytes.NewReader(data)).Master()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(test)
      for _, item := range master.Stream {
         fmt.Println(item)
      }
      for _, item := range master.Media {
         fmt.Println(item)
      }
      fmt.Println()
   }
}

func cbc_media(m Media) bool {
   return m.Type == "AUDIO"
}

func cbc_stream(s Stream) bool {
   return strings.Contains(s.Codecs, "avc1.")
}

func nbc_media(m Media) bool {
   return m.Type == "AUDIO"
}

type filters struct {
   media func(Media) bool
   stream func(Stream) bool
}
