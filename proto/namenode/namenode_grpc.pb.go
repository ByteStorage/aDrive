// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/namenode/namenode.proto

package namenode

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

// NameNodeServiceClient is the client API for NameNodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NameNodeServiceClient interface {
	HeartBeat(ctx context.Context, in *HeartBeatReq, opts ...grpc.CallOption) (*HeartBeatResp, error)
	RegisterDataNode(ctx context.Context, in *RegisterDataNodeReq, opts ...grpc.CallOption) (*RegisterDataNodeResp, error)
	JoinCluster(ctx context.Context, in *JoinClusterReq, opts ...grpc.CallOption) (*JoinClusterResp, error)
	Delete(ctx context.Context, in *DeleteDataReq, opts ...grpc.CallOption) (*DeleteDataResp, error)
	FindLeader(ctx context.Context, in *FindLeaderReq, opts ...grpc.CallOption) (*FindLeaderResp, error)
	List(ctx context.Context, in *ListReq, opts ...grpc.CallOption) (*ListResp, error)
	UpdateDataNodeMessage(ctx context.Context, in *UpdateDataNodeMessageReq, opts ...grpc.CallOption) (*UpdateDataNodeMessageResp, error)
	Put(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*PutResp, error)
	IsDir(ctx context.Context, in *IsDirReq, opts ...grpc.CallOption) (*IsDirResp, error)
	Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error)
	Mkdir(ctx context.Context, in *MkdirReq, opts ...grpc.CallOption) (*MkdirResp, error)
}

type nameNodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNameNodeServiceClient(cc grpc.ClientConnInterface) NameNodeServiceClient {
	return &nameNodeServiceClient{cc}
}

