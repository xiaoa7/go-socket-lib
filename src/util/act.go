//指令集
package util

const (
	//系统指令，请不要修改
	ACT_REQUEST_CLIENTID   = iota //请求用户ID
	ACT_RESPONSE_CLIENTID         //回应用户ID
	ACT_REQUEST_HEARTBEAT         //请求心跳
	ACT_RESPONSE_HEARTBEAT        //回应心跳
	ACT_MESSAGE                   //消息
)
