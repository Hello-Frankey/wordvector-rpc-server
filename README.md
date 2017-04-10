### 简介

词向量接口服务器可以将word2vec生成的词向量文件加载到内存并通过RPC调用的方式提供词向量的查询接口。目前开放的接口仅有一个GetVector，该接口接收一个词作为输入，返回该词的词向量以及其在词典中的位置索引。

### 接口服务器
#### 服务器构建
```
git clone https://github.com/Hello-Frankey/wordvector-rpc-server.git
cd wordvector-rpc-server/
go build
./wordvector-rpc-server -h
Usage of ./wordvector-rpc-server:
  -binary
        待加载词向量文件是否为二进制文件
  -filepath string
        词向量文件路径
  -port int
        服务器监听的端口号 (default 50051)
  -size int
        词向量维度
```
#### 服务器启动
```
加载文本格式的词向量文件，需要提供词向量维度参数
./wordvector-rpc-server -filepath sample_vector.txt -size 200
加载二进制格式的词向量文件，无需提供词向量维度参数
./wordvector-rpc-server -filepath sample_vector.bin -binary
```

### 接口说明

词向量接口服务器基于gRPC实现接口调用，描述接口的proto文件如下。当前仅开放一个查询接口GetVector，请求消息包含一个字段word（查询词），响应消息包含三个字段word（查询词）、index（查询词在词表中的位置索引）和features（查询词的向量值）。
```
syntax = "proto3";
option java_multiple_files = true;
option java_package = "wordvector";
option java_outer_classname = "WordVectorProto";
package wordvector;
// The wordvector service definition.
service WordVector {
  // Sends a word vector of given word.
  rpc GetVector (GetVectorRequest) returns (GetVectorReply) {}
}
// The request message containing the word.
message GetVectorRequest {
  string word = 1;
}
// The response message containing the word, word index in the
// vocabulary and its vector.
message GetVectorReply {
  string word = 1;
  int64 index = 2;
  repeated float features = 3;
}
```

### 客户端调用代码示例
#### python客户端调用
```
import grpc
import wordvector_pb2
def run():
    channel = grpc.insecure_channel('localhost:5051')
    stub = wordvector_pb2.WordVectorStub(channel)
    response = stub.GetVector(wordvector_pb2.GetVectorRequest(word='银行担保'))
    print response.index
    print response.features
if __name__ == '__main__':
    run()
```
#### go客户端调用
```
package main
import (
        "log"
        "golang.org/x/net/context"
        "google.golang.org/grpc"
        pb "wordvector/wordvector"
)
 
func main() {
        conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
        if err != nil {
                log.Fatalf("did not connect: %v", err)
        }
        defer conn.Close()
        c := pb.NewWordVectorClient(conn)
        r, err := c.GetVector(context.Background(), &pb.GetVectorRequest{Word: "银行担保"})
        if err != nil {
                log.Fatalf("could not get word vector: %v", err)
        }
        log.Println("word vector:", r.Word, r.Index, r.Features)
}
```