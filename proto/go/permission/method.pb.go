// Code generated by protoc-gen-go. DO NOT EDIT.
// source: method.proto

package permission

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Method token subject.
type Subject int32

const (
	Subject_NONE      Subject = 0
	Subject_GUEST     Subject = 1
	Subject_WEB       Subject = 8
	Subject_PC        Subject = 64
	Subject_MOBILE    Subject = 512
	Subject_LOGGED_IN Subject = 584
	Subject_CLIENT    Subject = 585
	Subject_SERVER    Subject = 4096
	Subject_ANY       Subject = 65535
)

var Subject_name = map[int32]string{
	0:     "NONE",
	1:     "GUEST",
	8:     "WEB",
	64:    "PC",
	512:   "MOBILE",
	584:   "LOGGED_IN",
	585:   "CLIENT",
	4096:  "SERVER",
	65535: "ANY",
}

var Subject_value = map[string]int32{
	"NONE":      0,
	"GUEST":     1,
	"WEB":       8,
	"PC":        64,
	"MOBILE":    512,
	"LOGGED_IN": 584,
	"CLIENT":    585,
	"SERVER":    4096,
	"ANY":       65535,
}

func (x Subject) String() string {
	return proto.EnumName(Subject_name, int32(x))
}

func (Subject) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4d10d17d32e60d7d, []int{0}
}

var E_Required = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: ([]Subject)(nil),
	Field:         2507,
	Name:          "appootb.permission.method.required",
	Tag:           "varint,2507,rep,name=required,enum=appootb.permission.method.Subject",
	Filename:      "method.proto",
}

func init() {
	proto.RegisterEnum("appootb.permission.method.Subject", Subject_name, Subject_value)
	proto.RegisterExtension(E_Required)
}

func init() { proto.RegisterFile("method.proto", fileDescriptor_4d10d17d32e60d7d) }

var fileDescriptor_4d10d17d32e60d7d = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x4d, 0x4f, 0xc2, 0x30,
	0x18, 0xc7, 0xdd, 0x86, 0xbc, 0x14, 0x43, 0x9a, 0x9a, 0x18, 0xf4, 0x60, 0x88, 0x27, 0xe3, 0xa1,
	0x8b, 0x7a, 0xf3, 0x24, 0xc3, 0x06, 0x49, 0x60, 0x5b, 0x00, 0x35, 0x1a, 0x12, 0xc2, 0x46, 0x1d,
	0x33, 0x8e, 0x67, 0x74, 0xdd, 0x9d, 0x2f, 0xe0, 0x97, 0xf0, 0xe8, 0x27, 0xf1, 0xe5, 0x0b, 0x79,
	0xc3, 0xb0, 0x22, 0x5c, 0xf4, 0xd8, 0x3e, 0xff, 0xff, 0xaf, 0xbf, 0x3e, 0x68, 0x27, 0xe2, 0x72,
	0x02, 0x63, 0x1a, 0x0b, 0x90, 0x40, 0xf6, 0x47, 0x71, 0x0c, 0x20, 0x3d, 0x1a, 0x73, 0x11, 0x85,
	0x49, 0x12, 0xc2, 0x94, 0xaa, 0xc0, 0x41, 0x2d, 0x00, 0x08, 0x9e, 0xb9, 0x99, 0x05, 0xbd, 0xf4,
	0xd1, 0x1c, 0xf3, 0xc4, 0x17, 0x61, 0x2c, 0x41, 0xa8, 0xf2, 0xc9, 0x0c, 0x15, 0x7a, 0xa9, 0xf7,
	0xc4, 0x7d, 0x49, 0x8a, 0x28, 0x67, 0x3b, 0x36, 0xc3, 0x5b, 0xa4, 0x84, 0xb6, 0x9b, 0x37, 0xac,
	0xd7, 0xc7, 0x1a, 0x29, 0x20, 0xe3, 0x8e, 0x59, 0xb8, 0x48, 0xf2, 0x48, 0x77, 0x1b, 0xf8, 0x92,
	0x94, 0x51, 0xbe, 0xe3, 0x58, 0xad, 0x36, 0xc3, 0xf3, 0x1c, 0xa9, 0xa0, 0x52, 0xdb, 0x69, 0x36,
	0xd9, 0xd5, 0xb0, 0x65, 0xe3, 0xf7, 0xdc, 0x72, 0xd8, 0x68, 0xb7, 0x98, 0xdd, 0xc7, 0x1f, 0xd9,
	0xa1, 0xc7, 0xba, 0xb7, 0xac, 0x8b, 0xe7, 0x35, 0x52, 0x42, 0x46, 0xdd, 0xbe, 0xc7, 0x8b, 0x85,
	0x71, 0x31, 0x44, 0x45, 0xc1, 0x67, 0x69, 0x28, 0xf8, 0x98, 0x1c, 0x52, 0x65, 0x48, 0x7f, 0x0d,
	0x69, 0x27, 0x33, 0x77, 0x62, 0x19, 0xc2, 0x34, 0xa9, 0x7e, 0xed, 0xd6, 0x8c, 0xe3, 0xca, 0xd9,
	0x11, 0xfd, 0xf7, 0x8f, 0x74, 0xa5, 0xdf, 0x5d, 0x43, 0xad, 0x17, 0x0d, 0xed, 0xf9, 0x10, 0xfd,
	0xd1, 0xb1, 0xca, 0x0a, 0xef, 0x2e, 0x5f, 0xbb, 0xd6, 0x5c, 0xed, 0xe1, 0x34, 0x08, 0xe5, 0x24,
	0xf5, 0xa8, 0x0f, 0x91, 0xb9, 0xca, 0x9b, 0x49, 0xea, 0x25, 0x52, 0x8c, 0x64, 0x1a, 0xa9, 0xc5,
	0x99, 0x01, 0x98, 0x1b, 0xc6, 0xb7, 0xa6, 0xbd, 0xea, 0x46, 0xc3, 0xb5, 0xde, 0x74, 0xe4, 0xae,
	0x6f, 0x3f, 0xf5, 0x6a, 0x5d, 0xd5, 0x07, 0x19, 0x7c, 0xb0, 0x19, 0x79, 0xf9, 0x8c, 0x72, 0xfe,
	0x13, 0x00, 0x00, 0xff, 0xff, 0x70, 0x47, 0x6d, 0x6b, 0xb7, 0x01, 0x00, 0x00,
}