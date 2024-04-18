// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lbm/fswap/v1/fswap.proto

package types

import (
	fmt "fmt"
	types "github.com/Finschia/finschia-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type Swapped struct {
	OldCoinAmount types.Coin `protobuf:"bytes,1,opt,name=old_coin_amount,json=oldCoinAmount,proto3,castrepeated=github.com/Finschia/finschia-sdk/types.Coin" json:"old_coin_amount"`
	NewCoinAmount types.Coin `protobuf:"bytes,2,opt,name=new_coin_amount,json=newCoinAmount,proto3,castrepeated=github.com/Finschia/finschia-sdk/types.Coin" json:"new_coin_amount"`
}

func (m *Swapped) Reset()         { *m = Swapped{} }
func (m *Swapped) String() string { return proto.CompactTextString(m) }
func (*Swapped) ProtoMessage()    {}
func (*Swapped) Descriptor() ([]byte, []int) {
	return fileDescriptor_42ca60eaf37a2b67, []int{0}
}
func (m *Swapped) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Swapped) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Swapped.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Swapped) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Swapped.Merge(m, src)
}
func (m *Swapped) XXX_Size() int {
	return m.Size()
}
func (m *Swapped) XXX_DiscardUnknown() {
	xxx_messageInfo_Swapped.DiscardUnknown(m)
}

var xxx_messageInfo_Swapped proto.InternalMessageInfo

func (m *Swapped) GetOldCoinAmount() types.Coin {
	if m != nil {
		return m.OldCoinAmount
	}
	return types.Coin{}
}

func (m *Swapped) GetNewCoinAmount() types.Coin {
	if m != nil {
		return m.NewCoinAmount
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*Swapped)(nil), "lbm.fswap.v1.Swapped")
}

func init() { proto.RegisterFile("lbm/fswap/v1/fswap.proto", fileDescriptor_42ca60eaf37a2b67) }

var fileDescriptor_42ca60eaf37a2b67 = []byte{
	// 262 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xc8, 0x49, 0xca, 0xd5,
	0x4f, 0x2b, 0x2e, 0x4f, 0x2c, 0xd0, 0x2f, 0x33, 0x84, 0x30, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2,
	0x85, 0x78, 0x72, 0x92, 0x72, 0xf5, 0x20, 0x02, 0x65, 0x86, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9,
	0x60, 0x09, 0x7d, 0x10, 0x0b, 0xa2, 0x46, 0x4a, 0x2e, 0x39, 0xbf, 0x38, 0x37, 0xbf, 0x58, 0x3f,
	0x29, 0xb1, 0x38, 0x55, 0xbf, 0xcc, 0x30, 0x29, 0xb5, 0x24, 0xd1, 0x50, 0x3f, 0x39, 0x3f, 0x33,
	0x0f, 0x22, 0xaf, 0xf4, 0x93, 0x91, 0x8b, 0x3d, 0xb8, 0x3c, 0xb1, 0xa0, 0x20, 0x35, 0x45, 0xa8,
	0x8c, 0x8b, 0x3f, 0x3f, 0x27, 0x25, 0x1e, 0x24, 0x1b, 0x9f, 0x98, 0x9b, 0x5f, 0x9a, 0x57, 0x22,
	0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x6d, 0x24, 0xa9, 0x07, 0x31, 0x45, 0x0f, 0x64, 0x8a, 0x1e, 0xd4,
	0x14, 0x3d, 0xe7, 0xfc, 0xcc, 0x3c, 0x27, 0xe3, 0x13, 0xf7, 0xe4, 0x19, 0x56, 0xdd, 0x97, 0xd7,
	0x4e, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x77, 0xcb, 0xcc, 0x2b, 0x4e,
	0xce, 0xc8, 0x4c, 0xd4, 0x4f, 0x83, 0x32, 0x74, 0x8b, 0x53, 0xb2, 0xf5, 0x4b, 0x2a, 0x0b, 0x52,
	0x8b, 0xc1, 0x9a, 0x82, 0x78, 0xf3, 0x73, 0x52, 0x40, 0x0c, 0x47, 0xb0, 0x25, 0x20, 0x7b, 0xf3,
	0x52, 0xcb, 0x51, 0xec, 0x65, 0xa2, 0x8d, 0xbd, 0x79, 0xa9, 0xe5, 0x08, 0x7b, 0x9d, 0x3c, 0x4e,
	0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18,
	0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a, 0x8f, 0xa0, 0xa9, 0x15, 0xd0, 0x28,
	0x01, 0x9b, 0x9e, 0xc4, 0x06, 0x0e, 0x4c, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5c, 0x78,
	0x6f, 0xce, 0xac, 0x01, 0x00, 0x00,
}

func (m *Swapped) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Swapped) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Swapped) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.NewCoinAmount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintFswap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.OldCoinAmount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintFswap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintFswap(dAtA []byte, offset int, v uint64) int {
	offset -= sovFswap(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Swapped) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.OldCoinAmount.Size()
	n += 1 + l + sovFswap(uint64(l))
	l = m.NewCoinAmount.Size()
	n += 1 + l + sovFswap(uint64(l))
	return n
}

func sovFswap(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFswap(x uint64) (n int) {
	return sovFswap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Swapped) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFswap
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
			return fmt.Errorf("proto: Swapped: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Swapped: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OldCoinAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFswap
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
				return ErrInvalidLengthFswap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFswap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OldCoinAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewCoinAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFswap
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
				return ErrInvalidLengthFswap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFswap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.NewCoinAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFswap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFswap
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
func skipFswap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFswap
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
					return 0, ErrIntOverflowFswap
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
					return 0, ErrIntOverflowFswap
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
				return 0, ErrInvalidLengthFswap
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFswap
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFswap
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFswap        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFswap          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFswap = fmt.Errorf("proto: unexpected end of group")
)
