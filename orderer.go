package fabric_sdk_go

import (
	"time"

	"fmt"
	"io"
	"strings"

	"github.com/hyperledger/fabric/protos/common"
	ab "github.com/hyperledger/fabric/protos/orderer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	config "sk-git.securekey.com/vme/fabric-sdk-go/config"
)

type Orderer struct {
	Url            string
	GrpcDialOption []grpc.DialOption
}

func CreateNewOrderer(url string) *Orderer {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	if config.IsTlsEnabled() {
		creds := credentials.NewClientTLSFromCert(config.GetTlsCACertPool(), config.GetTlsServerHostOverride())
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return &Orderer{Url: url, GrpcDialOption: opts}
}

func (o *Orderer) SendBroadcast(envelope *common.Envelope) error {
	conn, err := grpc.Dial(o.Url, o.GrpcDialOption...)
	if err != nil {
		return err
	}
	defer conn.Close()

	broadcastStream, err := ab.NewAtomicBroadcastClient(conn).Broadcast(context.Background())
	if err != nil {
		return fmt.Errorf("Error Create NewAtomicBroadcastClient %v", err)
	}
	done := make(chan bool)
	var broadcastErr error
	go func() {
		for {
			broadcastResponse, err := broadcastStream.Recv()
			logger.Debugf("Orderer.broadcastStream - response:%v, error:%v\n", broadcastResponse, err)
			if err != nil {
				if strings.Contains(err.Error(), io.EOF.Error()) {
					done <- true
					return
				}
				broadcastErr = fmt.Errorf("Error broadcast respone : %v\n", err)
				continue
			}
			if broadcastResponse.Status != common.Status_SUCCESS {
				broadcastErr = fmt.Errorf("broadcast respone is not success : %v\n", broadcastResponse.Status)
			}
		}
	}()
	if err := broadcastStream.Send(envelope); err != nil {
		return fmt.Errorf("Failed to send a envelope to orderer: %v", err)
	}
	broadcastStream.CloseSend()
	<-done
	return broadcastErr
}
