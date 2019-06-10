// Code generated by protoc-gen-go. DO NOT EDIT.
// source: upnp.proto

package upnp

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
	PortNumber           uint32   `protobuf:"varint,1,opt,name=portNumber,proto3" json:"portNumber,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralRequest) Reset()         { *m = GeneralRequest{} }
func (m *GeneralRequest) String() string { return proto.CompactTextString(m) }
func (*GeneralRequest) ProtoMessage()    {}
func (*GeneralRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6afd328382fff2d5, []int{0}
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

func (m *GeneralRequest) GetPortNumber() uint32 {
	if m != nil {
		return m.PortNumber
	}
	return 0
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
	return fileDescriptor_6afd328382fff2d5, []int{1}
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
	proto.RegisterType((*GeneralRequest)(nil), "upnp.GeneralRequest")
	proto.RegisterType((*GeneralResponse)(nil), "upnp.GeneralResponse")
}

func init() { proto.RegisterFile("upnp.proto", fileDescriptor_6afd328382fff2d5) }

var fileDescriptor_6afd328382fff2d5 = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2d, 0xc8, 0x2b,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x0c, 0xb8, 0xf8, 0xdc, 0x53,
	0xf3, 0x52, 0x8b, 0x12, 0x73, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xe4, 0xb8, 0xb8,
	0x0a, 0xf2, 0x8b, 0x4a, 0xfc, 0x4a, 0x73, 0x93, 0x52, 0x8b, 0x24, 0x18, 0x15, 0x18, 0x35, 0x78,
	0x83, 0x90, 0x44, 0x94, 0xb4, 0xb9, 0xf8, 0xe1, 0x3a, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85,
	0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0xc1, 0xea, 0x39, 0x83, 0x60, 0x5c,
	0xa3, 0x0e, 0x26, 0x2e, 0x96, 0xd0, 0x82, 0xbc, 0x02, 0x21, 0x6b, 0x2e, 0x2e, 0xf7, 0xd4, 0x12,
	0xf7, 0xc4, 0x92, 0xd4, 0xf2, 0xc4, 0x4a, 0x21, 0x11, 0x3d, 0xb0, 0x43, 0x50, 0x6d, 0x96, 0x12,
	0x45, 0x13, 0x85, 0x98, 0xae, 0xc4, 0x20, 0xe4, 0xc4, 0x25, 0xe8, 0x96, 0x5f, 0x54, 0x9e, 0x58,
	0x94, 0x12, 0x90, 0x5f, 0x54, 0x12, 0x9c, 0x99, 0x93, 0x9a, 0x57, 0x42, 0xaa, 0x19, 0x36, 0x5c,
	0xdc, 0x48, 0x66, 0x90, 0xaa, 0xdb, 0x91, 0x4b, 0x20, 0x28, 0x35, 0x37, 0xbf, 0x2c, 0x15, 0x6a,
	0x46, 0x66, 0x5e, 0x3a, 0x89, 0x46, 0x24, 0xb1, 0x81, 0x83, 0xdd, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0xf6, 0xd0, 0x8e, 0x39, 0x84, 0x01, 0x00, 0x00,
}
