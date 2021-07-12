// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo.exe -godefs ctypes/ctypes.go

package snap7go

type TS7CpuInfo struct {
	ModuleTypeName [33]int8
	SerialNumber   [25]int8
	ASName         [25]int8
	Copyright      [27]int8
	ModuleName     [25]int8
}
type Tm struct {
	Sec   int32
	Min   int32
	Hour  int32
	Mday  int32
	Mon   int32
	Year  int32
	Wday  int32
	Yday  int32
	Isdst int32
}

type TS7DataItem struct {
	Area     int32
	WordLen  int32
	Result   int32
	DBNumber int32
	Start    int32
	Amount   int32
	Pdata    *byte
}
