// Package handler 用于接口逻辑处理
package handler

//Run 执行函数
func Run(handler HandlerIntf) error {
	err := handler.HandleInput()
	if err != nil {
		return err
	}
	err = handler.HandleProcess()
	return err
}

// HandlerIntf handler接口
type HandlerIntf interface {
	HandleInput() error
	HandleProcess() error
}

//RpcIntf Rpc接口
type RpcIntf interface {
	Request(reqBody interface{}) (interface{}, error)
}

//RunHandler new version of handler, have cache
func RunHandler(handler Handler) error {
	err := handler.HandleInput()
	if err != nil {
		return err
	}
	if ok := handler.UseCache(); ok {
		return nil
	}
	err = handler.HandleProcess()
	if err != nil {
		return err
	}
	handler.SetCache()
	return err
}

//HandlerBase handler共有属性
type HandlerBase struct {
	CacheKey string
}

//UseCache 使用缓存
func (p *HandlerBase) UseCache() bool {
	return false
}

//SetCache 设置缓存
func (p *HandlerBase) SetCache() {

}

//Handler hanlder接口
type Handler interface {
	HandleInput() error
	HandleProcess() error
	UseCache() bool
	SetCache()
}
