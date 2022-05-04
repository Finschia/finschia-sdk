// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lbm/foundation/v1/genesis.proto

package foundation

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/line/lbm-sdk/codec/types"
	_ "github.com/regen-network/cosmos-proto"
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

// GenesisState defines the foundation module's genesis state.
type GenesisState struct {
	// params defines the module parameters at genesis.
	Params *Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	// allowed_validators defines the allowed validator addresses at genesis.
	// provided empty, the module gathers information from staking module.
	ValidatorAuths []ValidatorAuth `protobuf:"bytes,2,rep,name=validator_auths,json=validatorAuths,proto3" json:"validator_auths"`
	// foundation is the foundation info.
	Foundation *FoundationInfo `protobuf:"bytes,3,opt,name=foundation,proto3" json:"foundation,omitempty"`
	// members is the list of the foundation members.
	Members []Member `protobuf:"bytes,4,rep,name=members,proto3" json:"members"`
	// it is used to get the next proposal ID.
	PreviousProposalId uint64 `protobuf:"varint,5,opt,name=previous_proposal_id,json=previousProposalId,proto3" json:"previous_proposal_id,omitempty"`
	// proposals is the list of proposals.
	Proposals []Proposal `protobuf:"bytes,6,rep,name=proposals,proto3" json:"proposals"`
	// votes is the list of votes.
	Votes []Vote `protobuf:"bytes,7,rep,name=votes,proto3" json:"votes"`
	// grants
	Authorizations []GrantAuthorization `protobuf:"bytes,8,rep,name=authorizations,proto3" json:"authorizations"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5e13dd78b24d473, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

// GrantAuthorization defines authorization grant to grantee via route.
type GrantAuthorization struct {
	Route         string     `protobuf:"bytes,1,opt,name=route,proto3" json:"route,omitempty"`
	Grantee       string     `protobuf:"bytes,2,opt,name=grantee,proto3" json:"grantee,omitempty"`
	Authorization *types.Any `protobuf:"bytes,3,opt,name=authorization,proto3" json:"authorization,omitempty"`
}

func (m *GrantAuthorization) Reset()         { *m = GrantAuthorization{} }
func (m *GrantAuthorization) String() string { return proto.CompactTextString(m) }
func (*GrantAuthorization) ProtoMessage()    {}
func (*GrantAuthorization) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5e13dd78b24d473, []int{1}
}
func (m *GrantAuthorization) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GrantAuthorization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GrantAuthorization.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GrantAuthorization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrantAuthorization.Merge(m, src)
}
func (m *GrantAuthorization) XXX_Size() int {
	return m.Size()
}
func (m *GrantAuthorization) XXX_DiscardUnknown() {
	xxx_messageInfo_GrantAuthorization.DiscardUnknown(m)
}

var xxx_messageInfo_GrantAuthorization proto.InternalMessageInfo

func (m *GrantAuthorization) GetRoute() string {
	if m != nil {
		return m.Route
	}
	return ""
}

func (m *GrantAuthorization) GetGrantee() string {
	if m != nil {
		return m.Grantee
	}
	return ""
}

func (m *GrantAuthorization) GetAuthorization() *types.Any {
	if m != nil {
		return m.Authorization
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "lbm.foundation.v1.GenesisState")
	proto.RegisterType((*GrantAuthorization)(nil), "lbm.foundation.v1.GrantAuthorization")
}

func init() { proto.RegisterFile("lbm/foundation/v1/genesis.proto", fileDescriptor_c5e13dd78b24d473) }

var fileDescriptor_c5e13dd78b24d473 = []byte{
	// 505 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x31, 0x6f, 0xd3, 0x40,
	0x1c, 0xc5, 0xed, 0x36, 0x4d, 0xda, 0x2b, 0x14, 0x71, 0x8a, 0xc4, 0xb5, 0x48, 0x4e, 0x88, 0x40,
	0xca, 0x92, 0x3b, 0xd2, 0x4e, 0xb0, 0xa0, 0x64, 0xa0, 0xea, 0x80, 0xa8, 0x5c, 0x89, 0x81, 0x25,
	0x3a, 0xd7, 0x17, 0xc7, 0xc2, 0xf6, 0x3f, 0xf2, 0x9d, 0x2d, 0xda, 0x4f, 0xc0, 0xc8, 0xc6, 0xda,
	0x91, 0x1d, 0x3e, 0x44, 0xc5, 0xd4, 0x91, 0x09, 0xa1, 0x64, 0xe1, 0x63, 0xa0, 0xdc, 0xd9, 0x24,
	0x51, 0xcc, 0x96, 0xff, 0xfd, 0x7f, 0xef, 0xdd, 0xcb, 0xf3, 0xa1, 0x56, 0xe4, 0xc5, 0x6c, 0x0c,
	0x59, 0xe2, 0x73, 0x15, 0x42, 0xc2, 0xf2, 0x3e, 0x0b, 0x44, 0x22, 0x64, 0x28, 0xe9, 0x34, 0x05,
	0x05, 0xf8, 0x61, 0xe4, 0xc5, 0x74, 0x09, 0xd0, 0xbc, 0x7f, 0xd4, 0x0c, 0x20, 0x00, 0xbd, 0x65,
	0x8b, 0x5f, 0x06, 0x3c, 0xea, 0x6c, 0x3a, 0xad, 0xc8, 0x0c, 0x73, 0x78, 0x09, 0x32, 0x06, 0x39,
	0x32, 0x62, 0x33, 0x94, 0xab, 0x00, 0x20, 0x88, 0x04, 0xd3, 0x93, 0x97, 0x8d, 0x19, 0x4f, 0xae,
	0xcc, 0xaa, 0xf3, 0xa5, 0x86, 0xee, 0x9d, 0x9a, 0x50, 0x17, 0x8a, 0x2b, 0x81, 0xfb, 0xa8, 0x3e,
	0xe5, 0x29, 0x8f, 0x25, 0xb1, 0xdb, 0x76, 0x77, 0xff, 0xf8, 0x90, 0x6e, 0x84, 0xa4, 0xe7, 0x1a,
	0x70, 0x0b, 0x10, 0xbf, 0x45, 0x0f, 0x72, 0x1e, 0x85, 0x3e, 0x57, 0x90, 0x8e, 0x78, 0xa6, 0x26,
	0x92, 0x6c, 0xb5, 0xb7, 0xbb, 0xfb, 0xc7, 0xed, 0x0a, 0xed, 0xbb, 0x92, 0x1c, 0x64, 0x6a, 0x32,
	0xac, 0xdd, 0xfe, 0x6a, 0x59, 0xee, 0x41, 0xbe, 0x7a, 0x28, 0xf1, 0x00, 0xa1, 0xa5, 0x88, 0x6c,
	0xeb, 0x1c, 0x4f, 0x2a, 0xbc, 0x5e, 0xff, 0x9b, 0xce, 0x92, 0x31, 0xb8, 0x2b, 0x22, 0xfc, 0x02,
	0x35, 0x62, 0x11, 0x7b, 0x22, 0x95, 0xa4, 0xa6, 0xb3, 0x54, 0xfd, 0x8f, 0x37, 0x9a, 0x28, 0x42,
	0x94, 0x3c, 0x7e, 0x8e, 0x9a, 0xd3, 0x54, 0xe4, 0x21, 0x64, 0xba, 0xcc, 0x29, 0x48, 0x1e, 0x8d,
	0x42, 0x9f, 0xec, 0xb4, 0xed, 0x6e, 0xcd, 0xc5, 0xe5, 0xee, 0xbc, 0x58, 0x9d, 0xf9, 0xf8, 0x15,
	0xda, 0x2b, 0x41, 0x49, 0xea, 0xfa, 0xba, 0xc7, 0x55, 0xb5, 0x15, 0x4c, 0x71, 0xe1, 0x52, 0x83,
	0x4f, 0xd0, 0x4e, 0x0e, 0x4a, 0x48, 0xd2, 0xd0, 0xe2, 0x47, 0x55, 0xbd, 0x81, 0x12, 0x85, 0xd0,
	0xb0, 0xf8, 0x02, 0x1d, 0x2c, 0xca, 0x86, 0x34, 0xbc, 0xd6, 0x94, 0x24, 0xbb, 0x5a, 0xfd, 0xac,
	0x42, 0x7d, 0x9a, 0xf2, 0x44, 0x0d, 0x56, 0xe9, 0xb2, 0xfa, 0x75, 0x8b, 0x97, 0xbb, 0x9f, 0x6e,
	0x5a, 0xd6, 0x9f, 0x9b, 0x96, 0xd5, 0xf9, 0x66, 0x23, 0xbc, 0x29, 0xc3, 0x4d, 0xb4, 0x93, 0x42,
	0xa6, 0x84, 0x7e, 0x1e, 0x7b, 0xae, 0x19, 0x30, 0x41, 0x8d, 0x60, 0xc1, 0x0a, 0x41, 0xb6, 0xf4,
	0x79, 0x39, 0xe2, 0x18, 0xdd, 0x5f, 0xbb, 0xa2, 0xf8, 0x9c, 0x4d, 0x6a, 0xde, 0x24, 0x2d, 0xdf,
	0x24, 0x1d, 0x24, 0x57, 0xc3, 0xfe, 0x8f, 0xef, 0xbd, 0x5e, 0x10, 0xaa, 0x49, 0xe6, 0xd1, 0x4b,
	0x88, 0x59, 0x14, 0x26, 0x82, 0x45, 0x5e, 0xdc, 0x93, 0xfe, 0x07, 0xf6, 0x91, 0x2d, 0x8c, 0xae,
	0xe9, 0x5a, 0x1e, 0x77, 0xdd, 0x7d, 0x38, 0xfc, 0x3a, 0x73, 0xec, 0xdb, 0x99, 0x63, 0xdf, 0xcd,
	0x1c, 0xfb, 0xf7, 0xcc, 0xb1, 0x3f, 0xcf, 0x1d, 0xeb, 0x6e, 0xee, 0x58, 0x3f, 0xe7, 0x8e, 0xf5,
	0xfe, 0xe9, 0xff, 0xbd, 0x97, 0x7d, 0x79, 0x75, 0x9d, 0xe9, 0xe4, 0x6f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x85, 0x0c, 0xc1, 0x92, 0xc0, 0x03, 0x00, 0x00,
}

func (this *GrantAuthorization) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GrantAuthorization)
	if !ok {
		that2, ok := that.(GrantAuthorization)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Route != that1.Route {
		return false
	}
	if this.Grantee != that1.Grantee {
		return false
	}
	if !this.Authorization.Equal(that1.Authorization) {
		return false
	}
	return true
}
func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Authorizations) > 0 {
		for iNdEx := len(m.Authorizations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Authorizations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.Votes) > 0 {
		for iNdEx := len(m.Votes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Votes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.Proposals) > 0 {
		for iNdEx := len(m.Proposals) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Proposals[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.PreviousProposalId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.PreviousProposalId))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Members) > 0 {
		for iNdEx := len(m.Members) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Members[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Foundation != nil {
		{
			size, err := m.Foundation.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ValidatorAuths) > 0 {
		for iNdEx := len(m.ValidatorAuths) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorAuths[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GrantAuthorization) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GrantAuthorization) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GrantAuthorization) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Authorization != nil {
		{
			size, err := m.Authorization.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Grantee) > 0 {
		i -= len(m.Grantee)
		copy(dAtA[i:], m.Grantee)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Grantee)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Route) > 0 {
		i -= len(m.Route)
		copy(dAtA[i:], m.Route)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Route)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.ValidatorAuths) > 0 {
		for _, e := range m.ValidatorAuths {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.Foundation != nil {
		l = m.Foundation.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Members) > 0 {
		for _, e := range m.Members {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.PreviousProposalId != 0 {
		n += 1 + sovGenesis(uint64(m.PreviousProposalId))
	}
	if len(m.Proposals) > 0 {
		for _, e := range m.Proposals {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Votes) > 0 {
		for _, e := range m.Votes {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Authorizations) > 0 {
		for _, e := range m.Authorizations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *GrantAuthorization) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Route)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Grantee)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Authorization != nil {
		l = m.Authorization.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &Params{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAuths", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAuths = append(m.ValidatorAuths, ValidatorAuth{})
			if err := m.ValidatorAuths[len(m.ValidatorAuths)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Foundation", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Foundation == nil {
				m.Foundation = &FoundationInfo{}
			}
			if err := m.Foundation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Members", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Members = append(m.Members, Member{})
			if err := m.Members[len(m.Members)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousProposalId", wireType)
			}
			m.PreviousProposalId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PreviousProposalId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proposals", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Proposals = append(m.Proposals, Proposal{})
			if err := m.Proposals[len(m.Proposals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Votes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Votes = append(m.Votes, Vote{})
			if err := m.Votes[len(m.Votes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authorizations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authorizations = append(m.Authorizations, GrantAuthorization{})
			if err := m.Authorizations[len(m.Authorizations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GrantAuthorization) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GrantAuthorization: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GrantAuthorization: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Route", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Route = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Grantee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Grantee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authorization", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Authorization == nil {
				m.Authorization = &types.Any{}
			}
			if err := m.Authorization.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
