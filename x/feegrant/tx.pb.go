// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/feegrant/v1beta1/tx.proto

package feegrant

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/Finschia/finschia-sdk/codec/types"
	_ "github.com/regen-network/cosmos-proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
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

// MsgGrantAllowance adds permission for Grantee to spend up to Allowance
// of fees from the account of Granter.
type MsgGrantAllowance struct {
	// granter is the address of the user granting an allowance of their funds.
	Granter string `protobuf:"bytes,1,opt,name=granter,proto3" json:"granter,omitempty" yaml:"granter_address"`
	// grantee is the address of the user being granted an allowance of another user's funds.
	Grantee string `protobuf:"bytes,2,opt,name=grantee,proto3" json:"grantee,omitempty" yaml:"grantee_address"`
	// allowance can be any of basic and filtered fee allowance.
	Allowance *types.Any `protobuf:"bytes,3,opt,name=allowance,proto3" json:"allowance,omitempty"`
}

func (m *MsgGrantAllowance) Reset()         { *m = MsgGrantAllowance{} }
func (m *MsgGrantAllowance) String() string { return proto.CompactTextString(m) }
func (*MsgGrantAllowance) ProtoMessage()    {}
func (*MsgGrantAllowance) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd44ad7946dad783, []int{0}
}
func (m *MsgGrantAllowance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgGrantAllowance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgGrantAllowance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgGrantAllowance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgGrantAllowance.Merge(m, src)
}
func (m *MsgGrantAllowance) XXX_Size() int {
	return m.Size()
}
func (m *MsgGrantAllowance) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgGrantAllowance.DiscardUnknown(m)
}

var xxx_messageInfo_MsgGrantAllowance proto.InternalMessageInfo

func (m *MsgGrantAllowance) GetGranter() string {
	if m != nil {
		return m.Granter
	}
	return ""
}

func (m *MsgGrantAllowance) GetGrantee() string {
	if m != nil {
		return m.Grantee
	}
	return ""
}

func (m *MsgGrantAllowance) GetAllowance() *types.Any {
	if m != nil {
		return m.Allowance
	}
	return nil
}

// MsgGrantAllowanceResponse defines the Msg/GrantAllowanceResponse response type.
type MsgGrantAllowanceResponse struct {
}

func (m *MsgGrantAllowanceResponse) Reset()         { *m = MsgGrantAllowanceResponse{} }
func (m *MsgGrantAllowanceResponse) String() string { return proto.CompactTextString(m) }
func (*MsgGrantAllowanceResponse) ProtoMessage()    {}
func (*MsgGrantAllowanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd44ad7946dad783, []int{1}
}
func (m *MsgGrantAllowanceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgGrantAllowanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgGrantAllowanceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgGrantAllowanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgGrantAllowanceResponse.Merge(m, src)
}
func (m *MsgGrantAllowanceResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgGrantAllowanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgGrantAllowanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgGrantAllowanceResponse proto.InternalMessageInfo

// MsgRevokeAllowance removes any existing Allowance from Granter to Grantee.
type MsgRevokeAllowance struct {
	// granter is the address of the user granting an allowance of their funds.
	Granter string `protobuf:"bytes,1,opt,name=granter,proto3" json:"granter,omitempty" yaml:"granter_address"`
	// grantee is the address of the user being granted an allowance of another user's funds.
	Grantee string `protobuf:"bytes,2,opt,name=grantee,proto3" json:"grantee,omitempty" yaml:"grantee_address"`
}

func (m *MsgRevokeAllowance) Reset()         { *m = MsgRevokeAllowance{} }
func (m *MsgRevokeAllowance) String() string { return proto.CompactTextString(m) }
func (*MsgRevokeAllowance) ProtoMessage()    {}
func (*MsgRevokeAllowance) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd44ad7946dad783, []int{2}
}
func (m *MsgRevokeAllowance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRevokeAllowance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRevokeAllowance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRevokeAllowance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRevokeAllowance.Merge(m, src)
}
func (m *MsgRevokeAllowance) XXX_Size() int {
	return m.Size()
}
func (m *MsgRevokeAllowance) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRevokeAllowance.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRevokeAllowance proto.InternalMessageInfo

func (m *MsgRevokeAllowance) GetGranter() string {
	if m != nil {
		return m.Granter
	}
	return ""
}

func (m *MsgRevokeAllowance) GetGrantee() string {
	if m != nil {
		return m.Grantee
	}
	return ""
}

// MsgRevokeAllowanceResponse defines the Msg/RevokeAllowanceResponse response type.
type MsgRevokeAllowanceResponse struct {
}

