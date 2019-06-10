// Code generated by protoc-gen-go. DO NOT EDIT.
// source: accounts.proto

package accounts

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GeneralRequest struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PrivateKey           string   `protobuf:"bytes,2,opt,name=privateKey,proto3" json:"privateKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralRequest) Reset()         { *m = GeneralRequest{} }
func (m *GeneralRequest) String() string { return proto.CompactTextString(m) }
func (*GeneralRequest) ProtoMessage()    {}
func (*GeneralRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1e7723af4c007b7, []int{0}
}

func (m *GeneralRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeneralRequest.Unmarshal(m, b)
}

func (m *GeneralRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeneralRequest.Marshal(b, m, deterministic)
}

func (m *GeneralRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeneralRequest.Merge(m, src)
}

func (m *GeneralRequest) XXX_Size() int {
	return xxx_messageInfo_GeneralRequest.Size(m)
}

func (m *GeneralRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GeneralRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GeneralRequest proto.InternalMessageInfo

func (m *GeneralRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *GeneralRequest) GetPrivateKey() string {
	if m != nil {
		return m.PrivateKey
	}
	return ""
}

type GeneralResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralResponse) Reset()         { *m = GeneralResponse{} }
func (m *GeneralResponse) String() string { return proto.CompactTextString(m) }
func (*GeneralResponse) ProtoMessage()    {}
func (*GeneralResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1e7723af4c007b7, []int{1}
}

func (m *GeneralResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeneralResponse.Unmarshal(m, b)
}

func (m *GeneralResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeneralResponse.Marshal(b, m, deterministic)
}

func (m *GeneralResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeneralResponse.Merge(m, src)
}

func (m *GeneralResponse) XXX_Size() int {
	return xxx_messageInfo_GeneralResponse.Size(m)
}

func (m *GeneralResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GeneralResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GeneralResponse proto.InternalMessageInfo

func (m *GeneralResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*GeneralRequest)(nil), "accounts.GeneralRequest")
	proto.RegisterType((*GeneralResponse)(nil), "accounts.GeneralResponse")
}

func init() { proto.RegisterFile("accounts.proto", fileDescriptor_e1e7723af4c007b7) }

var fileDescriptor_e1e7723af4c007b7 = []byte{
	// 279 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0xd3, 0x4f, 0x4b, 0xc3, 0x30,
	0x18, 0x06, 0x70, 0x27, 0x6c, 0xce, 0xf7, 0x30, 0x25, 0x22, 0x54, 0x0f, 0x22, 0x3d, 0x09, 0xc2,
	0x0e, 0x7a, 0x57, 0xe6, 0xd0, 0xaa, 0x75, 0x3d, 0x74, 0x9f, 0x20, 0xa6, 0x8f, 0x63, 0xd8, 0x26,
	0x35, 0xc9, 0x36, 0xfa, 0x75, 0xfd, 0x24, 0xd2, 0xd9, 0x94, 0x09, 0x9e, 0xd2, 0xe3, 0xd3, 0x17,
	0x7e, 0x7d, 0xff, 0x10, 0x1a, 0x71, 0x21, 0xd4, 0x4a, 0x5a, 0x33, 0x2e, 0xb5, 0xb2, 0x8a, 0x0d,
	0x5d, 0x0e, 0x5f, 0x69, 0x14, 0x41, 0x42, 0xf3, 0x3c, 0xc5, 0xd7, 0x0a, 0xc6, 0xb2, 0x80, 0x0e,
	0x78, 0x96, 0x69, 0x18, 0x13, 0xf4, 0x2e, 0x7b, 0x57, 0x87, 0xa9, 0x8b, 0xec, 0x82, 0xa8, 0xd4,
	0xcb, 0x35, 0xb7, 0x88, 0x51, 0x05, 0xfb, 0xdb, 0xe2, 0xce, 0x97, 0xf0, 0x9a, 0x8e, 0x5a, 0xcb,
	0x94, 0x4a, 0x1a, 0xd4, 0x58, 0x01, 0x63, 0xf8, 0x02, 0x0e, 0x6b, 0xe2, 0xcd, 0x77, 0x9f, 0x86,
	0x93, 0xa6, 0x0b, 0x36, 0x25, 0x4a, 0xb0, 0x69, 0x22, 0x0b, 0xc6, 0x6d, 0xbb, 0x7f, 0x7b, 0x3b,
	0x3f, 0xfb, 0xa7, 0xf2, 0xfb, 0xa7, 0x70, 0x8f, 0xc5, 0xc4, 0x12, 0x6c, 0xa6, 0x4a, 0x5a, 0xcd,
	0x85, 0xed, 0x88, 0x45, 0x34, 0x6a, 0x84, 0x27, 0xad, 0x8a, 0x18, 0x55, 0x07, 0x28, 0x82, 0x9d,
	0xe4, 0x79, 0x3b, 0xac, 0x27, 0xf4, 0x5c, 0x6f, 0xb7, 0x86, 0xdc, 0x84, 0xde, 0xd2, 0x0b, 0x1d,
	0xcf, 0xf8, 0x27, 0x1e, 0xa5, 0x50, 0xd9, 0x52, 0x2e, 0xe6, 0xfc, 0x03, 0xbe, 0xd4, 0x1b, 0x9d,
	0xa4, 0x10, 0x6a, 0x0d, 0x5d, 0x2b, 0x4e, 0xf4, 0xd5, 0xee, 0x69, 0x30, 0xb7, 0xba, 0x03, 0x70,
	0x47, 0xfd, 0x87, 0xca, 0xc2, 0x7b, 0x33, 0x09, 0x9d, 0xa6, 0xe0, 0xd9, 0xce, 0xe5, 0x67, 0x28,
	0x94, 0xf6, 0x3d, 0xfe, 0xfb, 0x60, 0xfb, 0xdc, 0x6e, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4e,
	0xb7, 0xd3, 0xaa, 0x80, 0x03, 0x00, 0x00,
}
