package gocache

//PeerPicker 根据传入的key选择对应的节点
type PeerPicker interface {
	//PickPeer 根据传入的key选择对应的节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//PeerGetter 节点实体
type PeerGetter interface {
	// Get 从对应的group节点中获取缓存值
	Get(group string, key string) ([]byte, error)
}
