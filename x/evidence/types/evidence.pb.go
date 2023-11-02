// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/evidence/v1beta1/evidence.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Equivocation implements the Evidence interface and defines evidence of double
// signing misbehavior.
type Equivocation struct {
	Height           int64     `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Time             time.Time `protobuf:"bytes,2,opt,name=time,proto3,stdtime" json:"time"`
	Power            int64     `protobuf:"varint,3,opt,name=power,proto3" json:"power,omitempty"`
	ConsensusAddress string    `protobuf:"bytes,4,opt,name=consensus_address,json=consensusAddress,proto3" json:"consensus_address,omitempty" yaml:"consensus_address"`
}

func (m *Equivocation) Reset()      { *m = Equivocation{} }
func (*Equivocation) ProtoMessage() {}
func (*Equivocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd143e71a177f0dd, []int{0}
}
func (m *Equivocation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Equivocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Equivocation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Equivocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Equivocation.Merge(m, src)
}
func (m *Equivocation) XXX_Size() int {
	return m.Size()
}
func (m *Equivocation) XXX_DiscardUnknown() {
	xxx_messageInfo_Equivocation.DiscardUnknown(m)
}

var xxx_messageInfo_Equivocation proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Equivocation)(nil), "cosmos.evidence.v1beta1.Equivocation")
}

func init() {
	proto.RegisterFile("cosmos/evidence/v1beta1/evidence.proto", fileDescriptor_dd143e71a177f0dd)
}

var fileDescriptor_dd143e71a177f0dd = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x3f, 0x4f, 0xfa, 0x40,
	0x1c, 0xc6, 0xef, 0x7e, 0xf0, 0x23, 0x5a, 0x19, 0xb4, 0x21, 0xda, 0x10, 0x73, 0x47, 0x18, 0x0c,
	0x8b, 0xbd, 0xa0, 0x8b, 0x61, 0x93, 0x44, 0x13, 0x27, 0x13, 0xe2, 0xe4, 0x62, 0xfa, 0xe7, 0x68,
	0x2f, 0xd2, 0x7e, 0x2b, 0x77, 0x45, 0x79, 0x07, 0x8e, 0x8c, 0x8e, 0x8c, 0xbe, 0x14, 0x36, 0x19,
	0x9d, 0xd0, 0x94, 0xc5, 0xd9, 0x57, 0x60, 0xe8, 0x41, 0x1d, 0xdc, 0xbe, 0xcf, 0x93, 0xcf, 0xf3,
	0x49, 0x2e, 0x67, 0x1c, 0x79, 0x20, 0x23, 0x90, 0x8c, 0x8f, 0x84, 0xcf, 0x63, 0x8f, 0xb3, 0x51,
	0xdb, 0xe5, 0xca, 0x69, 0x17, 0x85, 0x9d, 0x0c, 0x41, 0x81, 0x79, 0xa0, 0x39, 0xbb, 0xa8, 0xd7,
	0x5c, 0xbd, 0x16, 0x40, 0x00, 0x39, 0xc3, 0x56, 0x97, 0xc6, 0xeb, 0x34, 0x00, 0x08, 0x06, 0x9c,
	0xe5, 0xc9, 0x4d, 0xfb, 0x4c, 0x89, 0x88, 0x4b, 0xe5, 0x44, 0x89, 0x06, 0x9a, 0x6f, 0xd8, 0xa8,
	0x5e, 0x3c, 0xa4, 0x62, 0x04, 0x9e, 0xa3, 0x04, 0xc4, 0xe6, 0xbe, 0x51, 0x09, 0xb9, 0x08, 0x42,
	0x65, 0xe1, 0x06, 0x6e, 0x95, 0x7a, 0xeb, 0x64, 0x9e, 0x19, 0xe5, 0xd5, 0xd6, 0xfa, 0xd7, 0xc0,
	0xad, 0x9d, 0x93, 0xba, 0xad, 0xc5, 0xf6, 0x46, 0x6c, 0xdf, 0x6c, 0xc4, 0xdd, 0xad, 0xd9, 0x82,
	0xa2, 0xc9, 0x07, 0xc5, 0xbd, 0x7c, 0x61, 0xd6, 0x8c, 0xff, 0x09, 0x3c, 0xf2, 0xa1, 0x55, 0xca,
	0x85, 0x3a, 0x98, 0x57, 0xc6, 0x9e, 0x07, 0xb1, 0xe4, 0xb1, 0x4c, 0xe5, 0x9d, 0xe3, 0xfb, 0x43,
	0x2e, 0xa5, 0x55, 0x6e, 0xe0, 0xd6, 0x76, 0xf7, 0xf0, 0x7b, 0x41, 0xad, 0xb1, 0x13, 0x0d, 0x3a,
	0xcd, 0x3f, 0x48, 0xb3, 0xb7, 0x5b, 0x74, 0xe7, 0xba, 0xea, 0x54, 0x9f, 0xa7, 0x14, 0xbd, 0x4c,
	0x29, 0xfa, 0x9a, 0x52, 0xd4, 0xbd, 0x7e, 0xcd, 0x08, 0x9e, 0x65, 0x04, 0xcf, 0x33, 0x82, 0x3f,
	0x33, 0x82, 0x27, 0x4b, 0x82, 0xe6, 0x4b, 0x82, 0xde, 0x97, 0x04, 0xdd, 0xb6, 0x03, 0xa1, 0xc2,
	0xd4, 0xb5, 0x3d, 0x88, 0xd8, 0xa5, 0x88, 0xa5, 0x17, 0x0a, 0x87, 0xf5, 0xd7, 0xc7, 0xb1, 0xf4,
	0xef, 0xd9, 0xd3, 0xef, 0x1f, 0xa8, 0x71, 0xc2, 0xa5, 0x5b, 0xc9, 0xdf, 0x78, 0xfa, 0x13, 0x00,
	0x00, 0xff, 0xff, 0x54, 0xd6, 0x7c, 0xfe, 0xa3, 0x01, 0x00, 0x00,
}

func (m *Equivocation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Equivocation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Equivocation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ConsensusAddress) > 0 {
		i -= len(m.ConsensusAddress)
		copy(dAtA[i:], m.ConsensusAddress)
		i = encodeVarintEvidence(dAtA, i, uint64(len(m.ConsensusAddress)))
		i--
		dAtA[i] = 0x22
	}
	if m.Power != 0 {
		i = encodeVarintEvidence(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x18
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Time, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintEvidence(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x12
	if m.Height != 0 {
		i = encodeVarintEvidence(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvidence(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvidence(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Equivocation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovEvidence(uint64(m.Height))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovEvidence(uint64(l))
	if m.Power != 0 {
		n += 1 + sovEvidence(uint64(m.Power))
	}
	l = len(m.ConsensusAddress)
	if l > 0 {
		n += 1 + l + sovEvidence(uint64(l))
	}
	return n
}

func sovEvidence(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvidence(x uint64) (n int) {
	return sovEvidence(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Equivocation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvidence
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
			return fmt.Errorf("proto: Equivocation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Equivocation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
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
				return ErrInvalidLengthEvidence
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvidence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsensusAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvidence
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
				return ErrInvalidLengthEvidence
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvidence
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConsensusAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvidence(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvidence
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
func skipEvidence(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvidence
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
					return 0, ErrIntOverflowEvidence
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
					return 0, ErrIntOverflowEvidence
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
				return 0, ErrInvalidLengthEvidence
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvidence
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvidence
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvidence        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvidence          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvidence = fmt.Errorf("proto: unexpected end of group")
)
