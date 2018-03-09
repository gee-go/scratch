package pprof_label

import "unsafe"

// LabelMap must match labelMap from runtime/pprof/label.go
type LabelMap map[string]string

// runtime_setProfLabel is defined in runtime/proflabel.go.
//go:linkname runtime_setProfLabel runtime/pprof.runtime_setProfLabel
func runtime_setProfLabel(labels unsafe.Pointer)

// runtime_getProfLabel is defined in runtime/proflabel.go.
//go:linkname runtime_getProfLabel runtime/pprof.runtime_getProfLabel
func runtime_getProfLabel() unsafe.Pointer

// Set pprof labels for the current go routine.
func Set(labels map[string]string) {
	runtime_setProfLabel(unsafe.Pointer(&labels))
}

// Get pprof labels for the current go routine.
func Get() LabelMap {
	l := (*LabelMap)(runtime_getProfLabel())
	if l == nil {
		return LabelMap{}
	}
	return *l
}
