import * as jspb from 'google-protobuf'

import * as types_pb from './types_pb';


export class AddUserRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): AddUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddUserRequest): AddUserRequest.AsObject;
  static serializeBinaryToWriter(message: AddUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUserRequest;
  static deserializeBinaryFromReader(message: AddUserRequest, reader: jspb.BinaryReader): AddUserRequest;
}

export namespace AddUserRequest {
  export type AsObject = {
    username: string,
  }
}

export class AddUserReply extends jspb.Message {
  getUser(): types_pb.User | undefined;
  setUser(value?: types_pb.User): AddUserReply;
  hasUser(): boolean;
  clearUser(): AddUserReply;

  getEnrollmentToken(): string;
  setEnrollmentToken(value: string): AddUserReply;

  getEnrollmentUrl(): string;
  setEnrollmentUrl(value: string): AddUserReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUserReply.AsObject;
  static toObject(includeInstance: boolean, msg: AddUserReply): AddUserReply.AsObject;
  static serializeBinaryToWriter(message: AddUserReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUserReply;
  static deserializeBinaryFromReader(message: AddUserReply, reader: jspb.BinaryReader): AddUserReply;
}

export namespace AddUserReply {
  export type AsObject = {
    user?: types_pb.User.AsObject,
    enrollmentToken: string,
    enrollmentUrl: string,
  }
}

export class StartLoginRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartLoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartLoginRequest): StartLoginRequest.AsObject;
  static serializeBinaryToWriter(message: StartLoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartLoginRequest;
  static deserializeBinaryFromReader(message: StartLoginRequest, reader: jspb.BinaryReader): StartLoginRequest;
}

export namespace StartLoginRequest {
  export type AsObject = {
  }
}

export class StartLoginReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartLoginReply.AsObject;
  static toObject(includeInstance: boolean, msg: StartLoginReply): StartLoginReply.AsObject;
  static serializeBinaryToWriter(message: StartLoginReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartLoginReply;
  static deserializeBinaryFromReader(message: StartLoginReply, reader: jspb.BinaryReader): StartLoginReply;
}

export namespace StartLoginReply {
  export type AsObject = {
  }
}

export class StartEnrollmentRequest extends jspb.Message {
  getEnrollmentToken(): string;
  setEnrollmentToken(value: string): StartEnrollmentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartEnrollmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartEnrollmentRequest): StartEnrollmentRequest.AsObject;
  static serializeBinaryToWriter(message: StartEnrollmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartEnrollmentRequest;
  static deserializeBinaryFromReader(message: StartEnrollmentRequest, reader: jspb.BinaryReader): StartEnrollmentRequest;
}

export namespace StartEnrollmentRequest {
  export type AsObject = {
    enrollmentToken: string,
  }
}

export class StartEnrollmentReply extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartEnrollmentReply.AsObject;
  static toObject(includeInstance: boolean, msg: StartEnrollmentReply): StartEnrollmentReply.AsObject;
  static serializeBinaryToWriter(message: StartEnrollmentReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartEnrollmentReply;
  static deserializeBinaryFromReader(message: StartEnrollmentReply, reader: jspb.BinaryReader): StartEnrollmentReply;
}

export namespace StartEnrollmentReply {
  export type AsObject = {
  }
}

