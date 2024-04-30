// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lbm/fswap/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	types "github.com/Finschia/finschia-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type QuerySwappedRequest struct {
}

func (m *QuerySwappedRequest) Reset()         { *m = QuerySwappedRequest{} }
func (m *QuerySwappedRequest) String() string { return proto.CompactTextString(m) }
func (*QuerySwappedRequest) ProtoMessage()    {}
func (*QuerySwappedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_01deae9da7816d6a, []int{0}
}
func (m *QuerySwappedRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySwappedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySwappedRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySwappedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySwappedRequest.Merge(m, src)
}
func (m *QuerySwappedRequest) XXX_Size() int {
	return m.Size()
}
func (m *QuerySwappedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySwappedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySwappedRequest proto.InternalMessageInfo

type QuerySwappedResponse struct {
	FromCoinAmount types.Coin `protobuf:"bytes,1,opt,name=from_coin_amount,json=fromCoinAmount,proto3,castrepeated=github.com/Finschia/finschia-sdk/types.Coin" json:"from_coin_amount"`
	ToCoinAmount   types.Coin `protobuf:"bytes,2,opt,name=to_coin_amount,json=toCoinAmount,proto3,castrepeated=github.com/Finschia/finschia-sdk/types.Coin" json:"to_coin_amount"`
}

func (m *QuerySwappedResponse) Reset()         { *m = QuerySwappedResponse{} }
func (m *QuerySwappedResponse) String() string { return proto.CompactTextString(m) }
func (*QuerySwappedResponse) ProtoMessage()    {}
func (*QuerySwappedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_01deae9da7816d6a, []int{1}
}
func (m *QuerySwappedResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySwappedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySwappedResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySwappedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySwappedResponse.Merge(m, src)
}
func (m *QuerySwappedResponse) XXX_Size() int {
	return m.Size()
}
func (m *QuerySwappedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySwappedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySwappedResponse proto.InternalMessageInfo

func (m *QuerySwappedResponse) GetFromCoinAmount() types.Coin {
	if m != nil {
		return m.FromCoinAmount
	}
	return types.Coin{}
}

func (m *QuerySwappedResponse) GetToCoinAmount() types.Coin {
	if m != nil {
		return m.ToCoinAmount
	}
	return types.Coin{}
}

type QueryTotalSwappableToCoinAmountRequest struct {
}

func (m *QueryTotalSwappableToCoinAmountRequest) Reset() {
	*m = QueryTotalSwappableToCoinAmountRequest{}
}
func (m *QueryTotalSwappableToCoinAmountRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTotalSwappableToCoinAmountRequest) ProtoMessage()    {}
func (*QueryTotalSwappableToCoinAmountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_01deae9da7816d6a, []int{2}
}
func (m *QueryTotalSwappableToCoinAmountRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTotalSwappableToCoinAmountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTotalSwappableToCoinAmountRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTotalSwappableToCoinAmountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTotalSwappableToCoinAmountRequest.Merge(m, src)
}
func (m *QueryTotalSwappableToCoinAmountRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTotalSwappableToCoinAmountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTotalSwappableToCoinAmountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTotalSwappableToCoinAmountRequest proto.InternalMessageInfo

type QueryTotalSwappableToCoinAmountResponse struct {
	SwappableAmount types.Coin `protobuf:"bytes,1,opt,name=swappable_amount,json=swappableAmount,proto3,castrepeated=github.com/Finschia/finschia-sdk/types.Coin" json:"swappable_amount"`
}

func (m *QueryTotalSwappableToCoinAmountResponse) Reset() {
	*m = QueryTotalSwappableToCoinAmountResponse{}
}
func (m *QueryTotalSwappableToCoinAmountResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTotalSwappableToCoinAmountResponse) ProtoMessage()    {}
func (*QueryTotalSwappableToCoinAmountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_01deae9da7816d6a, []int{3}
}
func (m *QueryTotalSwappableToCoinAmountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTotalSwappableToCoinAmountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTotalSwappableToCoinAmountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTotalSwappableToCoinAmountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTotalSwappableToCoinAmountResponse.Merge(m, src)
}
func (m *QueryTotalSwappableToCoinAmountResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTotalSwappableToCoinAmountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTotalSwappableToCoinAmountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTotalSwappableToCoinAmountResponse proto.InternalMessageInfo

func (m *QueryTotalSwappableToCoinAmountResponse) GetSwappableAmount() types.Coin {
	if m != nil {
		return m.SwappableAmount
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*QuerySwappedRequest)(nil), "lbm.fswap.v1.QuerySwappedRequest")
	proto.RegisterType((*QuerySwappedResponse)(nil), "lbm.fswap.v1.QuerySwappedResponse")
	proto.RegisterType((*QueryTotalSwappableToCoinAmountRequest)(nil), "lbm.fswap.v1.QueryTotalSwappableToCoinAmountRequest")
	proto.RegisterType((*QueryTotalSwappableToCoinAmountResponse)(nil), "lbm.fswap.v1.QueryTotalSwappableToCoinAmountResponse")
}

func init() { proto.RegisterFile("lbm/fswap/v1/query.proto", fileDescriptor_01deae9da7816d6a) }

var fileDescriptor_01deae9da7816d6a = []byte{
	// 444 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x3f, 0x8f, 0xd3, 0x30,
	0x18, 0xc6, 0x93, 0x48, 0x80, 0x64, 0x4e, 0xc7, 0xc9, 0xdc, 0x89, 0x23, 0x82, 0x1c, 0x64, 0xe0,
	0x0e, 0x21, 0x6c, 0xa5, 0x85, 0x0f, 0x40, 0x91, 0x10, 0x2b, 0xa5, 0x13, 0x4b, 0xe5, 0xa4, 0x6e,
	0x1a, 0x91, 0xf8, 0x4d, 0x6b, 0xa7, 0x7f, 0x56, 0x06, 0x66, 0x24, 0xbe, 0x02, 0x13, 0x3b, 0x0b,
	0x9f, 0xa0, 0x63, 0x25, 0x16, 0x26, 0x40, 0x2d, 0x1f, 0x04, 0xc5, 0x71, 0xa1, 0x91, 0x0a, 0x15,
	0x43, 0xb7, 0x57, 0x7e, 0x1f, 0xe7, 0xf7, 0xe4, 0x79, 0x5f, 0xa3, 0xd3, 0x34, 0xcc, 0x68, 0x5f,
	0x4e, 0x58, 0x4e, 0xc7, 0x01, 0x1d, 0x16, 0x7c, 0x34, 0x23, 0xf9, 0x08, 0x14, 0xe0, 0x83, 0x34,
	0xcc, 0x88, 0xee, 0x90, 0x71, 0xe0, 0xde, 0x8a, 0x01, 0xe2, 0x94, 0x53, 0x96, 0x27, 0x94, 0x09,
	0x01, 0x8a, 0xa9, 0x04, 0x84, 0xac, 0xb4, 0xee, 0x71, 0x0c, 0x31, 0xe8, 0x92, 0x96, 0x95, 0x39,
	0xf5, 0x22, 0x90, 0x19, 0x48, 0x1a, 0x32, 0xc9, 0xe9, 0x38, 0x08, 0xb9, 0x62, 0x01, 0x8d, 0x20,
	0x11, 0xa6, 0x5f, 0x67, 0x57, 0x28, 0xdd, 0xf1, 0x4f, 0xd0, 0xf5, 0x17, 0xa5, 0x95, 0x97, 0x13,
	0x96, 0xe7, 0xbc, 0xd7, 0xe6, 0xc3, 0x82, 0x4b, 0xe5, 0xbf, 0x75, 0xd0, 0x71, 0xfd, 0x5c, 0xe6,
	0x20, 0x24, 0xc7, 0x53, 0x74, 0xd4, 0x1f, 0x41, 0xd6, 0x2d, 0x3f, 0xde, 0x65, 0x19, 0x14, 0x42,
	0x9d, 0xda, 0x77, 0xec, 0x8b, 0xab, 0x8d, 0x9b, 0xa4, 0x32, 0x41, 0x4a, 0x13, 0xc4, 0x98, 0x20,
	0x4f, 0x21, 0x11, 0xad, 0xe6, 0xfc, 0xdb, 0x99, 0xf5, 0xf1, 0xfb, 0xd9, 0x83, 0x38, 0x51, 0x83,
	0x22, 0x24, 0x11, 0x64, 0xf4, 0x59, 0x22, 0x64, 0x34, 0x48, 0x18, 0xed, 0x9b, 0xe2, 0xa1, 0xec,
	0xbd, 0xa6, 0x6a, 0x96, 0x73, 0xa9, 0x2f, 0xb5, 0x0f, 0x4b, 0x4e, 0x59, 0x3d, 0xd1, 0x14, 0xac,
	0xd0, 0xa1, 0x82, 0x1a, 0xd7, 0xd9, 0x0b, 0xf7, 0x40, 0xc1, 0x1f, 0xaa, 0x7f, 0x81, 0xee, 0xe9,
	0x1c, 0x3a, 0xa0, 0x58, 0xaa, 0xc3, 0x60, 0x61, 0xca, 0x3b, 0x1b, 0x92, 0x75, 0x64, 0x1f, 0x6c,
	0x74, 0xbe, 0x53, 0x6a, 0x52, 0x9c, 0xa1, 0x23, 0xb9, 0x16, 0xec, 0x37, 0xc5, 0x6b, 0xbf, 0x39,
	0x95, 0x85, 0xc6, 0x27, 0x07, 0x5d, 0xd2, 0x36, 0x31, 0xa0, 0x2b, 0x66, 0xba, 0xf8, 0x2e, 0xd9,
	0x5c, 0x41, 0xb2, 0x65, 0x23, 0x5c, 0xff, 0x5f, 0x92, 0xea, 0xb7, 0xfc, 0xdb, 0x6f, 0xbe, 0xfc,
	0x7c, 0xef, 0xdc, 0xc0, 0x27, 0xb4, 0xb6, 0x6f, 0xd2, 0x50, 0x3e, 0xdb, 0xc8, 0xfd, 0x7b, 0x38,
	0xf8, 0xd1, 0x16, 0xc2, 0xce, 0xd8, 0xdd, 0xc7, 0xff, 0x79, 0xcb, 0x58, 0xa5, 0xda, 0xea, 0x7d,
	0x7c, 0xbe, 0xc5, 0xaa, 0x9e, 0x8a, 0xe0, 0x93, 0xcd, 0x65, 0x6b, 0x3d, 0x9f, 0x2f, 0x3d, 0x7b,
	0xb1, 0xf4, 0xec, 0x1f, 0x4b, 0xcf, 0x7e, 0xb7, 0xf2, 0xac, 0xc5, 0xca, 0xb3, 0xbe, 0xae, 0x3c,
	0xeb, 0x15, 0xd9, 0x39, 0x8f, 0xa9, 0x01, 0xe8, 0xb9, 0x84, 0x97, 0xf5, 0xcb, 0x6b, 0xfe, 0x0a,
	0x00, 0x00, 0xff, 0xff, 0x8c, 0x3e, 0xc5, 0x8e, 0x11, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Swapped queries the current swapped status that includes a burnt amount of from-coin and a minted amount of
	// to-coin.
	Swapped(ctx context.Context, in *QuerySwappedRequest, opts ...grpc.CallOption) (*QuerySwappedResponse, error)
	// TotalSwappableToCoinAmount queries the current swappable amount for to-coin.
	TotalSwappableToCoinAmount(ctx context.Context, in *QueryTotalSwappableToCoinAmountRequest, opts ...grpc.CallOption) (*QueryTotalSwappableToCoinAmountResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Swapped(ctx context.Context, in *QuerySwappedRequest, opts ...grpc.CallOption) (*QuerySwappedResponse, error) {
	out := new(QuerySwappedResponse)
	err := c.cc.Invoke(ctx, "/lbm.fswap.v1.Query/Swapped", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) TotalSwappableToCoinAmount(ctx context.Context, in *QueryTotalSwappableToCoinAmountRequest, opts ...grpc.CallOption) (*QueryTotalSwappableToCoinAmountResponse, error) {
	out := new(QueryTotalSwappableToCoinAmountResponse)
	err := c.cc.Invoke(ctx, "/lbm.fswap.v1.Query/TotalSwappableToCoinAmount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Swapped queries the current swapped status that includes a burnt amount of from-coin and a minted amount of
	// to-coin.
	Swapped(context.Context, *QuerySwappedRequest) (*QuerySwappedResponse, error)
	// TotalSwappableToCoinAmount queries the current swappable amount for to-coin.
	TotalSwappableToCoinAmount(context.Context, *QueryTotalSwappableToCoinAmountRequest) (*QueryTotalSwappableToCoinAmountResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Swapped(ctx context.Context, req *QuerySwappedRequest) (*QuerySwappedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Swapped not implemented")
}
func (*UnimplementedQueryServer) TotalSwappableToCoinAmount(ctx context.Context, req *QueryTotalSwappableToCoinAmountRequest) (*QueryTotalSwappableToCoinAmountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TotalSwappableToCoinAmount not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Swapped_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySwappedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Swapped(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbm.fswap.v1.Query/Swapped",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Swapped(ctx, req.(*QuerySwappedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TotalSwappableToCoinAmount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTotalSwappableToCoinAmountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TotalSwappableToCoinAmount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbm.fswap.v1.Query/TotalSwappableToCoinAmount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TotalSwappableToCoinAmount(ctx, req.(*QueryTotalSwappableToCoinAmountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "lbm.fswap.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Swapped",
			Handler:    _Query_Swapped_Handler,
		},
		{
			MethodName: "TotalSwappableToCoinAmount",
			Handler:    _Query_TotalSwappableToCoinAmount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbm/fswap/v1/query.proto",
}

func (m *QuerySwappedRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySwappedRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySwappedRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QuerySwappedResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySwappedResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySwappedResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.ToCoinAmount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.FromCoinAmount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryTotalSwappableToCoinAmountRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTotalSwappableToCoinAmountRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTotalSwappableToCoinAmountRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryTotalSwappableToCoinAmountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTotalSwappableToCoinAmountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTotalSwappableToCoinAmountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.SwappableAmount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QuerySwappedRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QuerySwappedResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.FromCoinAmount.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.ToCoinAmount.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryTotalSwappableToCoinAmountRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryTotalSwappableToCoinAmountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.SwappableAmount.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QuerySwappedRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QuerySwappedRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySwappedRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QuerySwappedResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QuerySwappedResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySwappedResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromCoinAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.FromCoinAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToCoinAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ToCoinAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTotalSwappableToCoinAmountRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTotalSwappableToCoinAmountRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTotalSwappableToCoinAmountRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTotalSwappableToCoinAmountResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTotalSwappableToCoinAmountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTotalSwappableToCoinAmountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SwappableAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SwappableAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
