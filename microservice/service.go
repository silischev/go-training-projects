package main

import (
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type BizServerHandler interface {
	Check(context.Context, Nothing) Nothing
	Add(context.Context, Nothing) Nothing
	Test(context.Context, Nothing) Nothing
}

type Biz struct {
	rules []AclRule
}

func (b Biz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	log.Println("In Check()")
	return &Nothing{Dummy: true}, nil
}

func (b Biz) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	log.Println("In Add()")
	return &Nothing{Dummy: true}, nil
}

func (b Biz) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return &Nothing{Dummy: true}, nil
}

func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)

	rules, err := CreateRulesFromIncomingMessage([]byte(ACLData))
	if err != nil {
		return err
	}

	biz := Biz{rules}

	go func(ctx context.Context) error {
		lis, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				lis.Close()
				server.Stop()

				return nil
			default:
				RegisterBizServer(server, biz)

				err = server.Serve(lis)
				if err != nil {
					log.Println("Cant serve: ", err)
					return err
				}
			}
		}
	}(ctx)

	return nil
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	bizServer := info.Server.(Biz)

	md, _ := metadata.FromIncomingContext(ctx)
	consumer, ok := md["consumer"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Field not exist")
	}

	hasAccess := hasAccess(strings.Join(consumer, ","), info.FullMethod, bizServer.rules)
	if !hasAccess {
		return nil, status.Errorf(codes.Unauthenticated, "Access denied")
	}

	return handler(ctx, req)
}
