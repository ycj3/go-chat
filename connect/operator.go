package connect

import "go-chat/proto"

type Operator interface {
	Connect(conn *proto.ConnectRequest) (int, error)
	DisConnect(disConn *proto.DisConnectRequest) (err error)
}

type DefaultOperator struct {
}

// rpc call logic layer
func (o *DefaultOperator) Connect(conn *proto.ConnectRequest) (uid int, err error) {
	rpcConnect := new(RpcConnect)
	uid, err = rpcConnect.Connect(conn)
	return
}

// rpc call logic layer
func (o *DefaultOperator) DisConnect(disConn *proto.DisConnectRequest) (err error) {
	rpcConnect := new(RpcConnect)
	err = rpcConnect.DisConnect(disConn)
	return
}
