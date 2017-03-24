package main

import (
	"flag"
	"log"

	pb "wordvector-rpc-server/wordvector"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	address = flag.String("address", "localhost:50051", "接口服务器地址")
	word    = flag.String("word", "", "查询词")
)

func main() {
	flag.Parse()
	if *address == "" {
		log.Fatalln("missing rpc server address")
	}
	if *word == "" {
		log.Fatalln("missing word to query")
	}
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faield to connect to rpc server: %v", err)
	}
	defer conn.Close()
	c := pb.NewWordVectorClient(conn)
	r, err := c.GetVector(context.Background(), &pb.GetVectorRequest{Word: *word})
	if err != nil {
		log.Fatalf("could not get word vector: %v", err)
	}
	log.Println("word vector:", r.Word, r.Index, r.Features)
}
