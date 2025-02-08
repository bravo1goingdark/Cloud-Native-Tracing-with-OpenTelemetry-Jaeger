package main

import (
	"context"
	"log"
	"net"

	"github.com/bravo1goingdark/jaegaurd/internal/middleware"
	"github.com/bravo1goingdark/jaegaurd/internal/tracing"
	pb "github.com/bravo1goingdark/jaegaurd/proto"
	"go.opentelemetry.io/otel/sdk/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	tr := tracing.Tracer()
	ctx, span := tr.Start(ctx, "SayHello")
	defer span.End()

	log.Printf("📡 Received request: %s", req.Name)
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	// Initialize tracing
	tp, err := tracing.InitTracer("service-a")
	if err != nil {
		log.Fatalf("❌ Failed to initialize tracer: %v", err)
	}
	defer func(tp *trace.TracerProvider, ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("⚠️ Error shutting down tracer: %v", err)
		}
	}(tp, context.Background())

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryTracingInterceptor()),
	)
	pb.RegisterGreeterServer(grpcServer, &greeterServer{})

	reflection.Register(grpcServer)

	log.Println("🚀 Starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to start gRPC server: %v", err)
	}
}
