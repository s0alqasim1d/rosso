package slices

import "sort"

// github.com/golang/go/blob/0df6812/src/slices/slices.go 
func Clone[E any](s []E) []E {
   return append([]E{}, s...)
}

// github.com/golang/go/blob/0df6812/src/slices/slices.go 
func Delete[E any](s []E, del func(E) bool) []E {
   for i, v := range s {
      if del(v) {
         j := i
         for i++; i < len(s); i++ {
            v = s[i]
            if !del(v) {
               s[j] = v
               j++
            }
         }
         return s[:j]
      }
   }
   return s
}

// github.com/golang/go/blob/0df6812/src/slices/slices.go 
func Index[E any](s []E, f func(E) bool) int {
   for i, v := range s {
      if f(v) {
         return i
      }
   }
   return -1
}

// godocs.io/sort#Slice
func Sort[E any](x []E, less func(a, b E) bool) {
   sort.Slice(x, func(i, j int) bool {
      return less(x[i], x[j])
   })
}
