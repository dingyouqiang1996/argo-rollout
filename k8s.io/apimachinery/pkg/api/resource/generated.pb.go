// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: k8s.io/apimachinery/pkg/api/resource/generated.proto

package resource

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Quantity is a fixed-point representation of a number.
// It provides convenient marshaling/unmarshaling in JSON and YAML,
// in addition to String() and AsInt64() accessors.
//
// The serialization format is:
//
// ```
// <quantity>        ::= <signedNumber><suffix>
//
//	(Note that <suffix> may be empty, from the "" case in <decimalSI>.)
//
// <digit>           ::= 0 | 1 | ... | 9
// <digits>          ::= <digit> | <digit><digits>
// <number>          ::= <digits> | <digits>.<digits> | <digits>. | .<digits>
// <sign>            ::= "+" | "-"
// <signedNumber>    ::= <number> | <sign><number>
// <suffix>          ::= <binarySI> | <decimalExponent> | <decimalSI>
// <binarySI>        ::= Ki | Mi | Gi | Ti | Pi | Ei
//
//	(International System of units; See: http://physics.nist.gov/cuu/Units/binary.html)
//
// <decimalSI>       ::= m | "" | k | M | G | T | P | E
//
//	(Note that 1024 = 1Ki but 1000 = 1k; I didn't choose the capitalization.)
//
// <decimalExponent> ::= "e" <signedNumber> | "E" <signedNumber>
// ```
//
// No matter which of the three exponent forms is used, no quantity may represent
// a number greater than 2^63-1 in magnitude, nor may it have more than 3 decimal
// places. Numbers larger or more precise will be capped or rounded up.
// (E.g.: 0.1m will rounded up to 1m.)
// This may be extended in the future if we require larger or smaller quantities.
//
// When a Quantity is parsed from a string, it will remember the type of suffix
// it had, and will use the same type again when it is serialized.
//
// Before serializing, Quantity will be put in "canonical form".
// This means that Exponent/suffix will be adjusted up or down (with a
// corresponding increase or decrease in Mantissa) such that:
//
// - No precision is lost
// - No fractional digits will be emitted
// - The exponent (or suffix) is as large as possible.
//
// The sign will be omitted unless the number is negative.
//
// Examples:
//
// - 1.5 will be serialized as "1500m"
// - 1.5Gi will be serialized as "1536Mi"
//
// Note that the quantity will NEVER be internally represented by a
// floating point number. That is the whole point of this exercise.
//
// Non-canonical values will still parse as long as they are well formed,
// but will be re-emitted in their canonical form. (So always use canonical
// form, or don't diff.)
//
// This format is intended to make it difficult to use these numbers without
// writing some sort of special handling code in the hopes that that will
// cause implementors to also use a fixed point implementation.
//
// +protobuf=true
// +protobuf.embed=string
// +protobuf.options.marshal=false
// +protobuf.options.(gogoproto.goproto_stringer)=false
// +k8s:deepcopy-gen=true
// +k8s:openapi-gen=true
type Quantity struct {
	String_              *string  `protobuf:"bytes,1,opt,name=string" json:"string,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Quantity) Reset()      { *m = Quantity{} }
func (*Quantity) ProtoMessage() {}
func (*Quantity) Descriptor() ([]byte, []int) {
	return fileDescriptor_7288c78ff45111e9, []int{0}
}
func (m *Quantity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Quantity.Unmarshal(m, b)
}
func (m *Quantity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Quantity.Marshal(b, m, deterministic)
}
func (m *Quantity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Quantity.Merge(m, src)
}
func (m *Quantity) XXX_Size() int {
	return xxx_messageInfo_Quantity.Size(m)
}
func (m *Quantity) XXX_DiscardUnknown() {
	xxx_messageInfo_Quantity.DiscardUnknown(m)
}

var xxx_messageInfo_Quantity proto.InternalMessageInfo

// QuantityValue makes it possible to use a Quantity as value for a command
// line parameter.
//
// +protobuf=true
// +protobuf.embed=string
// +protobuf.options.marshal=false
// +protobuf.options.(gogoproto.goproto_stringer)=false
// +k8s:deepcopy-gen=true
type QuantityValue struct {
	String_              *string  `protobuf:"bytes,1,opt,name=string" json:"string,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QuantityValue) Reset()      { *m = QuantityValue{} }
func (*QuantityValue) ProtoMessage() {}
func (*QuantityValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_7288c78ff45111e9, []int{1}
}
func (m *QuantityValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QuantityValue.Unmarshal(m, b)
}
func (m *QuantityValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QuantityValue.Marshal(b, m, deterministic)
}
func (m *QuantityValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuantityValue.Merge(m, src)
}
func (m *QuantityValue) XXX_Size() int {
	return xxx_messageInfo_QuantityValue.Size(m)
}
func (m *QuantityValue) XXX_DiscardUnknown() {
	xxx_messageInfo_QuantityValue.DiscardUnknown(m)
}

var xxx_messageInfo_QuantityValue proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Quantity)(nil), "k8s.io.apimachinery.pkg.api.resource.Quantity")
	proto.RegisterType((*QuantityValue)(nil), "k8s.io.apimachinery.pkg.api.resource.QuantityValue")
}

