package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "strings"
   "testing"
)

func Test_Audio(t *testing.T) {
   for _, name := range tests {
      data, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var pre Presentation
      if err := xml.Unmarshal(data, &pre); err != nil {
         t.Fatal(err)
      }
      reps := Audio(pre.Representation())
      target := reps.Index(func(carry, item Representation) bool {
         if !strings.HasPrefix(item.Adaptation.Lang, "en") {
            return false
         }
         if !strings.Contains(item.Codecs, "mp4a.") {
            return false
         }
         if item.Role() == "description" {
            return false
         }
         return true
      })
      fmt.Println(name)
      for i, rep := range reps {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}
var test_reps = []Representation{
   {Width: 1920, Height: 1080, Bandwidth: 9_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 8_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 7_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 6_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 5_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 4_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 3_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 2_999_999},
   {Width: 1920, Height: 1080, Bandwidth: 1_999_999},
   {Width: 1280, Height:  720, Bandwidth: 1_099_999},
}

func Test_Dash(t *testing.T) {
   {
      index := Index_Bandwidth(test_reps, 2_000_000)
      fmt.Println(index)
   }
   {
      index := Index_Bandwidth(test_reps, 10_000_000)
      fmt.Println(index)
   }
}

func Test_Video(t *testing.T) {
   for _, name := range tests {
      data, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var pre Presentation
      if err := xml.Unmarshal(data, &pre); err != nil {
         t.Fatal(err)
      }
      reps := Video(pre.Representation())
      fmt.Println(name)
      for i, rep := range reps {
         if i == reps.Bandwidth(0) {
            fmt.Print("!")
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Info(t *testing.T) {
   for _, name := range tests {
      data, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var pre Presentation
      if err := xml.Unmarshal(data, &pre); err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      reps := pre.Representation()
      
      for _, rep := range Audio(reps) {
         fmt.Println(rep)
      }
      for _, rep := range Video(reps) {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

