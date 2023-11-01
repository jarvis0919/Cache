package peer

import pb "cache/api"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool, peers string)
	Set(peers ...string)
	RemovePeer(peer string)
}

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
