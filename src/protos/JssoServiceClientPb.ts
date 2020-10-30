/**
 * @fileoverview gRPC-Web generated client stub for jsso
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as jsso_pb from './jsso_pb';


export class UserClient {
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

  methodInfoEdit = new grpcWeb.AbstractClientBase.MethodInfo(
    jsso_pb.EditUserReply,
    (request: jsso_pb.EditUserRequest) => {
      return request.serializeBinary();
    },
    jsso_pb.EditUserReply.deserializeBinary
  );

  edit(
    request: jsso_pb.EditUserRequest,
    metadata: grpcWeb.Metadata | null): Promise<jsso_pb.EditUserReply>;

  edit(
    request: jsso_pb.EditUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: jsso_pb.EditUserReply) => void): grpcWeb.ClientReadableStream<jsso_pb.EditUserReply>;

  edit(
    request: jsso_pb.EditUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: jsso_pb.EditUserReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/jsso.User/Edit',
        request,
        metadata || {},
        this.methodInfoEdit,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/jsso.User/Edit',
    request,
    metadata || {},
    this.methodInfoEdit);
  }

  methodInfoGenerateEnrollmentLink = new grpcWeb.AbstractClientBase.MethodInfo(
    jsso_pb.GenerateEnrollmentLinkReply,
    (request: jsso_pb.GenerateEnrollmentLinkRequest) => {
      return request.serializeBinary();
    },
    jsso_pb.GenerateEnrollmentLinkReply.deserializeBinary
  );

  generateEnrollmentLink(
    request: jsso_pb.GenerateEnrollmentLinkRequest,
    metadata: grpcWeb.Metadata | null): Promise<jsso_pb.GenerateEnrollmentLinkReply>;

  generateEnrollmentLink(
    request: jsso_pb.GenerateEnrollmentLinkRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: jsso_pb.GenerateEnrollmentLinkReply) => void): grpcWeb.ClientReadableStream<jsso_pb.GenerateEnrollmentLinkReply>;

  generateEnrollmentLink(
    request: jsso_pb.GenerateEnrollmentLinkRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: jsso_pb.GenerateEnrollmentLinkReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/jsso.User/GenerateEnrollmentLink',
        request,
        metadata || {},
        this.methodInfoGenerateEnrollmentLink,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/jsso.User/GenerateEnrollmentLink',
    request,
    metadata || {},
    this.methodInfoGenerateEnrollmentLink);
  }

}

export class LoginClient {
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

  methodInfoStart = new grpcWeb.AbstractClientBase.MethodInfo(
    jsso_pb.StartLoginReply,
    (request: jsso_pb.StartLoginRequest) => {
      return request.serializeBinary();
    },
    jsso_pb.StartLoginReply.deserializeBinary
  );

  start(
    request: jsso_pb.StartLoginRequest,
    metadata: grpcWeb.Metadata | null): Promise<jsso_pb.StartLoginReply>;

  start(
    request: jsso_pb.StartLoginRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: jsso_pb.StartLoginReply) => void): grpcWeb.ClientReadableStream<jsso_pb.StartLoginReply>;

  start(
    request: jsso_pb.StartLoginRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: jsso_pb.StartLoginReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/jsso.Login/Start',
        request,
        metadata || {},
        this.methodInfoStart,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/jsso.Login/Start',
    request,
    metadata || {},
    this.methodInfoStart);
  }

}

export class EnrollmentClient {
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

  methodInfoStart = new grpcWeb.AbstractClientBase.MethodInfo(
    jsso_pb.StartEnrollmentReply,
    (request: jsso_pb.StartEnrollmentRequest) => {
      return request.serializeBinary();
    },
    jsso_pb.StartEnrollmentReply.deserializeBinary
  );

  start(
    request: jsso_pb.StartEnrollmentRequest,
    metadata: grpcWeb.Metadata | null): Promise<jsso_pb.StartEnrollmentReply>;

  start(
    request: jsso_pb.StartEnrollmentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: jsso_pb.StartEnrollmentReply) => void): grpcWeb.ClientReadableStream<jsso_pb.StartEnrollmentReply>;

  start(
    request: jsso_pb.StartEnrollmentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: jsso_pb.StartEnrollmentReply) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/jsso.Enrollment/Start',
        request,
        metadata || {},
        this.methodInfoStart,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/jsso.Enrollment/Start',
    request,
    metadata || {},
    this.methodInfoStart);
  }

}

