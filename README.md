### 简介

词向量接口服务器可以将word2vec生成的词向量文件加载到内存并通过RPC调用的方式提供词向量的查询接口。目前开放的接口仅有一个GetVector，该接口接收一个词作为输入，返回该词的词向量以及其在词典中的位置索引。

### 接口服务器
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

./wordvector-rpc-server -filepath sample_vector.txt -size 200
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
#### 调用结果示例
```
./client -address localhost:5051 -word 银行担保
2017/04/10 17:22:34 word vector: 银行担保 9 [0.112975 0.038483 0.004365 0.040772 -0.147109 0.00387 -0.01505 -0.077319 -0.111126 0.020161 -0.094328 -0.040327 0.021586 -0.077863 0.102627 -0.100328 -0.015392 0.000837 -0.012477 0.033762 -0.027628 0.010177 0.02421 0.139064 -0.004234 0.03329 0.110185 -0.044716 0.071849 0.044547 0.046687 0.037414 -0.152806 -0.008204 -0.068997 -0.030079 0.064869 -0.069475 -0.222064 0.087641 0.015623 -0.015399 -0.028867 0.076204 -0.016811 -0.043338 -0.016375 0.028917 0.073044 -0.076721 -0.011282 -0.032553 -0.060786 -0.071592 -0.09752 -0.052884 -0.008509 0.112547 0.00121 -0.041229 0.043674 0.03513 0.087865 -0.122061 -0.007484 0.026341 -0.048112 -0.076131 0.101304 -0.154075 0.008777 -0.007455 0.005158 -0.06447 0.008536 0.073898 0.060431 -0.071171 0.01475 -0.035088 -0.094213 -0.101527 0.056793 0.005966 -0.084752 -0.045544 0.03262 0.093777 -0.181556 0.05288 -0.018762 0.055837 -0.038792 -0.093862 -0.038243 -0.052604 -0.009202 0.059273 0.018261 -0.113781 0.084002 -0.081499 0.091183 -0.031993 0.039732 0.014032 0.138285 -0.032791 0.084675 0.025166 -0.058029 0.007452 -0.020193 -0.07053 0.014399 0.018984 0.028766 -0.125199 -0.028797 -0.041677 -0.026934 0.039082 0.060499 0.063914 0.09968 0.211986 0.098151 0.051541 -0.053377 -0.024781 -0.07679 0.000475 0.027537 0.063294 0.069619 0.121042 -0.03279 -0.14487 -0.012769 -0.136599 0.034089 0.121043 0.066555 -0.019764 -0.053804 -0.040012 -0.044327 0.093057 -0.103222 -0.039553 -0.003958 0.006198 0.015406 0.056962 0.052272 -0.072491 -0.06891 0.132141 -0.048858 0.021628 0.065523 -0.003949 0.0237 0.123135 0.036168 0.072514 -0.031365 -0.057712 0.02702 -0.068734 -0.040546 -0.014049 0.151302 -0.02961 -0.020679 -0.081156 -0.007406 0.116467 0.080552 -0.05146 -0.019194 0.076374 0.082635 0.020141 0.143519 0.040339 0.124316 0.0009 0.003159 0.01807 0.031587 0.050386 -0.026184 -0.05656 0.016478 0.00594 0.032288 -0.054161 -0.151346 -0.088725]
```