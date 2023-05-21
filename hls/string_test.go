package hls

import (
   "bytes"
   "fmt"
   "os"
   "sort"
   "testing"
)

var segment_tests = []string{
   "m3u8/cbc-audio.m3u8",
   "m3u8/cbc-video.m3u8",
   "m3u8/nbc-segment.m3u8",
   "m3u8/roku-segment.m3u8",
}

func Test_Reverse(t *testing.T) {
   for _, test := range segment_tests {
      fmt.Println(test + ":")
      data, err := os.ReadFile(test + ".txt")
      if err != nil {
         t.Fatal(err)
      }
      sort.SliceStable(data, func(int, int) bool {
         return true
      })
      os.Stdout.Write(data)
   }
}

func Test_Segment(t *testing.T) {
   for _, test := range segment_tests {
      data, err := os.ReadFile(test + ".txt")
      if err != nil {
         t.Fatal(err)
      }
      sort.SliceStable(data, func(int, int) bool {
         return true
      })
      seg, err := New_Scanner(bytes.NewReader(data)).Segment()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n\n", seg)
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
      for _, item := range master.Streams {
         fmt.Println(item)
      }
      for _, item := range master.Media {
         fmt.Println(item)
      }
      fmt.Println()
   }
}