func (c *nameNodeServiceClient) HeartBeat(ctx context.Context, in *HeartBeatReq, opts ...grpc.CallOption) (*HeartBeatResp, error) {
	out := new(HeartBeatResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) RegisterDataNode(ctx context.Context, in *RegisterDataNodeReq, opts ...grpc.CallOption) (*RegisterDataNodeResp, error) {
	out := new(RegisterDataNodeResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/RegisterDataNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) JoinCluster(ctx context.Context, in *JoinClusterReq, opts ...grpc.CallOption) (*JoinClusterResp, error) {
	out := new(JoinClusterResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/JoinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) Delete(ctx context.Context, in *DeleteDataReq, opts ...grpc.CallOption) (*DeleteDataResp, error) {
	out := new(DeleteDataResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) FindLeader(ctx context.Context, in *FindLeaderReq, opts ...grpc.CallOption) (*FindLeaderResp, error) {
	out := new(FindLeaderResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/FindLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) List(ctx context.Context, in *ListReq, opts ...grpc.CallOption) (*ListResp, error) {
	out := new(ListResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) UpdateDataNodeMessage(ctx context.Context, in *UpdateDataNodeMessageReq, opts ...grpc.CallOption) (*UpdateDataNodeMessageResp, error) {
	out := new(UpdateDataNodeMessageResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/UpdateDataNodeMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) Put(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*PutResp, error) {
	out := new(PutResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) IsDir(ctx context.Context, in *IsDirReq, opts ...grpc.CallOption) (*IsDirResp, error) {
	out := new(IsDirResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/IsDir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error) {
	out := new(GetResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameNodeServiceClient) Mkdir(ctx context.Context, in *MkdirReq, opts ...grpc.CallOption) (*MkdirResp, error) {
	out := new(MkdirResp)
	err := c.cc.Invoke(ctx, "/namenode_.NameNodeService/Mkdir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NameNodeServiceServer is the server API for NameNodeService service.
// All implementations must embed UnimplementedNameNodeServiceServer
// for forward compatibility
type NameNodeServiceServer interface {
	HeartBeat(context.Context, *HeartBeatReq) (*HeartBeatResp, error)
	RegisterDataNode(context.Context, *RegisterDataNodeReq) (*RegisterDataNodeResp, error)
	JoinCluster(context.Context, *JoinClusterReq) (*JoinClusterResp, error)
	Delete(context.Context, *DeleteDataReq) (*DeleteDataResp, error)
	FindLeader(context.Context, *FindLeaderReq) (*FindLeaderResp, error)
	List(context.Context, *ListReq) (*ListResp, error)
	UpdateDataNodeMessage(context.Context, *UpdateDataNodeMessageReq) (*UpdateDataNodeMessageResp, error)
	Put(context.Context, *PutReq) (*PutResp, error)
	IsDir(context.Context, *IsDirReq) (*IsDirResp, error)
	Get(context.Context, *GetReq) (*GetResp, error)
	Mkdir(context.Context, *MkdirReq) (*MkdirResp, error)
	mustEmbedUnimplementedNameNodeServiceServer()
}

// UnimplementedNameNodeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNameNodeServiceServer struct {
}

func (UnimplementedNameNodeServiceServer) HeartBeat(context.Context, *HeartBeatReq) (*HeartBeatResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}
func (UnimplementedNameNodeServiceServer) RegisterDataNode(context.Context, *RegisterDataNodeReq) (*RegisterDataNodeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDataNode not implemented")
}
func (UnimplementedNameNodeServiceServer) JoinCluster(context.Context, *JoinClusterReq) (*JoinClusterResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}
func (UnimplementedNameNodeServiceServer) Delete(context.Context, *DeleteDataReq) (*DeleteDataResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedNameNodeServiceServer) FindLeader(context.Context, *FindLeaderReq) (*FindLeaderResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindLeader not implemented")
}
func (UnimplementedNameNodeServiceServer) List(context.Context, *ListReq) (*ListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedNameNodeServiceServer) UpdateDataNodeMessage(context.Context, *UpdateDataNodeMessageReq) (*UpdateDataNodeMessageResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDataNodeMessage not implemented")
}
func (UnimplementedNameNodeServiceServer) Put(context.Context, *PutReq) (*PutResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedNameNodeServiceServer) IsDir(context.Context, *IsDirReq) (*IsDirResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsDir not implemented")
}
func (UnimplementedNameNodeServiceServer) Get(context.Context, *GetReq) (*GetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedNameNodeServiceServer) Mkdir(context.Context, *MkdirReq) (*MkdirResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mkdir not implemented")
}
func (UnimplementedNameNodeServiceServer) mustEmbedUnimplementedNameNodeServiceServer() {}

// UnsafeNameNodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NameNodeServiceServer will
// result in compilation errors.
type UnsafeNameNodeServiceServer interface {
	mustEmbedUnimplementedNameNodeServiceServer()
}

func RegisterNameNodeServiceServer(s grpc.ServiceRegistrar, srv NameNodeServiceServer) {
	s.RegisterService(&NameNodeService_ServiceDesc, srv)
}

func _NameNodeService_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).HeartBeat(ctx, req.(*HeartBeatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_RegisterDataNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterDataNodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).RegisterDataNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/RegisterDataNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).RegisterDataNode(ctx, req.(*RegisterDataNodeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinClusterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).JoinCluster(ctx, req.(*JoinClusterReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).Delete(ctx, req.(*DeleteDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_FindLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindLeaderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).FindLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/FindLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).FindLeader(ctx, req.(*FindLeaderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).List(ctx, req.(*ListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_UpdateDataNodeMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDataNodeMessageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).UpdateDataNodeMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/UpdateDataNodeMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).UpdateDataNodeMessage(ctx, req.(*UpdateDataNodeMessageReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).Put(ctx, req.(*PutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_IsDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsDirReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).IsDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/IsDir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).IsDir(ctx, req.(*IsDirReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).Get(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NameNodeService_Mkdir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MkdirReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NameNodeServiceServer).Mkdir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/namenode_.NameNodeService/Mkdir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NameNodeServiceServer).Mkdir(ctx, req.(*MkdirReq))
	}
	return interceptor(ctx, in, info, handler)
}

// NameNodeService_ServiceDesc is the grpc.ServiceDesc for NameNodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NameNodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "namenode_.NameNodeService",
	HandlerType: (*NameNodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HeartBeat",
			Handler:    _NameNodeService_HeartBeat_Handler,
		},
		{
			MethodName: "RegisterDataNode",
			Handler:    _NameNodeService_RegisterDataNode_Handler,
		},
		{
			MethodName: "JoinCluster",
			Handler:    _NameNodeService_JoinCluster_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _NameNodeService_Delete_Handler,
		},
		{
			MethodName: "FindLeader",
			Handler:    _NameNodeService_FindLeader_Handler,
		},
		{
			MethodName: "List",
			Handler:    _NameNodeService_List_Handler,
		},
		{
			MethodName: "UpdateDataNodeMessage",
			Handler:    _NameNodeService_UpdateDataNodeMessage_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _NameNodeService_Put_Handler,
		},
		{
			MethodName: "IsDir",
			Handler:    _NameNodeService_IsDir_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _NameNodeService_Get_Handler,
		},
		{
			MethodName: "Mkdir",
			Handler:    _NameNodeService_Mkdir_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/namenode/namenode.proto",
}