func init() {
	proto.RegisterFile("k8s.io/apimachinery/pkg/api/resource/generated.proto", fileDescriptor_7288c78ff45111e9)
}

var fileDescriptor_7288c78ff45111e9 = []byte{
	// 234 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0xc9, 0xb6, 0x28, 0xd6,
	0xcb, 0xcc, 0xd7, 0x4f, 0x2c, 0xc8, 0xcc, 0x4d, 0x4c, 0xce, 0xc8, 0xcc, 0x4b, 0x2d, 0xaa, 0xd4,
	0x2f, 0xc8, 0x4e, 0x07, 0x09, 0xe8, 0x17, 0xa5, 0x16, 0xe7, 0x97, 0x16, 0x25, 0xa7, 0xea, 0xa7,
	0xa7, 0xe6, 0xa5, 0x16, 0x25, 0x96, 0xa4, 0xa6, 0xe8, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0xa9,
	0x40, 0x74, 0xe9, 0x21, 0xeb, 0xd2, 0x2b, 0xc8, 0x4e, 0x07, 0x09, 0xe8, 0xc1, 0x74, 0x49, 0xe9,
	0xa6, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0xa7, 0xe7, 0xa7, 0xe7, 0xeb,
	0x83, 0x35, 0x27, 0x95, 0xa6, 0x81, 0x79, 0x60, 0x0e, 0x98, 0x05, 0x31, 0x54, 0xc9, 0x82, 0x8b,
	0x23, 0xb0, 0x34, 0x31, 0xaf, 0x24, 0xb3, 0xa4, 0x52, 0x48, 0x8c, 0x8b, 0xad, 0xb8, 0xa4, 0x28,
	0x33, 0x2f, 0x5d, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xca, 0xb3, 0x12, 0x99, 0xb1, 0x40,
	0x9e, 0xa1, 0x63, 0xa1, 0x3c, 0xc3, 0x84, 0x85, 0xf2, 0x0c, 0x0b, 0x16, 0xca, 0x33, 0x34, 0xdc,
	0x51, 0x60, 0x50, 0xb2, 0xe5, 0xe2, 0x85, 0xe9, 0x0c, 0x4b, 0xcc, 0x29, 0x4d, 0x25, 0x4d, 0xbb,
	0x93, 0xd7, 0x89, 0x87, 0x72, 0x0c, 0x17, 0x1e, 0xca, 0x31, 0xdc, 0x78, 0x28, 0xc7, 0xd0, 0xf0,
	0x48, 0x8e, 0xf1, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x6f, 0x3c, 0x92, 0x63, 0x7c,
	0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x28, 0x15, 0x62, 0x42, 0x0a, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x50, 0x91, 0xd0, 0x9c, 0x50, 0x01, 0x00, 0x00,
}
