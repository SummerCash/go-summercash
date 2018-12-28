// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chainconfig.proto

package config

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GeneralRequest struct {
	GenesisPath          string   `protobuf:"bytes,1,opt,name=genesisPath,proto3" json:"genesisPath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralRequest) Reset()         { *m = GeneralRequest{} }
func (m *GeneralRequest) String() string { return proto.CompactTextString(m) }
func (*GeneralRequest) ProtoMessage()    {}
func (*GeneralRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ba5639a1e213847, []int{0}
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

func (m *GeneralRequest) GetGenesisPath() string {
	if m != nil {
		return m.GenesisPath
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
	return fileDescriptor_2ba5639a1e213847, []int{1}
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
	proto.RegisterType((*GeneralRequest)(nil), "config.GeneralRequest")
	proto.RegisterType((*GeneralResponse)(nil), "config.GeneralResponse")
}

func init() { proto.RegisterFile("chainconfig.proto", fileDescriptor_2ba5639a1e213847) }

var fileDescriptor_2ba5639a1e213847 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xce, 0x48, 0xcc,
	0xcc, 0x4b, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83,
	0xf0, 0x94, 0x8c, 0xb8, 0xf8, 0xdc, 0x53, 0xf3, 0x52, 0x8b, 0x12, 0x73, 0x82, 0x52, 0x0b, 0x4b,
	0x53, 0x8b, 0x4b, 0x84, 0x14, 0xb8, 0xb8, 0xd3, 0x53, 0xf3, 0x52, 0x8b, 0x33, 0x8b, 0x03, 0x12,
	0x4b, 0x32, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x90, 0x85, 0x94, 0xb4, 0xb9, 0xf8, 0xe1,
	0x7a, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13,
	0xd3, 0x53, 0xa1, 0x1a, 0x60, 0x5c, 0xa3, 0x5b, 0x4c, 0x5c, 0x6c, 0xce, 0x60, 0xbb, 0x84, 0x9c,
	0xb9, 0xf8, 0xfc, 0x52, 0xcb, 0x9d, 0x41, 0x6e, 0x81, 0x8a, 0x88, 0xe9, 0x41, 0x1d, 0x85, 0xea,
	0x06, 0x29, 0x71, 0x0c, 0x71, 0x88, 0x3d, 0x4a, 0x0c, 0x42, 0x56, 0x5c, 0xac, 0x4e, 0x95, 0x25,
	0xa9, 0xc5, 0xe4, 0xe8, 0xb5, 0xe6, 0x62, 0x0b, 0x2e, 0x29, 0xca, 0xcc, 0x23, 0xcb, 0x62, 0x27,
	0x2e, 0xde, 0xf0, 0xa2, 0xcc, 0x92, 0xd4, 0x90, 0x7c, 0xdf, 0xd4, 0xdc, 0xfc, 0xa2, 0x4a, 0x72,
	0xcc, 0xf0, 0xe3, 0x92, 0x0c, 0x4a, 0x4d, 0x4c, 0x41, 0x0a, 0x02, 0xb7, 0xa2, 0xfc, 0x5c, 0xb2,
	0xcd, 0x4b, 0x62, 0x03, 0x47, 0xa6, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xc9, 0x41, 0xd2, 0x17,
	0xe1, 0x01, 0x00, 0x00,
}
