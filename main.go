package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "wordvector-rpc-server/wordvector"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	filepath = flag.String("filepath", "", "词向量文件路径")
	port     = flag.Int("port", 50051, "服务器监听的端口号")
	binary   = flag.Bool("binary", false, "待加载词向量文件是否为二进制文件")
	size     = flag.Int("size", 0, "词向量维度")
)

var (
	vocabulary map[string]int64
	vectors    [][]float32
)

// server is used to implement word vector server.
type server struct{}

// GetVector implements the get vector method.
func (s *server) GetVector(ctx context.Context, in *pb.GetVectorRequest) (*pb.GetVectorReply, error) {
	if idx, ok := vocabulary[in.Word]; ok {
		v := &pb.GetVectorReply{
			Word:     in.Word,
			Index:    idx,
			Features: vectors[idx],
		}
		return v, nil
	}
	return &pb.GetVectorReply{Word: "", Index: -1, Features: []float32{}}, nil
}

// isFileExists checks if the given path exists and it's a file.
func isFileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		if f.IsDir() {
			return false, nil
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// initConfig initializes the server config with the command line parameter
// and checks if the input parameters are valid.
func initConfig() error {
	flag.Parse()
	if *filepath == "" {
		return errors.New("missing word vector file")
	}
	if ok, _ := isFileExists(*filepath); !ok {
		return errors.New("word vector file not exist")
	}
	if *port <= 0 {
		return errors.New("missing valid port number")
	}
	if *size <= 0 {
		return errors.New("missing valid vector size")
	}
	return nil
}

func loadBinaryWordVector(file string, size int) error {
	return nil
}

func loadTextWordVector(file string, size int) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer in.Close()
	vocabulary = make(map[string]int64)
	var index int64
	startTime := time.Now()
	for {
		var word string
		_, err := fmt.Fscanf(in, "%s", &word)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if word == "" {
			break
		}
		vocabulary[word] = index
		features := make([]float32, size)
		for i := 0; i < size; i++ {
			fmt.Fscanf(in, "%f", &features[i])
		}
		vectors = append(vectors, features)
		index++
	}
	endTime := time.Now()
	log.Printf("word vector loaded, time consumed: %fs, vocabulary size: %d\n",
		endTime.Sub(startTime).Seconds(), index)
	return nil
}

// loadWordVector loads the word vectors from the given file.
func loadWordVector(file string, size int, binary bool) error {
	if binary {
		return loadBinaryWordVector(file, size)
	}
	return loadTextWordVector(file, size)
}

func main() {
	err := initConfig()
	if err != nil {
		log.Fatalf("failed to init word vector server, %v\n", err)
	}
	err = loadWordVector(*filepath, *size, *binary)
	if err != nil {
		log.Fatalf("failed to load word vector, %v\n", err)
	}
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, %v\n", *port, err)
	} else {
		log.Printf("word vector server listening on port: %d\n", *port)
	}
	s := grpc.NewServer()
	pb.RegisterWordVectorServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve request, %v\n", err)
	}
}
