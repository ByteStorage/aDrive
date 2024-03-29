// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/datanode/datanode.proto

package datanode

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

// DataNodeClient is the client API for DataNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataNodeClient interface {
	Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingResp, error)
	HeartBeat(ctx context.Context, in *HeartBeatReq, opts ...grpc.CallOption) (*HeartBeatResp, error)
	Put(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*PutResp, error)
	Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error)
	Delete(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*DeleteResp, error)
	Stat(ctx context.Context, in *StatReq, opts ...grpc.CallOption) (*StatResp, error)
	List(ctx context.Context, in *ListReq, opts ...grpc.CallOption) (*ListResp, error)
	Rename(ctx context.Context, in *RenameReq, opts ...grpc.CallOption) (*RenameResp, error)
	Mkdir(ctx context.Context, in *MkdirReq, opts ...grpc.CallOption) (*MkdirResp, error)
}

type dataNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewDataNodeClient(cc grpc.ClientConnInterface) DataNodeClient {
	return &dataNodeClient{cc}
}

func (c *dataNodeClient) Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingResp, error) {
	out := new(PingResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) HeartBeat(ctx context.Context, in *HeartBeatReq, opts ...grpc.CallOption) (*HeartBeatResp, error) {
	out := new(HeartBeatResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) Put(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*PutResp, error) {
	out := new(PutResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error) {
	out := new(GetResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) Delete(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*DeleteResp, error) {
	out := new(DeleteResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) Stat(ctx context.Context, in *StatReq, opts ...grpc.CallOption) (*StatResp, error) {
	out := new(StatResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/Stat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) List(ctx context.Context, in *ListReq, opts ...grpc.CallOption) (*ListResp, error) {
	out := new(ListResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) Rename(ctx context.Context, in *RenameReq, opts ...grpc.CallOption) (*RenameResp, error) {
	out := new(RenameResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/Rename", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataNodeClient) Mkdir(ctx context.Context, in *MkdirReq, opts ...grpc.CallOption) (*MkdirResp, error) {
	out := new(MkdirResp)
	err := c.cc.Invoke(ctx, "/datanode.DataNode/Mkdir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataNodeServer is the server API for DataNode service.
// All implementations must embed UnimplementedDataNodeServer
// for forward compatibility
type DataNodeServer interface {
	Ping(context.Context, *PingReq) (*PingResp, error)
	HeartBeat(context.Context, *HeartBeatReq) (*HeartBeatResp, error)
	Put(context.Context, *PutReq) (*PutResp, error)
	Get(context.Context, *GetReq) (*GetResp, error)
	Delete(context.Context, *DeleteReq) (*DeleteResp, error)
	Stat(context.Context, *StatReq) (*StatResp, error)
	List(context.Context, *ListReq) (*ListResp, error)
	Rename(context.Context, *RenameReq) (*RenameResp, error)
	Mkdir(context.Context, *MkdirReq) (*MkdirResp, error)
	mustEmbedUnimplementedDataNodeServer()
}

// UnimplementedDataNodeServer must be embedded to have forward compatible implementations.
type UnimplementedDataNodeServer struct {
}

func (UnimplementedDataNodeServer) Ping(context.Context, *PingReq) (*PingResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDataNodeServer) HeartBeat(context.Context, *HeartBeatReq) (*HeartBeatResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}
func (UnimplementedDataNodeServer) Put(context.Context, *PutReq) (*PutResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedDataNodeServer) Get(context.Context, *GetReq) (*GetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDataNodeServer) Delete(context.Context, *DeleteReq) (*DeleteResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDataNodeServer) Stat(context.Context, *StatReq) (*StatResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stat not implemented")
}
func (UnimplementedDataNodeServer) List(context.Context, *ListReq) (*ListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDataNodeServer) Rename(context.Context, *RenameReq) (*RenameResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rename not implemented")
}
func (UnimplementedDataNodeServer) Mkdir(context.Context, *MkdirReq) (*MkdirResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mkdir not implemented")
}
func (UnimplementedDataNodeServer) mustEmbedUnimplementedDataNodeServer() {}

// UnsafeDataNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataNodeServer will
// result in compilation errors.
type UnsafeDataNodeServer interface {
	mustEmbedUnimplementedDataNodeServer()
}

func RegisterDataNodeServer(s grpc.ServiceRegistrar, srv DataNodeServer) {
	s.RegisterService(&DataNode_ServiceDesc, srv)
}

func _DataNode_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).Ping(ctx, req.(*PingReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).HeartBeat(ctx, req.(*HeartBeatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).Put(ctx, req.(*PutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).Get(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).Delete(ctx, req.(*DeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_Stat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).Stat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/Stat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).Stat(ctx, req.(*StatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).List(ctx, req.(*ListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_Rename_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenameReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).Rename(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/Rename",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).Rename(ctx, req.(*RenameReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataNode_Mkdir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MkdirReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataNodeServer).Mkdir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datanode.DataNode/Mkdir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataNodeServer).Mkdir(ctx, req.(*MkdirReq))
	}
	return interceptor(ctx, in, info, handler)
}

// DataNode_ServiceDesc is the grpc.ServiceDesc for DataNode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataNode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datanode.DataNode",
	HandlerType: (*DataNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _DataNode_Ping_Handler,
		},
		{
			MethodName: "HeartBeat",
			Handler:    _DataNode_HeartBeat_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _DataNode_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DataNode_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DataNode_Delete_Handler,
		},
		{
			MethodName: "Stat",
			Handler:    _DataNode_Stat_Handler,
		},
		{
			MethodName: "List",
			Handler:    _DataNode_List_Handler,
		},
		{
			MethodName: "Rename",
			Handler:    _DataNode_Rename_Handler,
		},
		{
			MethodName: "Mkdir",
			Handler:    _DataNode_Mkdir_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/datanode/datanode.proto",
}
