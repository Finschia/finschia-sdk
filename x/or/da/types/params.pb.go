// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: finschia/or/da/v1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

// Params defines the parameters for the module.
type Params struct {
	// 1. CC-related
	CCBatchMaxBytes uint64 `protobuf:"varint,1,opt,name=cc_batch_max_bytes,json=ccBatchMaxBytes,proto3" json:"cc_batch_max_bytes,omitempty"`
	// max_queue_tx_size is the maximum queue tx size that can be submitted.
	MaxQueueTxSize uint64 `protobuf:"varint,2,opt,name=max_queue_tx_size,json=maxQueueTxSize,proto3" json:"max_queue_tx_size,omitempty"`
	// min_queue_tx_gas is the minimum gas that must be specified for a queue tx.
	MinQueueTxGas uint64 `protobuf:"varint,3,opt,name=min_queue_tx_gas,json=minQueueTxGas,proto3" json:"min_queue_tx_gas,omitempty"`
	// l2gas_discount_divisor is the ratio between the cost of gas on L1 and L2.
	// This is a positive integer, meaning we assume L2 gas is always less costly.
	L2GasDiscountDivisor uint64 `protobuf:"varint,4,opt,name=l2gas_discount_divisor,json=l2gasDiscountDivisor,proto3" json:"l2gas_discount_divisor,omitempty"`
	// enqueue_l2gas_prepaid is the base cost of calling enqueue function.
	EnqueueL2GasPrepaid uint64 `protobuf:"varint,5,opt,name=enqueue_l2gas_prepaid,json=enqueueL2gasPrepaid,proto3" json:"enqueue_l2gas_prepaid,omitempty"`
	// A sequencer must submit a queue tx to L2 before this time.
	QueueTxExpirationWindow uint64 `protobuf:"varint,6,opt,name=queue_tx_expiration_window,json=queueTxExpirationWindow,proto3" json:"queue_tx_expiration_window,omitempty"`
	// 2. SCC-related
	SCCBatchMaxBytes uint64 `protobuf:"varint,7,opt,name=scc_batch_max_bytes,json=sccBatchMaxBytes,proto3" json:"scc_batch_max_bytes,omitempty"`
	// Number of seconds that the verifier is allowed to submit a fraud proof.
	// Currnet scc batch header timestamp + fraud_proof_window(seconds) = challenge period
	FraudProofWindow int64 `protobuf:"varint,8,opt,name=fraud_proof_window,json=fraudProofWindow,proto3" json:"fraud_proof_window,omitempty"`
	// Number of seconds that the sequencer is exclusively allowed to post state roots
	SequencerPublishWindow int64 `protobuf:"varint,9,opt,name=sequencer_publish_window,json=sequencerPublishWindow,proto3" json:"sequencer_publish_window,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_208dad78ec5197d3, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetCCBatchMaxBytes() uint64 {
	if m != nil {
		return m.CCBatchMaxBytes
	}
	return 0
}

func (m *Params) GetMaxQueueTxSize() uint64 {
	if m != nil {
		return m.MaxQueueTxSize
	}
	return 0
}

func (m *Params) GetMinQueueTxGas() uint64 {
	if m != nil {
		return m.MinQueueTxGas
	}
	return 0
}

func (m *Params) GetL2GasDiscountDivisor() uint64 {
	if m != nil {
		return m.L2GasDiscountDivisor
	}
	return 0
}

func (m *Params) GetEnqueueL2GasPrepaid() uint64 {
	if m != nil {
		return m.EnqueueL2GasPrepaid
	}
	return 0
}

func (m *Params) GetQueueTxExpirationWindow() uint64 {
	if m != nil {
		return m.QueueTxExpirationWindow
	}
	return 0
}

func (m *Params) GetSCCBatchMaxBytes() uint64 {
	if m != nil {
		return m.SCCBatchMaxBytes
	}
	return 0
}

func (m *Params) GetFraudProofWindow() int64 {
	if m != nil {
		return m.FraudProofWindow
	}
	return 0
}

func (m *Params) GetSequencerPublishWindow() int64 {
	if m != nil {
		return m.SequencerPublishWindow
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "finschia.or.da.v1.Params")
}

func init() { proto.RegisterFile("finschia/or/da/v1/params.proto", fileDescriptor_208dad78ec5197d3) }

var fileDescriptor_208dad78ec5197d3 = []byte{
	// 466 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x1b, 0x5a, 0x0a, 0x58, 0x82, 0x75, 0x6e, 0x19, 0x51, 0x0f, 0xe9, 0xc4, 0x85, 0x21,
	0x41, 0xac, 0x0d, 0x0e, 0x08, 0x2e, 0xa8, 0x1d, 0x3f, 0x0e, 0x20, 0x95, 0x0e, 0x09, 0x89, 0x8b,
	0xe5, 0x38, 0x6e, 0x6a, 0xd1, 0xc4, 0xae, 0xed, 0x74, 0xd9, 0xfe, 0x0a, 0x8e, 0x1c, 0xf9, 0x73,
	0x38, 0xee, 0xc8, 0x69, 0x42, 0xe9, 0x3f, 0xc1, 0x11, 0xc5, 0x4e, 0x8a, 0x98, 0xb8, 0x59, 0xef,
	0xf3, 0xfd, 0x3c, 0x3f, 0x3d, 0x3d, 0x10, 0xcc, 0x79, 0xa6, 0xe9, 0x82, 0x13, 0x24, 0x14, 0x8a,
	0x09, 0x5a, 0x1f, 0x22, 0x49, 0x14, 0x49, 0x75, 0x28, 0x95, 0x30, 0x02, 0xee, 0x36, 0x3c, 0x14,
	0x2a, 0x8c, 0x49, 0xb8, 0x3e, 0x1c, 0x8e, 0x12, 0x21, 0x92, 0x25, 0x43, 0x36, 0x10, 0xe5, 0x73,
	0x64, 0x78, 0xca, 0xb4, 0x21, 0xa9, 0x74, 0xce, 0x70, 0x90, 0x88, 0x44, 0xd8, 0x27, 0xaa, 0x5e,
	0xae, 0x7a, 0xff, 0x77, 0x1b, 0x74, 0xa7, 0xb6, 0x35, 0x7c, 0x09, 0x20, 0xa5, 0x38, 0x22, 0x86,
	0x2e, 0x70, 0x4a, 0x0a, 0x1c, 0x9d, 0x19, 0xa6, 0x7d, 0x6f, 0xdf, 0x3b, 0xe8, 0x8c, 0xfb, 0xe5,
	0xe5, 0x68, 0x67, 0x32, 0x19, 0x57, 0xf0, 0x3d, 0x29, 0xc6, 0x15, 0x9a, 0xed, 0x50, 0xfa, 0x4f,
	0x01, 0x3e, 0x04, 0xbb, 0x95, 0xb8, 0xca, 0x59, 0xce, 0xb0, 0x29, 0xb0, 0xe6, 0xe7, 0xcc, 0xbf,
	0x56, 0x35, 0x98, 0xdd, 0x49, 0x49, 0xf1, 0xa1, 0xaa, 0x7f, 0x2c, 0x4e, 0xf8, 0x39, 0x83, 0x0f,
	0x40, 0x2f, 0xe5, 0xd9, 0xdf, 0x68, 0x42, 0xb4, 0xdf, 0xb6, 0xc9, 0xdb, 0x29, 0xcf, 0xea, 0xe4,
	0x1b, 0xa2, 0xe1, 0x53, 0xb0, 0xb7, 0x3c, 0x4a, 0x88, 0xc6, 0x31, 0xd7, 0x54, 0xe4, 0x99, 0xc1,
	0x31, 0x5f, 0x73, 0x2d, 0x94, 0xdf, 0xb1, 0xf1, 0x81, 0xa5, 0xc7, 0x35, 0x3c, 0x76, 0x0c, 0x1e,
	0x81, 0xbb, 0x2c, 0x73, 0xcd, 0x9d, 0x2d, 0x15, 0x93, 0x84, 0xc7, 0xfe, 0x75, 0x2b, 0xf5, 0x6b,
	0xf8, 0xae, 0x62, 0x53, 0x87, 0xe0, 0x0b, 0x30, 0xdc, 0x8e, 0xc3, 0x0a, 0xc9, 0x15, 0x31, 0x5c,
	0x64, 0xf8, 0x94, 0x67, 0xb1, 0x38, 0xf5, 0xbb, 0x56, 0xbc, 0xb7, 0x72, 0x93, 0xbd, 0xda, 0xf2,
	0x4f, 0x16, 0xc3, 0x09, 0xe8, 0xeb, 0xff, 0x6c, 0xef, 0x86, 0xdd, 0xde, 0xa0, 0xbc, 0x1c, 0xf5,
	0x4e, 0xae, 0xae, 0xaf, 0xa7, 0xaf, 0xee, 0xef, 0x11, 0x80, 0x73, 0x45, 0xf2, 0x18, 0x4b, 0x25,
	0xc4, 0xbc, 0xf9, 0xf9, 0xe6, 0xbe, 0x77, 0xd0, 0x9e, 0xf5, 0x2c, 0x99, 0x56, 0xa0, 0xfe, 0xf2,
	0x19, 0xf0, 0x35, 0x5b, 0xe5, 0x2c, 0xa3, 0x4c, 0x61, 0x99, 0x47, 0x4b, 0xae, 0x17, 0x8d, 0x73,
	0xcb, 0x3a, 0x7b, 0x5b, 0x3e, 0x75, 0xd8, 0x99, 0xcf, 0x3b, 0xdf, 0xbe, 0x8f, 0x5a, 0xe3, 0xb7,
	0x3f, 0xca, 0xc0, 0xbb, 0x28, 0x03, 0xef, 0x57, 0x19, 0x78, 0x5f, 0x37, 0x41, 0xeb, 0x62, 0x13,
	0xb4, 0x7e, 0x6e, 0x82, 0xd6, 0xe7, 0x30, 0xe1, 0x66, 0x91, 0x47, 0x21, 0x15, 0x29, 0x7a, 0xdd,
	0x5c, 0x62, 0x73, 0x72, 0x8f, 0x75, 0xfc, 0x05, 0x15, 0xf5, 0x61, 0x9a, 0x33, 0xc9, 0x74, 0xd4,
	0xb5, 0xb7, 0xf4, 0xe4, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x85, 0xda, 0xaf, 0xb7, 0x02,
	0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SequencerPublishWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SequencerPublishWindow))
		i--
		dAtA[i] = 0x48
	}
	if m.FraudProofWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.FraudProofWindow))
		i--
		dAtA[i] = 0x40
	}
	if m.SCCBatchMaxBytes != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SCCBatchMaxBytes))
		i--
		dAtA[i] = 0x38
	}
	if m.QueueTxExpirationWindow != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.QueueTxExpirationWindow))
		i--
		dAtA[i] = 0x30
	}
	if m.EnqueueL2GasPrepaid != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.EnqueueL2GasPrepaid))
		i--
		dAtA[i] = 0x28
	}
	if m.L2GasDiscountDivisor != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.L2GasDiscountDivisor))
		i--
		dAtA[i] = 0x20
	}
	if m.MinQueueTxGas != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MinQueueTxGas))
		i--
		dAtA[i] = 0x18
	}
	if m.MaxQueueTxSize != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxQueueTxSize))
		i--
		dAtA[i] = 0x10
	}
	if m.CCBatchMaxBytes != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.CCBatchMaxBytes))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CCBatchMaxBytes != 0 {
		n += 1 + sovParams(uint64(m.CCBatchMaxBytes))
	}
	if m.MaxQueueTxSize != 0 {
		n += 1 + sovParams(uint64(m.MaxQueueTxSize))
	}
	if m.MinQueueTxGas != 0 {
		n += 1 + sovParams(uint64(m.MinQueueTxGas))
	}
	if m.L2GasDiscountDivisor != 0 {
		n += 1 + sovParams(uint64(m.L2GasDiscountDivisor))
	}
	if m.EnqueueL2GasPrepaid != 0 {
		n += 1 + sovParams(uint64(m.EnqueueL2GasPrepaid))
	}
	if m.QueueTxExpirationWindow != 0 {
		n += 1 + sovParams(uint64(m.QueueTxExpirationWindow))
	}
	if m.SCCBatchMaxBytes != 0 {
		n += 1 + sovParams(uint64(m.SCCBatchMaxBytes))
	}
	if m.FraudProofWindow != 0 {
		n += 1 + sovParams(uint64(m.FraudProofWindow))
	}
	if m.SequencerPublishWindow != 0 {
		n += 1 + sovParams(uint64(m.SequencerPublishWindow))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CCBatchMaxBytes", wireType)
			}
			m.CCBatchMaxBytes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CCBatchMaxBytes |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxQueueTxSize", wireType)
			}
			m.MaxQueueTxSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxQueueTxSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinQueueTxGas", wireType)
			}
			m.MinQueueTxGas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinQueueTxGas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field L2GasDiscountDivisor", wireType)
			}
			m.L2GasDiscountDivisor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.L2GasDiscountDivisor |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnqueueL2GasPrepaid", wireType)
			}
			m.EnqueueL2GasPrepaid = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EnqueueL2GasPrepaid |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueueTxExpirationWindow", wireType)
			}
			m.QueueTxExpirationWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.QueueTxExpirationWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SCCBatchMaxBytes", wireType)
			}
			m.SCCBatchMaxBytes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SCCBatchMaxBytes |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FraudProofWindow", wireType)
			}
			m.FraudProofWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FraudProofWindow |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SequencerPublishWindow", wireType)
			}
			m.SequencerPublishWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SequencerPublishWindow |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
