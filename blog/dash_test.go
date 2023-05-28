package dash

import (
   "fmt"
   "testing"
)

var test_reps = []representation{
   {width: 1920, height: 1080, bandwidth: 9_999_999},
   {width: 1920, height: 1080, bandwidth: 8_999_999},
   {width: 1920, height: 1080, bandwidth: 7_999_999},
   {width: 1920, height: 1080, bandwidth: 6_999_999},
   {width: 1920, height: 1080, bandwidth: 5_999_999},
   {width: 1920, height: 1080, bandwidth: 4_999_999},
   {width: 1920, height: 1080, bandwidth: 3_999_999},
   {width: 1920, height: 1080, bandwidth: 2_999_999},
   {width: 1920, height: 1080, bandwidth: 1_999_999},
   {width: 1280, height:  720, bandwidth: 1_099_999},
}

func Test_Dash(t *testing.T) {
   {
      index := index_bandwidth(test_reps, 2_000_000)
      fmt.Println(index)
   }
   {
      index := index_bandwidth(test_reps, 10_000_000)
      fmt.Println(index)
   }
}
