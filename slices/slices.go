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
   for i := range s {
      if f(s[i]) {
         return i
      }
   }
   return -1
}

// github.com/golang/go/blob/0df6812/src/slices/sort.go
func Sort[E any](x []E, cmp func(a, b E) int) {
   sort.Slice(x, func(i, j int) bool {
      return cmp(x[i], x[j]) >= 1
   })
}
