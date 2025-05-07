package grpc2

import (
	"fmt"
	rpc "vk-test/grpcServer"
	"vk-test/subpub"
    "vk-test/logger"

	context "context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)
type Server struct {
    Subpub subpub.SubPub
    l logger.Logger
    rpc.UnimplementedPubSubServer
}

func senderAsMsgHandler(PubSubSS rpc.PubSub_SubscribeServer) func(interface{}) {
    return func(msg interface{}) {
        event := new(rpc.Event)
        event.Data = msg.(string)
        PubSubSS.Send(event)
    }
}

func (s *Server) Subscribe(sr *rpc.SubscribeRequest, PubSubSS rpc.PubSub_SubscribeServer) error {
    sub, err := s.Subpub.Subscribe(sr.Key, senderAsMsgHandler(PubSubSS))

    if err != nil {
        return err
    }

    for subpub.IsSubValid(sub) {}
    return nil
}

func (s *Server) Publish(ctx context.Context, pr *rpc.PublishRequest) (*emptypb.Empty, error) {
    var err error
    select {
    case <-ctx.Done():
        return &emptypb.Empty{}, status.Error(codes.Canceled, "Request cancelled by client")
    default:
        err = s.Subpub.Publish(pr.GetKey(), pr.GetData())
    }
    
    if err != nil {    
        return &emptypb.Empty{}, status.Error(status.Code(err), fmt.Sprintf("%v", err))
    }

    return &emptypb.Empty{}, nil 
}

func NewServer(l logger.Logger) *Server {
    return &Server {
        Subpub: subpub.NewSubPub(),
        l: l,
    }
}