func (m *MsgRevokeAllowanceResponse) Reset()         { *m = MsgRevokeAllowanceResponse{} }
func (m *MsgRevokeAllowanceResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRevokeAllowanceResponse) ProtoMessage()    {}
func (*MsgRevokeAllowanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd44ad7946dad783, []int{3}
}
func (m *MsgRevokeAllowanceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRevokeAllowanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRevokeAllowanceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRevokeAllowanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRevokeAllowanceResponse.Merge(m, src)
}
func (m *MsgRevokeAllowanceResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgRevokeAllowanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRevokeAllowanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRevokeAllowanceResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgGrantAllowance)(nil), "cosmos.feegrant.v1beta1.MsgGrantAllowance")
	proto.RegisterType((*MsgGrantAllowanceResponse)(nil), "cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse")
	proto.RegisterType((*MsgRevokeAllowance)(nil), "cosmos.feegrant.v1beta1.MsgRevokeAllowance")
	proto.RegisterType((*MsgRevokeAllowanceResponse)(nil), "cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse")
}

func init() { proto.RegisterFile("cosmos/feegrant/v1beta1/tx.proto", fileDescriptor_dd44ad7946dad783) }

var fileDescriptor_dd44ad7946dad783 = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x93, 0x31, 0x4f, 0xc2, 0x40,
	0x14, 0xc7, 0x39, 0x49, 0x34, 0x9c, 0x51, 0x43, 0x43, 0x14, 0xaa, 0xa9, 0xa4, 0x13, 0xd1, 0x70,
	0x17, 0xc0, 0xc9, 0xb8, 0x40, 0xa2, 0xc6, 0x81, 0xa5, 0xa3, 0x0b, 0xb9, 0xc2, 0xe3, 0x24, 0xb4,
	0x3d, 0xd2, 0x2b, 0x08, 0x9b, 0x1f, 0xc1, 0x0f, 0xe3, 0x67, 0x30, 0xc6, 0x89, 0xd1, 0xc9, 0x18,
	0x58, 0x9d, 0xfc, 0x04, 0x86, 0x96, 0x83, 0x84, 0xaa, 0xd1, 0xc9, 0xad, 0x2f, 0xef, 0xf7, 0x7f,
	0xff, 0xf7, 0x6f, 0x5f, 0x71, 0xbe, 0x29, 0xa4, 0x2b, 0x24, 0x6d, 0x03, 0x70, 0x9f, 0x79, 0x01,
	0x1d, 0x94, 0x6c, 0x08, 0x58, 0x89, 0x06, 0x43, 0xd2, 0xf3, 0x45, 0x20, 0xb4, 0xbd, 0x88, 0x20,
	0x8a, 0x20, 0x73, 0x42, 0xcf, 0x70, 0xc1, 0x45, 0xc8, 0xd0, 0xd9, 0x53, 0x84, 0xeb, 0x39, 0x2e,
	0x04, 0x77, 0x80, 0x86, 0x95, 0xdd, 0x6f, 0x53, 0xe6, 0x8d, 0x54, 0x2b, 0x9a, 0xd4, 0x88, 0x34,
	0xf3, 0xb1, 0x61, 0x61, 0x3e, 0x22, 0x9c, 0xae, 0x4b, 0x7e, 0x39, 0x33, 0xa8, 0x3a, 0x8e, 0xb8,
	0x65, 0x5e, 0x13, 0xb4, 0x13, 0xbc, 0x11, 0x5a, 0x82, 0x9f, 0x45, 0x79, 0x54, 0x48, 0xd5, 0xf4,
	0x8f, 0xd7, 0xc3, 0xdd, 0x11, 0x73, 0x9d, 0x53, 0x73, 0xde, 0x68, 0xb0, 0x56, 0xcb, 0x07, 0x29,
	0x4d, 0x4b, 0xa1, 0x4b, 0x15, 0x64, 0xd7, 0xbe, 0x56, 0x41, 0x4c, 0x05, 0xda, 0x39, 0x4e, 0x31,
	0x65, 0x9c, 0x4d, 0xe6, 0x51, 0x61, 0xb3, 0x9c, 0x21, 0x51, 0x16, 0xa2, 0xb2, 0x90, 0xaa, 0x37,
	0xaa, 0xa5, 0x9f, 0x1f, 0x8a, 0x5b, 0x17, 0x00, 0x8b, 0x35, 0xaf, 0xac, 0xa5, 0xd2, 0xdc, 0xc7,
	0xb9, 0x58, 0x0e, 0x0b, 0x64, 0x4f, 0x78, 0x12, 0xcc, 0x3b, 0x84, 0xb5, 0xba, 0xe4, 0x16, 0x0c,
	0x44, 0x17, 0xfe, 0x25, 0xa6, 0x79, 0x80, 0xf5, 0xf8, 0x06, 0x6a, 0xc1, 0xf2, 0x3b, 0xc2, 0xc9,
	0xba, 0xe4, 0x5a, 0x0f, 0x6f, 0xaf, 0x7c, 0x8a, 0x23, 0xf2, 0xcd, 0x19, 0x90, 0x58, 0x5c, 0xbd,
	0xfc, 0x7b, 0x56, 0x39, 0x6b, 0x12, 0xef, 0xac, 0xbe, 0x96, 0xe3, 0x9f, 0xc6, 0xac, 0xc0, 0x7a,
	0xe5, 0x0f, 0xb0, 0x32, 0xad, 0x9d, 0x3d, 0x4d, 0x0c, 0x34, 0x9e, 0x18, 0xe8, 0x6d, 0x62, 0xa0,
	0xfb, 0xa9, 0x91, 0x18, 0x4f, 0x8d, 0xc4, 0xcb, 0xd4, 0x48, 0x5c, 0x9b, 0xbc, 0x13, 0xdc, 0xf4,
	0x6d, 0xd2, 0x14, 0x2e, 0x75, 0x3a, 0x1e, 0x50, 0xc7, 0x76, 0x8b, 0xb2, 0xd5, 0xa5, 0xc3, 0xc5,
	0x9f, 0x62, 0xaf, 0x87, 0x67, 0x51, 0xf9, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x60, 0x30, 0x87, 0x1f,
	0x43, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// GrantAllowance grants fee allowance to the grantee on the granter's
	// account with the provided expiration time.
	GrantAllowance(ctx context.Context, in *MsgGrantAllowance, opts ...grpc.CallOption) (*MsgGrantAllowanceResponse, error)
	// RevokeAllowance revokes any fee allowance of granter's account that
	// has been granted to the grantee.
	RevokeAllowance(ctx context.Context, in *MsgRevokeAllowance, opts ...grpc.CallOption) (*MsgRevokeAllowanceResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) GrantAllowance(ctx context.Context, in *MsgGrantAllowance, opts ...grpc.CallOption) (*MsgGrantAllowanceResponse, error) {
	out := new(MsgGrantAllowanceResponse)
	err := c.cc.Invoke(ctx, "/cosmos.feegrant.v1beta1.Msg/GrantAllowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RevokeAllowance(ctx context.Context, in *MsgRevokeAllowance, opts ...grpc.CallOption) (*MsgRevokeAllowanceResponse, error) {
	out := new(MsgRevokeAllowanceResponse)
	err := c.cc.Invoke(ctx, "/cosmos.feegrant.v1beta1.Msg/RevokeAllowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// GrantAllowance grants fee allowance to the grantee on the granter's
	// account with the provided expiration time.
	GrantAllowance(context.Context, *MsgGrantAllowance) (*MsgGrantAllowanceResponse, error)
	// RevokeAllowance revokes any fee allowance of granter's account that
	// has been granted to the grantee.
	RevokeAllowance(context.Context, *MsgRevokeAllowance) (*MsgRevokeAllowanceResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) GrantAllowance(ctx context.Context, req *MsgGrantAllowance) (*MsgGrantAllowanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrantAllowance not implemented")
}
func (*UnimplementedMsgServer) RevokeAllowance(ctx context.Context, req *MsgRevokeAllowance) (*MsgRevokeAllowanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeAllowance not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_GrantAllowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgGrantAllowance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).GrantAllowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.feegrant.v1beta1.Msg/GrantAllowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).GrantAllowance(ctx, req.(*MsgGrantAllowance))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RevokeAllowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRevokeAllowance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RevokeAllowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.feegrant.v1beta1.Msg/RevokeAllowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RevokeAllowance(ctx, req.(*MsgRevokeAllowance))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cosmos.feegrant.v1beta1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GrantAllowance",
			Handler:    _Msg_GrantAllowance_Handler,
		},
		{
			MethodName: "RevokeAllowance",
			Handler:    _Msg_RevokeAllowance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cosmos/feegrant/v1beta1/tx.proto",
}

func (m *MsgGrantAllowance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgGrantAllowance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgGrantAllowance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Allowance != nil {
		{
			size, err := m.Allowance.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Grantee) > 0 {
		i -= len(m.Grantee)
		copy(dAtA[i:], m.Grantee)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Grantee)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Granter) > 0 {
		i -= len(m.Granter)
		copy(dAtA[i:], m.Granter)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Granter)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgGrantAllowanceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgGrantAllowanceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgGrantAllowanceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgRevokeAllowance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRevokeAllowance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRevokeAllowance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Grantee) > 0 {
		i -= len(m.Grantee)
		copy(dAtA[i:], m.Grantee)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Grantee)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Granter) > 0 {
		i -= len(m.Granter)
		copy(dAtA[i:], m.Granter)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Granter)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgRevokeAllowanceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRevokeAllowanceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRevokeAllowanceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgGrantAllowance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Granter)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Grantee)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Allowance != nil {
		l = m.Allowance.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgGrantAllowanceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgRevokeAllowance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Granter)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Grantee)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgRevokeAllowanceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgGrantAllowance) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgGrantAllowance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgGrantAllowance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Granter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Granter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Grantee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Grantee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Allowance", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Allowance == nil {
				m.Allowance = &types.Any{}
			}
			if err := m.Allowance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgGrantAllowanceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgGrantAllowanceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgGrantAllowanceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgRevokeAllowance) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRevokeAllowance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRevokeAllowance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Granter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Granter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Grantee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Grantee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgRevokeAllowanceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRevokeAllowanceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRevokeAllowanceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
