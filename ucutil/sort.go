package ucutil

import "sort"

// sort int64
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func SortInt64(a []int64)               { sort.Sort(Int64Slice(a)) }
func RevSortInt64(a []int64)            { sort.Sort(sort.Reverse(Int64Slice(a))) }

// sort uint64
type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func SortUint64(a []uint64)              { sort.Sort(Uint64Slice(a)) }
func RevSortUint64(a []uint64)           { sort.Sort(sort.Reverse(Uint64Slice(a))) }

// sort uint32
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func SortUint32(a []uint32)              { sort.Sort(Uint32Slice(a)) }
func RevSortUint32(a []uint32)           { sort.Sort(sort.Reverse(Uint32Slice(a))) }

// sort string
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func SortString(a []string)              { sort.Sort(StringSlice(a)) }
func RevSortString(a []string)           { sort.Sort(sort.Reverse(StringSlice(a))) }
