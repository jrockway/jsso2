/**
 * @fileoverview gRPC-Web generated client stub for foo
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as foo_pb from './foo_pb';


export class NameServiceClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoTransformName = new grpcWeb.AbstractClientBase.MethodInfo(
    foo_pb.TransformNameReply,
    (request: foo_pb.TransformNameRequest) => {
      return request.serializeBinary();
    },
    foo_pb.TransformNameReply.deserializeBinary
  );

  transformName(
    request: foo_pb.TransformNameRequest,
    metadata: grpcWeb.Metadata | null): Promise<foo_pb.TransformNameReply>;

  transformName(
    request: foo_pb.TransformNameRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: foo_pb.TransformNameReply) => void): grpcWeb.ClientReadableStream<foo_pb.TransformNameReply>;

  transformName(
    request: foo_pb.TransformNameRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: foo_pb.TransformNameReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/foo.NameService/TransformName',
        request,
        metadata || {},
        this.methodInfoTransformName,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/foo.NameService/TransformName',
    request,
    metadata || {},
    this.methodInfoTransformName);
  }

}

