package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/calvinbui/teleport-sre-challenge/services/deployments"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/calvinbui/teleport-sre-challenge/proto/gen"
)

func (s *server) ServeGRPC() error {
	listener, err := net.Listen("tcp", s.config.GRPCAddress)
	if err != nil {
		return fmt.Errorf("Error creating listener on %s: %v", s.config.GRPCAddress, err)
	}

	// setup mTLS
	tlsConfig, err := mTLSConfig(s.config.CertFile, s.config.KeyFile, s.config.CACertFile)
	if err != nil {
		return fmt.Errorf("Error configuring gRPC TLS: %w", err)
	}

	// configure the grpc server
	opts := []grpc.ServerOption{
		grpc.Creds(tlsConfig),
	}

	grpcServer := grpc.NewServer(opts...)

	// register the deployment service from our protobuf generated files
	pb.RegisterDeploymentsServiceServer(grpcServer, deployments.New(s.services, s.config))

	// start the grpc server
	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func mTLSConfig(certFile, keyFile, caFile string) (credentials.TransportCredentials, error) {
	// load the server crt and private key
	crt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// create a cert pool and load the ca cert into it
	certPool := x509.NewCertPool()

	chain, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(chain); !ok {
		return nil, fmt.Errorf("Failed to add server CA certificates")
	}

	// create the tls config which requires mtls
	tlsConfig := &tls.Config{
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{crt},
	}

	return credentials.NewTLS(tlsConfig), nil
}
