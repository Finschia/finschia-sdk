// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: finschia/or/da/v1/genesis.proto

package types

import (
	fmt "fmt"
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

// GenesisState defines the da module's genesis state.
type GenesisState struct {
	Params  Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	CCList  []CC   `protobuf:"bytes,2,rep,name=cc_list,json=ccList,proto3" json:"cc_list"`
	SCCList []SCC  `protobuf:"bytes,3,rep,name=scc_list,json=sccList,proto3" json:"scc_list"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f212539e5ab9f0a, []int{0}
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

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetCCList() []CC {
	if m != nil {
		return m.CCList
	}
	return nil
}

func (m *GenesisState) GetSCCList() []SCC {
	if m != nil {
		return m.SCCList
	}
	return nil
}

type CC struct {
	RollupName   string        `protobuf:"bytes,1,opt,name=rollup_name,json=rollupName,proto3" json:"rollup_name,omitempty"`
	CCState      CCState       `protobuf:"bytes,2,opt,name=cc_state,json=ccState,proto3" json:"cc_state"`
	History      []CCRef       `protobuf:"bytes,3,rep,name=history,proto3" json:"history"`
	QueueTxState QueueTxState  `protobuf:"bytes,4,opt,name=queue_tx_state,json=queueTxState,proto3" json:"queue_tx_state"`
	QueueList    []L1ToL2Queue `protobuf:"bytes,5,rep,name=queue_list,json=queueList,proto3" json:"queue_list"`
	L2BatchMap   []L2BatchMap  `protobuf:"bytes,6,rep,name=l2_batch_map,json=l2BatchMap,proto3" json:"l2_batch_map"`
}

func (m *CC) Reset()         { *m = CC{} }
func (m *CC) String() string { return proto.CompactTextString(m) }
func (*CC) ProtoMessage()    {}
func (*CC) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f212539e5ab9f0a, []int{1}
}
func (m *CC) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CC.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CC.Merge(m, src)
}
func (m *CC) XXX_Size() int {
	return m.Size()
}
func (m *CC) XXX_DiscardUnknown() {
	xxx_messageInfo_CC.DiscardUnknown(m)
}

var xxx_messageInfo_CC proto.InternalMessageInfo

func (m *CC) GetRollupName() string {
	if m != nil {
		return m.RollupName
	}
	return ""
}

func (m *CC) GetCCState() CCState {
	if m != nil {
		return m.CCState
	}
	return CCState{}
}

func (m *CC) GetHistory() []CCRef {
	if m != nil {
		return m.History
	}
	return nil
}

func (m *CC) GetQueueTxState() QueueTxState {
	if m != nil {
		return m.QueueTxState
	}
	return QueueTxState{}
}

func (m *CC) GetQueueList() []L1ToL2Queue {
	if m != nil {
		return m.QueueList
	}
	return nil
}

func (m *CC) GetL2BatchMap() []L2BatchMap {
	if m != nil {
		return m.L2BatchMap
	}
	return nil
}

type SCC struct {
	RollupName string   `protobuf:"bytes,1,opt,name=rollup_name,json=rollupName,proto3" json:"rollup_name,omitempty"`
	State      SCCState `protobuf:"bytes,2,opt,name=state,proto3" json:"state"`
	History    []SCCRef `protobuf:"bytes,3,rep,name=history,proto3" json:"history"`
}

func (m *SCC) Reset()         { *m = SCC{} }
func (m *SCC) String() string { return proto.CompactTextString(m) }
func (*SCC) ProtoMessage()    {}
func (*SCC) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f212539e5ab9f0a, []int{2}
}
func (m *SCC) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SCC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SCC.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SCC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SCC.Merge(m, src)
}
func (m *SCC) XXX_Size() int {
	return m.Size()
}
func (m *SCC) XXX_DiscardUnknown() {
	xxx_messageInfo_SCC.DiscardUnknown(m)
}

var xxx_messageInfo_SCC proto.InternalMessageInfo

func (m *SCC) GetRollupName() string {
	if m != nil {
		return m.RollupName
	}
	return ""
}

func (m *SCC) GetState() SCCState {
	if m != nil {
		return m.State
	}
	return SCCState{}
}

func (m *SCC) GetHistory() []SCCRef {
	if m != nil {
		return m.History
	}
	return nil
}

type L2BatchMap struct {
	L2Height uint64 `protobuf:"varint,1,opt,name=l2_height,json=l2Height,proto3" json:"l2_height,omitempty"`
	BatchIdx uint64 `protobuf:"varint,2,opt,name=batch_idx,json=batchIdx,proto3" json:"batch_idx,omitempty"`
}

func (m *L2BatchMap) Reset()         { *m = L2BatchMap{} }
func (m *L2BatchMap) String() string { return proto.CompactTextString(m) }
func (*L2BatchMap) ProtoMessage()    {}
func (*L2BatchMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f212539e5ab9f0a, []int{3}
}
func (m *L2BatchMap) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *L2BatchMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_L2BatchMap.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *L2BatchMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_L2BatchMap.Merge(m, src)
}
func (m *L2BatchMap) XXX_Size() int {
	return m.Size()
}
func (m *L2BatchMap) XXX_DiscardUnknown() {
	xxx_messageInfo_L2BatchMap.DiscardUnknown(m)
}

var xxx_messageInfo_L2BatchMap proto.InternalMessageInfo

func (m *L2BatchMap) GetL2Height() uint64 {
	if m != nil {
		return m.L2Height
	}
	return 0
}

func (m *L2BatchMap) GetBatchIdx() uint64 {
	if m != nil {
		return m.BatchIdx
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "finschia.or.da.v1.GenesisState")
	proto.RegisterType((*CC)(nil), "finschia.or.da.v1.CC")
	proto.RegisterType((*SCC)(nil), "finschia.or.da.v1.SCC")
	proto.RegisterType((*L2BatchMap)(nil), "finschia.or.da.v1.L2BatchMap")
}

func init() { proto.RegisterFile("finschia/or/da/v1/genesis.proto", fileDescriptor_5f212539e5ab9f0a) }

var fileDescriptor_5f212539e5ab9f0a = []byte{
	// 523 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x3f, 0x6f, 0xd3, 0x40,
	0x18, 0xc6, 0xe3, 0x24, 0xcd, 0x9f, 0x37, 0x51, 0x11, 0x27, 0x40, 0x69, 0x2a, 0x9c, 0x2a, 0x53,
	0x17, 0x6c, 0xc5, 0x0c, 0x85, 0x85, 0xc1, 0x16, 0xa1, 0x88, 0x80, 0xc0, 0x61, 0x62, 0xb1, 0x2e,
	0xe7, 0x6b, 0x6c, 0xe1, 0xe4, 0x5c, 0xdf, 0xa5, 0x4a, 0x37, 0x3e, 0x02, 0x33, 0x9f, 0xa8, 0x12,
	0x4b, 0x47, 0xa6, 0x0a, 0x25, 0x5f, 0x04, 0xf9, 0xee, 0x4c, 0x23, 0xe2, 0x81, 0x29, 0xf1, 0xfb,
	0x3e, 0xef, 0xef, 0x9e, 0xf7, 0x39, 0x1d, 0x0c, 0x2e, 0xe2, 0x25, 0x27, 0x51, 0x8c, 0x6d, 0x96,
	0xd9, 0x21, 0xb6, 0xaf, 0x46, 0xf6, 0x9c, 0x2e, 0x29, 0x8f, 0xb9, 0x95, 0x66, 0x4c, 0x30, 0xf4,
	0xb0, 0x10, 0x58, 0x2c, 0xb3, 0x42, 0x6c, 0x5d, 0x8d, 0xfa, 0x8f, 0xe6, 0x6c, 0xce, 0x64, 0xd7,
	0xce, 0xff, 0x29, 0x61, 0xdf, 0xdc, 0x27, 0xa5, 0x38, 0xc3, 0x0b, 0x0d, 0xea, 0xf7, 0xf7, 0xfb,
	0x21, 0x56, 0xbd, 0xe1, 0x4f, 0x03, 0xba, 0x6f, 0xd4, 0xb1, 0x53, 0x81, 0x05, 0x45, 0x67, 0xd0,
	0x50, 0xc3, 0x3d, 0xe3, 0xc4, 0x38, 0xed, 0x38, 0x47, 0xd6, 0x9e, 0x0d, 0xeb, 0xa3, 0x14, 0xb8,
	0xf5, 0x9b, 0xbb, 0x41, 0xc5, 0xd7, 0x72, 0xf4, 0x0a, 0x9a, 0x84, 0x04, 0x49, 0xcc, 0x45, 0xaf,
	0x7a, 0x52, 0x3b, 0xed, 0x38, 0x8f, 0x4b, 0x26, 0x3d, 0xcf, 0x3d, 0xcc, 0xa7, 0x36, 0x77, 0x83,
	0x86, 0xe7, 0x4d, 0x62, 0x2e, 0xfc, 0x06, 0x21, 0xf9, 0x2f, 0x72, 0xa1, 0xc5, 0x0b, 0x40, 0x4d,
	0x02, 0x9e, 0x94, 0x00, 0xa6, 0x9e, 0xe7, 0x3e, 0xd0, 0x84, 0xe6, 0x54, 0x23, 0x9a, 0x5c, 0x31,
	0x86, 0xdf, 0x6a, 0x50, 0xf5, 0x3c, 0x34, 0x80, 0x4e, 0xc6, 0x92, 0x64, 0x95, 0x06, 0x4b, 0xbc,
	0xa0, 0x72, 0x91, 0xb6, 0x0f, 0xaa, 0xf4, 0x01, 0x2f, 0x28, 0x1a, 0x43, 0x8b, 0x90, 0x80, 0xe7,
	0x0b, 0xf7, 0xaa, 0x72, 0xcd, 0x7e, 0xa9, 0x59, 0x19, 0xc9, 0xfd, 0x79, 0xba, 0xe0, 0x37, 0x09,
	0x51, 0x61, 0xbd, 0x80, 0x66, 0x14, 0x73, 0xc1, 0xb2, 0x6b, 0x6d, 0xb9, 0x57, 0x8a, 0xf1, 0xe9,
	0x85, 0x0e, 0xab, 0x90, 0xa3, 0x77, 0x70, 0x78, 0xb9, 0xa2, 0x2b, 0x1a, 0x88, 0xb5, 0xf6, 0x51,
	0x97, 0x3e, 0x06, 0x25, 0x80, 0x4f, 0xb9, 0xf0, 0xf3, 0x5a, 0x99, 0x51, 0x9c, 0xee, 0xe5, 0x4e,
	0x0d, 0x79, 0x00, 0x0a, 0x26, 0xc3, 0x3b, 0x90, 0x4e, 0xcc, 0x12, 0xd0, 0x64, 0x24, 0xd8, 0xc4,
	0x91, 0x38, 0xcd, 0x69, 0xcb, 0x39, 0x99, 0xff, 0x6b, 0xe8, 0x26, 0x4e, 0x30, 0xc3, 0x82, 0x44,
	0xc1, 0x02, 0xa7, 0xbd, 0x86, 0xc4, 0x3c, 0x2d, 0xc3, 0x38, 0x6e, 0xae, 0x7a, 0x8f, 0x53, 0x4d,
	0x81, 0xe4, 0x6f, 0x65, 0xf8, 0xc3, 0x80, 0xda, 0xf4, 0x7f, 0xee, 0xe0, 0x0c, 0x0e, 0x76, 0x2f,
	0xe0, 0xb8, 0xfc, 0xb2, 0x77, 0x97, 0x56, 0x7a, 0xf4, 0xf2, 0xdf, 0xd0, 0x8f, 0xca, 0x47, 0xf7,
	0x53, 0x1f, 0x8e, 0x01, 0xee, 0xcd, 0xa3, 0x63, 0x68, 0x27, 0x4e, 0x10, 0xd1, 0x78, 0x1e, 0x09,
	0x69, 0xb0, 0xee, 0xb7, 0x12, 0xe7, 0x5c, 0x7e, 0xe7, 0x4d, 0x95, 0x45, 0x1c, 0xae, 0xa5, 0xc5,
	0xba, 0xdf, 0x92, 0x85, 0xb7, 0xe1, 0xda, 0x3d, 0xbf, 0xd9, 0x98, 0xc6, 0xed, 0xc6, 0x34, 0x7e,
	0x6f, 0x4c, 0xe3, 0xfb, 0xd6, 0xac, 0xdc, 0x6e, 0xcd, 0xca, 0xaf, 0xad, 0x59, 0xf9, 0x62, 0xcd,
	0x63, 0x11, 0xad, 0x66, 0x16, 0x61, 0x0b, 0x7b, 0x5c, 0x3c, 0xbb, 0xc2, 0xde, 0x33, 0x1e, 0x7e,
	0xb5, 0xd7, 0xfa, 0x15, 0x8a, 0xeb, 0x94, 0xf2, 0x59, 0x43, 0x3e, 0xc3, 0xe7, 0x7f, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x7b, 0x85, 0x7f, 0xc9, 0x0e, 0x04, 0x00, 0x00,
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
	if len(m.SCCList) > 0 {
		for iNdEx := len(m.SCCList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SCCList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.CCList) > 0 {
		for iNdEx := len(m.CCList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.CCList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	return len(dAtA) - i, nil
}

func (m *CC) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CC) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CC) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.L2BatchMap) > 0 {
		for iNdEx := len(m.L2BatchMap) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.L2BatchMap[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.QueueList) > 0 {
		for iNdEx := len(m.QueueList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.QueueList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	{
		size, err := m.QueueTxState.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.History) > 0 {
		for iNdEx := len(m.History) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.History[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size, err := m.CCState.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.RollupName) > 0 {
		i -= len(m.RollupName)
		copy(dAtA[i:], m.RollupName)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.RollupName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SCC) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SCC) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SCC) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.History) > 0 {
		for iNdEx := len(m.History) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.History[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size, err := m.State.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.RollupName) > 0 {
		i -= len(m.RollupName)
		copy(dAtA[i:], m.RollupName)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.RollupName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *L2BatchMap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *L2BatchMap) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *L2BatchMap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BatchIdx != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.BatchIdx))
		i--
		dAtA[i] = 0x10
	}
	if m.L2Height != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.L2Height))
		i--
		dAtA[i] = 0x8
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
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.CCList) > 0 {
		for _, e := range m.CCList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.SCCList) > 0 {
		for _, e := range m.SCCList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *CC) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RollupName)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = m.CCState.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.History) > 0 {
		for _, e := range m.History {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.QueueTxState.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.QueueList) > 0 {
		for _, e := range m.QueueList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.L2BatchMap) > 0 {
		for _, e := range m.L2BatchMap {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *SCC) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RollupName)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = m.State.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.History) > 0 {
		for _, e := range m.History {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *L2BatchMap) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.L2Height != 0 {
		n += 1 + sovGenesis(uint64(m.L2Height))
	}
	if m.BatchIdx != 0 {
		n += 1 + sovGenesis(uint64(m.BatchIdx))
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
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CCList", wireType)
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
			m.CCList = append(m.CCList, CC{})
			if err := m.CCList[len(m.CCList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SCCList", wireType)
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
			m.SCCList = append(m.SCCList, SCC{})
			if err := m.SCCList[len(m.SCCList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *CC) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: CC: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CC: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RollupName", wireType)
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
			m.RollupName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CCState", wireType)
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
			if err := m.CCState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field History", wireType)
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
			m.History = append(m.History, CCRef{})
			if err := m.History[len(m.History)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueueTxState", wireType)
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
			if err := m.QueueTxState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueueList", wireType)
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
			m.QueueList = append(m.QueueList, L1ToL2Queue{})
			if err := m.QueueList[len(m.QueueList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field L2BatchMap", wireType)
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
			m.L2BatchMap = append(m.L2BatchMap, L2BatchMap{})
			if err := m.L2BatchMap[len(m.L2BatchMap)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *SCC) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SCC: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SCC: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RollupName", wireType)
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
			m.RollupName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
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
			if err := m.State.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field History", wireType)
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
			m.History = append(m.History, SCCRef{})
			if err := m.History[len(m.History)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *L2BatchMap) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: L2BatchMap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: L2BatchMap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field L2Height", wireType)
			}
			m.L2Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.L2Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchIdx", wireType)
			}
			m.BatchIdx = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BatchIdx |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
