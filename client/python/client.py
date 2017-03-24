#!/usr/bin/env python
# -*- encoding: utf-8 -*-
"""
python client example of word vector.
"""
import argparse
import grpc
import wordvector_pb2

def run():
    """
    run function
    """
    parser = argparse.ArgumentParser(description="word vector python client")
    parser.add_argument("--address", default="localhost:50051", help="接口服务器地址")
    parser.add_argument("--word", help="查询词")
    arguments = parser.parse_args()
    if not arguments.address:
        print "missing rpc server address"
        exit(0)
    if not arguments.word:
        print "missing word to query"
        exit(0)
    channel = grpc.insecure_channel(arguments.address)
    stub = wordvector_pb2.WordVectorStub(channel)
    response = stub.GetVector(wordvector_pb2.GetVectorRequest(word=arguments.word))
    print response.word, response.index, response.features

if __name__ == '__main__':
    run()
