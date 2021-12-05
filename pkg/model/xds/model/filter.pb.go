// Code generated by protoc-gen-go. DO NOT EDIT.
// source: filter.proto

package model

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Filter core struct, filter is extend by user
type Filter struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()    {}
func (*Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f5303cab7a20d6f, []int{0}
}

func (m *Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Filter.Unmarshal(m, b)
}
func (m *Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Filter.Marshal(b, m, deterministic)
}
func (m *Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter.Merge(m, src)
}
func (m *Filter) XXX_Size() int {
	return xxx_messageInfo_Filter.Size(m)
}
func (m *Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Filter proto.InternalMessageInfo

func (m *Filter) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// FilterChain filter chain
type FilterChain struct {
	FilterChainMatch     *FilterChainMatch `protobuf:"bytes,1,opt,name=filter_chain_match,json=filterChainMatch,proto3" json:"filter_chain_match,omitempty"`
	Filters              []*Filter         `protobuf:"bytes,2,rep,name=filters,proto3" json:"filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *FilterChain) Reset()         { *m = FilterChain{} }
func (m *FilterChain) String() string { return proto.CompactTextString(m) }
func (*FilterChain) ProtoMessage()    {}
func (*FilterChain) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f5303cab7a20d6f, []int{1}
}

func (m *FilterChain) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilterChain.Unmarshal(m, b)
}
func (m *FilterChain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilterChain.Marshal(b, m, deterministic)
}
func (m *FilterChain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilterChain.Merge(m, src)
}
func (m *FilterChain) XXX_Size() int {
	return xxx_messageInfo_FilterChain.Size(m)
}
func (m *FilterChain) XXX_DiscardUnknown() {
	xxx_messageInfo_FilterChain.DiscardUnknown(m)
}

var xxx_messageInfo_FilterChain proto.InternalMessageInfo

func (m *FilterChain) GetFilterChainMatch() *FilterChainMatch {
	if m != nil {
		return m.FilterChainMatch
	}
	return nil
}

func (m *FilterChain) GetFilters() []*Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

// FilterChainMatch
type FilterChainMatch struct {
	Domains              []string `protobuf:"bytes,1,rep,name=domains,proto3" json:"domains,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FilterChainMatch) Reset()         { *m = FilterChainMatch{} }
func (m *FilterChainMatch) String() string { return proto.CompactTextString(m) }
func (*FilterChainMatch) ProtoMessage()    {}
func (*FilterChainMatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f5303cab7a20d6f, []int{2}
}

func (m *FilterChainMatch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilterChainMatch.Unmarshal(m, b)
}
func (m *FilterChainMatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilterChainMatch.Marshal(b, m, deterministic)
}
func (m *FilterChainMatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilterChainMatch.Merge(m, src)
}
func (m *FilterChainMatch) XXX_Size() int {
	return xxx_messageInfo_FilterChainMatch.Size(m)
}
func (m *FilterChainMatch) XXX_DiscardUnknown() {
	xxx_messageInfo_FilterChainMatch.DiscardUnknown(m)
}

var xxx_messageInfo_FilterChainMatch proto.InternalMessageInfo

func (m *FilterChainMatch) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

// HTTPFilter http filter
type HTTPFilter struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HTTPFilter) Reset()         { *m = HTTPFilter{} }
func (m *HTTPFilter) String() string { return proto.CompactTextString(m) }
func (*HTTPFilter) ProtoMessage()    {}
func (*HTTPFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f5303cab7a20d6f, []int{3}
}

func (m *HTTPFilter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPFilter.Unmarshal(m, b)
}
func (m *HTTPFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPFilter.Marshal(b, m, deterministic)
}
func (m *HTTPFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPFilter.Merge(m, src)
}
func (m *HTTPFilter) XXX_Size() int {
	return xxx_messageInfo_HTTPFilter.Size(m)
}
func (m *HTTPFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPFilter.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPFilter proto.InternalMessageInfo

func (m *HTTPFilter) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Filter)(nil), "pixiu.api.v1.Filter")
	proto.RegisterType((*FilterChain)(nil), "pixiu.api.v1.FilterChain")
	proto.RegisterType((*FilterChainMatch)(nil), "pixiu.api.v1.FilterChainMatch")
	proto.RegisterType((*HTTPFilter)(nil), "pixiu.api.v1.HTTPFilter")
}

func init() { proto.RegisterFile("filter.proto", fileDescriptor_1f5303cab7a20d6f) }

var fileDescriptor_1f5303cab7a20d6f = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xcb, 0xcc, 0x29,
	0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x29, 0xc8, 0xac, 0xc8, 0x2c, 0xd5,
	0x4b, 0x2c, 0xc8, 0xd4, 0x2b, 0x33, 0x54, 0x92, 0xe1, 0x62, 0x73, 0x03, 0xcb, 0x0a, 0x09, 0x71,
	0xb1, 0xe4, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0xdd,
	0x8c, 0x5c, 0xdc, 0x10, 0x69, 0xe7, 0x8c, 0xc4, 0xcc, 0x3c, 0x21, 0x1f, 0x2e, 0x21, 0x88, 0x59,
	0xf1, 0xc9, 0x20, 0x7e, 0x7c, 0x6e, 0x62, 0x49, 0x72, 0x06, 0x58, 0x07, 0xb7, 0x91, 0x9c, 0x1e,
	0xb2, 0xc1, 0x7a, 0x48, 0xda, 0x7c, 0x41, 0xaa, 0x82, 0x04, 0xd2, 0xd0, 0x44, 0x84, 0xf4, 0xb8,
	0xd8, 0x21, 0x62, 0xc5, 0x12, 0x4c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0x22, 0xd8, 0x8c, 0x08, 0x82,
	0x29, 0x52, 0xd2, 0xe1, 0x12, 0x40, 0x37, 0x55, 0x48, 0x82, 0x8b, 0x3d, 0x25, 0x3f, 0x37, 0x31,
	0x33, 0xaf, 0x58, 0x82, 0x51, 0x81, 0x59, 0x83, 0x33, 0x08, 0xc6, 0x55, 0x52, 0xe0, 0xe2, 0xf2,
	0x08, 0x09, 0x09, 0xc0, 0xed, 0xbb, 0x24, 0x36, 0x70, 0x80, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x5b, 0x35, 0x60, 0x78, 0x20, 0x01, 0x00, 0x00,
}
