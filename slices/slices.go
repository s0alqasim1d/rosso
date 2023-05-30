package slices

import "sort"

// github.com/golang/go/blob/go1.20.4/src/internal/types/testdata/check/slices.go
func Filter[T any](s []T, f func(T) bool) []T {
   var values []T
   for _, value := range s {
      if f(value) {
         values = append(values, value)
      }
   }
   return values
}

// github.com/golang/exp/blob/2e198f4/slices/slices.go
func Index_Func[T any](s []T, f func(T) bool) int {
   for i, value := range s {
      if f(value) {
         return i
      }
   }
   return -1
}

// github.com/golang/exp/blob/2e198f4/slices/sort.go
func Sort_Func[T any](s []T, less func(a, b T) bool) {
   sort.Slice(s, func(i, j int) bool {
      return less(s[i], s[j])
   })
}
