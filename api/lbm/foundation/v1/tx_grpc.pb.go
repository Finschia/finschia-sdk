// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: lbm/foundation/v1/tx.proto

package foundationv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Msg_FundTreasury_FullMethodName         = "/lbm.foundation.v1.Msg/FundTreasury"
	Msg_WithdrawFromTreasury_FullMethodName = "/lbm.foundation.v1.Msg/WithdrawFromTreasury"
	Msg_UpdateMembers_FullMethodName        = "/lbm.foundation.v1.Msg/UpdateMembers"
	Msg_UpdateDecisionPolicy_FullMethodName = "/lbm.foundation.v1.Msg/UpdateDecisionPolicy"
	Msg_SubmitProposal_FullMethodName       = "/lbm.foundation.v1.Msg/SubmitProposal"
	Msg_WithdrawProposal_FullMethodName     = "/lbm.foundation.v1.Msg/WithdrawProposal"
	Msg_Vote_FullMethodName                 = "/lbm.foundation.v1.Msg/Vote"
	Msg_Exec_FullMethodName                 = "/lbm.foundation.v1.Msg/Exec"
	Msg_LeaveFoundation_FullMethodName      = "/lbm.foundation.v1.Msg/LeaveFoundation"
	Msg_UpdateCensorship_FullMethodName     = "/lbm.foundation.v1.Msg/UpdateCensorship"
	Msg_Grant_FullMethodName                = "/lbm.foundation.v1.Msg/Grant"
	Msg_Revoke_FullMethodName               = "/lbm.foundation.v1.Msg/Revoke"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// FundTreasury defines a method to fund the treasury.
	FundTreasury(ctx context.Context, in *MsgFundTreasury, opts ...grpc.CallOption) (*MsgFundTreasuryResponse, error)
	// WithdrawFromTreasury defines a method to withdraw coins from the treasury.
	WithdrawFromTreasury(ctx context.Context, in *MsgWithdrawFromTreasury, opts ...grpc.CallOption) (*MsgWithdrawFromTreasuryResponse, error)
	// UpdateMembers updates the foundation members.
	UpdateMembers(ctx context.Context, in *MsgUpdateMembers, opts ...grpc.CallOption) (*MsgUpdateMembersResponse, error)
	// UpdateDecisionPolicy allows a group policy's decision policy to be updated.
	UpdateDecisionPolicy(ctx context.Context, in *MsgUpdateDecisionPolicy, opts ...grpc.CallOption) (*MsgUpdateDecisionPolicyResponse, error)
	// SubmitProposal submits a new proposal.
	SubmitProposal(ctx context.Context, in *MsgSubmitProposal, opts ...grpc.CallOption) (*MsgSubmitProposalResponse, error)
	// WithdrawProposal aborts a proposal.
	WithdrawProposal(ctx context.Context, in *MsgWithdrawProposal, opts ...grpc.CallOption) (*MsgWithdrawProposalResponse, error)
	// Vote allows a voter to vote on a proposal.
	Vote(ctx context.Context, in *MsgVote, opts ...grpc.CallOption) (*MsgVoteResponse, error)
	// Exec executes a proposal.
	Exec(ctx context.Context, in *MsgExec, opts ...grpc.CallOption) (*MsgExecResponse, error)
	// LeaveFoundation allows a member to leave the foundation.
	LeaveFoundation(ctx context.Context, in *MsgLeaveFoundation, opts ...grpc.CallOption) (*MsgLeaveFoundationResponse, error)
	// UpdateCensorship updates censorship information.
	UpdateCensorship(ctx context.Context, in *MsgUpdateCensorship, opts ...grpc.CallOption) (*MsgUpdateCensorshipResponse, error)
	// Grant grants the provided authorization to the grantee with authority of
	// the foundation. If there is already a grant for the given
	// (grantee, Authorization) tuple, then the grant will be overwritten.
	Grant(ctx context.Context, in *MsgGrant, opts ...grpc.CallOption) (*MsgGrantResponse, error)
	// Revoke revokes any authorization corresponding to the provided method name
	// that has been granted to the grantee.
	Revoke(ctx context.Context, in *MsgRevoke, opts ...grpc.CallOption) (*MsgRevokeResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) FundTreasury(ctx context.Context, in *MsgFundTreasury, opts ...grpc.CallOption) (*MsgFundTreasuryResponse, error) {
	out := new(MsgFundTreasuryResponse)
	err := c.cc.Invoke(ctx, Msg_FundTreasury_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) WithdrawFromTreasury(ctx context.Context, in *MsgWithdrawFromTreasury, opts ...grpc.CallOption) (*MsgWithdrawFromTreasuryResponse, error) {
	out := new(MsgWithdrawFromTreasuryResponse)
	err := c.cc.Invoke(ctx, Msg_WithdrawFromTreasury_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateMembers(ctx context.Context, in *MsgUpdateMembers, opts ...grpc.CallOption) (*MsgUpdateMembersResponse, error) {
	out := new(MsgUpdateMembersResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateMembers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateDecisionPolicy(ctx context.Context, in *MsgUpdateDecisionPolicy, opts ...grpc.CallOption) (*MsgUpdateDecisionPolicyResponse, error) {
	out := new(MsgUpdateDecisionPolicyResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateDecisionPolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SubmitProposal(ctx context.Context, in *MsgSubmitProposal, opts ...grpc.CallOption) (*MsgSubmitProposalResponse, error) {
	out := new(MsgSubmitProposalResponse)
	err := c.cc.Invoke(ctx, Msg_SubmitProposal_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) WithdrawProposal(ctx context.Context, in *MsgWithdrawProposal, opts ...grpc.CallOption) (*MsgWithdrawProposalResponse, error) {
	out := new(MsgWithdrawProposalResponse)
	err := c.cc.Invoke(ctx, Msg_WithdrawProposal_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Vote(ctx context.Context, in *MsgVote, opts ...grpc.CallOption) (*MsgVoteResponse, error) {
	out := new(MsgVoteResponse)
	err := c.cc.Invoke(ctx, Msg_Vote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Exec(ctx context.Context, in *MsgExec, opts ...grpc.CallOption) (*MsgExecResponse, error) {
	out := new(MsgExecResponse)
	err := c.cc.Invoke(ctx, Msg_Exec_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) LeaveFoundation(ctx context.Context, in *MsgLeaveFoundation, opts ...grpc.CallOption) (*MsgLeaveFoundationResponse, error) {
	out := new(MsgLeaveFoundationResponse)
	err := c.cc.Invoke(ctx, Msg_LeaveFoundation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateCensorship(ctx context.Context, in *MsgUpdateCensorship, opts ...grpc.CallOption) (*MsgUpdateCensorshipResponse, error) {
	out := new(MsgUpdateCensorshipResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateCensorship_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Grant(ctx context.Context, in *MsgGrant, opts ...grpc.CallOption) (*MsgGrantResponse, error) {
	out := new(MsgGrantResponse)
	err := c.cc.Invoke(ctx, Msg_Grant_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Revoke(ctx context.Context, in *MsgRevoke, opts ...grpc.CallOption) (*MsgRevokeResponse, error) {
	out := new(MsgRevokeResponse)
	err := c.cc.Invoke(ctx, Msg_Revoke_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// FundTreasury defines a method to fund the treasury.
	FundTreasury(context.Context, *MsgFundTreasury) (*MsgFundTreasuryResponse, error)
	// WithdrawFromTreasury defines a method to withdraw coins from the treasury.
	WithdrawFromTreasury(context.Context, *MsgWithdrawFromTreasury) (*MsgWithdrawFromTreasuryResponse, error)
	// UpdateMembers updates the foundation members.
	UpdateMembers(context.Context, *MsgUpdateMembers) (*MsgUpdateMembersResponse, error)
	// UpdateDecisionPolicy allows a group policy's decision policy to be updated.
	UpdateDecisionPolicy(context.Context, *MsgUpdateDecisionPolicy) (*MsgUpdateDecisionPolicyResponse, error)
	// SubmitProposal submits a new proposal.
	SubmitProposal(context.Context, *MsgSubmitProposal) (*MsgSubmitProposalResponse, error)
	// WithdrawProposal aborts a proposal.
	WithdrawProposal(context.Context, *MsgWithdrawProposal) (*MsgWithdrawProposalResponse, error)
	// Vote allows a voter to vote on a proposal.
	Vote(context.Context, *MsgVote) (*MsgVoteResponse, error)
	// Exec executes a proposal.
	Exec(context.Context, *MsgExec) (*MsgExecResponse, error)
	// LeaveFoundation allows a member to leave the foundation.
	LeaveFoundation(context.Context, *MsgLeaveFoundation) (*MsgLeaveFoundationResponse, error)
	// UpdateCensorship updates censorship information.
	UpdateCensorship(context.Context, *MsgUpdateCensorship) (*MsgUpdateCensorshipResponse, error)
	// Grant grants the provided authorization to the grantee with authority of
	// the foundation. If there is already a grant for the given
	// (grantee, Authorization) tuple, then the grant will be overwritten.
	Grant(context.Context, *MsgGrant) (*MsgGrantResponse, error)
	// Revoke revokes any authorization corresponding to the provided method name
	// that has been granted to the grantee.
	Revoke(context.Context, *MsgRevoke) (*MsgRevokeResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) FundTreasury(context.Context, *MsgFundTreasury) (*MsgFundTreasuryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FundTreasury not implemented")
}
func (UnimplementedMsgServer) WithdrawFromTreasury(context.Context, *MsgWithdrawFromTreasury) (*MsgWithdrawFromTreasuryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawFromTreasury not implemented")
}
func (UnimplementedMsgServer) UpdateMembers(context.Context, *MsgUpdateMembers) (*MsgUpdateMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMembers not implemented")
}
func (UnimplementedMsgServer) UpdateDecisionPolicy(context.Context, *MsgUpdateDecisionPolicy) (*MsgUpdateDecisionPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDecisionPolicy not implemented")
}
func (UnimplementedMsgServer) SubmitProposal(context.Context, *MsgSubmitProposal) (*MsgSubmitProposalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitProposal not implemented")
}
func (UnimplementedMsgServer) WithdrawProposal(context.Context, *MsgWithdrawProposal) (*MsgWithdrawProposalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawProposal not implemented")
}
func (UnimplementedMsgServer) Vote(context.Context, *MsgVote) (*MsgVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
func (UnimplementedMsgServer) Exec(context.Context, *MsgExec) (*MsgExecResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exec not implemented")
}
func (UnimplementedMsgServer) LeaveFoundation(context.Context, *MsgLeaveFoundation) (*MsgLeaveFoundationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveFoundation not implemented")
}
func (UnimplementedMsgServer) UpdateCensorship(context.Context, *MsgUpdateCensorship) (*MsgUpdateCensorshipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCensorship not implemented")
}
func (UnimplementedMsgServer) Grant(context.Context, *MsgGrant) (*MsgGrantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Grant not implemented")
}
func (UnimplementedMsgServer) Revoke(context.Context, *MsgRevoke) (*MsgRevokeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revoke not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_FundTreasury_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFundTreasury)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FundTreasury(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_FundTreasury_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FundTreasury(ctx, req.(*MsgFundTreasury))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawFromTreasury_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawFromTreasury)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawFromTreasury(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_WithdrawFromTreasury_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawFromTreasury(ctx, req.(*MsgWithdrawFromTreasury))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateMembers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateMembers(ctx, req.(*MsgUpdateMembers))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateDecisionPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateDecisionPolicy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateDecisionPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateDecisionPolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateDecisionPolicy(ctx, req.(*MsgUpdateDecisionPolicy))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SubmitProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitProposal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitProposal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SubmitProposal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitProposal(ctx, req.(*MsgSubmitProposal))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawProposal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawProposal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_WithdrawProposal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawProposal(ctx, req.(*MsgWithdrawProposal))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Vote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Vote(ctx, req.(*MsgVote))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Exec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Exec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Exec_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Exec(ctx, req.(*MsgExec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_LeaveFoundation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgLeaveFoundation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).LeaveFoundation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_LeaveFoundation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).LeaveFoundation(ctx, req.(*MsgLeaveFoundation))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateCensorship_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateCensorship)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateCensorship(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateCensorship_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateCensorship(ctx, req.(*MsgUpdateCensorship))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Grant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgGrant)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Grant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Grant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Grant(ctx, req.(*MsgGrant))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Revoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRevoke)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Revoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Revoke_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Revoke(ctx, req.(*MsgRevoke))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lbm.foundation.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FundTreasury",
			Handler:    _Msg_FundTreasury_Handler,
		},
		{
			MethodName: "WithdrawFromTreasury",
			Handler:    _Msg_WithdrawFromTreasury_Handler,
		},
		{
			MethodName: "UpdateMembers",
			Handler:    _Msg_UpdateMembers_Handler,
		},
		{
			MethodName: "UpdateDecisionPolicy",
			Handler:    _Msg_UpdateDecisionPolicy_Handler,
		},
		{
			MethodName: "SubmitProposal",
			Handler:    _Msg_SubmitProposal_Handler,
		},
		{
			MethodName: "WithdrawProposal",
			Handler:    _Msg_WithdrawProposal_Handler,
		},
		{
			MethodName: "Vote",
			Handler:    _Msg_Vote_Handler,
		},
		{
			MethodName: "Exec",
			Handler:    _Msg_Exec_Handler,
		},
		{
			MethodName: "LeaveFoundation",
			Handler:    _Msg_LeaveFoundation_Handler,
		},
		{
			MethodName: "UpdateCensorship",
			Handler:    _Msg_UpdateCensorship_Handler,
		},
		{
			MethodName: "Grant",
			Handler:    _Msg_Grant_Handler,
		},
		{
			MethodName: "Revoke",
			Handler:    _Msg_Revoke_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbm/foundation/v1/tx.proto",
}
