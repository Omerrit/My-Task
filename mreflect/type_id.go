package mreflect

import (
	"unsafe"
)

//NOTE: stolen from reflect, check for validity every new go version

type tflag uint8

type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype
type textOff int32 // offset from top of text section

type rtype struct {
	size       uintptr
	ptrdata    uintptr  // number of bytes in the type that can contain pointers
	hash       uint32   // hash of type; avoids computation in hash tables
	tflag      tflag    // extra type information flags
	align      uint8    // alignment of variable with this type
	fieldAlign uint8    // alignment of struct field with this type
	kind       uint8    // enumeration for C
	alg        *typeAlg // algorithm table
	gcdata     *byte    // garbage collection data
	str        nameOff  // string form
	ptrToThis  typeOff  // type for pointer to this type, may be zero
}

// a copy of runtime.typeAlg
type typeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}

type TypeId uintptr

func GetTypeId(value interface{}) TypeId {
	return TypeId(unsafe.Pointer((*emptyInterface)(unsafe.Pointer(&value)).typ))
}
