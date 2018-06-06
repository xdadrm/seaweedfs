package filesys

import (
	"bazil.org/fuse/fs"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/karlseguin/ccache"
	"google.golang.org/grpc"
)

type WFS struct {
	filerGrpcAddress          string
	listDirectoryEntriesCache *ccache.Cache
	collection                string
	replication               string
	chunkSizeLimit            int64
}

func NewSeaweedFileSystem(filerGrpcAddress string, collection string, replication string, chunkSizeLimitMB int) *WFS {
	return &WFS{
		filerGrpcAddress:          filerGrpcAddress,
		listDirectoryEntriesCache: ccache.New(ccache.Configure().MaxSize(6000).ItemsToPrune(100)),
		collection:                collection,
		replication:               replication,
		chunkSizeLimit:            int64(chunkSizeLimitMB) * 1024 * 1024,
	}
}

func (wfs *WFS) Root() (fs.Node, error) {
	return &Dir{Path: "/", wfs: wfs}, nil
}

func (wfs *WFS) withFilerClient(fn func(filer_pb.SeaweedFilerClient) error) error {

	grpcConnection, err := grpc.Dial(wfs.filerGrpcAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("fail to dial %s: %v", wfs.filerGrpcAddress, err)
	}
	defer grpcConnection.Close()

	client := filer_pb.NewSeaweedFilerClient(grpcConnection)

	return fn(client)
}
