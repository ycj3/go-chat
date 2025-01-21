package connect

type Operator interface {
	Connect(conn *ConnectRequest) (int, error)
	DisConnect(disConn *DisConnectRequest) (err error)
}

type DefaultOperator struct {
}

// rpc call logic layer
func (o *DefaultOperator) Connect(conn *ConnectRequest) (uid int, err error) {
	rpcConnect := new(RpcConnect)
	uid, err = rpcConnect.Connect(conn)
	return
}

// rpc call logic layer
func (o *DefaultOperator) DisConnect(disConn *DisConnectRequest) (err error) {
	rpcConnect := new(RpcConnect)
	err = rpcConnect.DisConnect(disConn)
	return
}
