import * as jspb from 'google-protobuf'

import * as types_pb from './types_pb';


export class EditUserRequest extends jspb.Message {
  getUser(): types_pb.User | undefined;
  setUser(value?: types_pb.User): EditUserRequest;
  hasUser(): boolean;
  clearUser(): EditUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EditUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EditUserRequest): EditUserRequest.AsObject;
  static serializeBinaryToWriter(message: EditUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EditUserRequest;
  static deserializeBinaryFromReader(message: EditUserRequest, reader: jspb.BinaryReader): EditUserRequest;
}

export namespace EditUserRequest {
  export type AsObject = {
    user?: types_pb.User.AsObject,
  }
}

export class EditUserReply extends jspb.Message {
  getUser(): types_pb.User | undefined;
  setUser(value?: types_pb.User): EditUserReply;
  hasUser(): boolean;
  clearUser(): EditUserReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EditUserReply.AsObject;
  static toObject(includeInstance: boolean, msg: EditUserReply): EditUserReply.AsObject;
  static serializeBinaryToWriter(message: EditUserReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EditUserReply;
  static deserializeBinaryFromReader(message: EditUserReply, reader: jspb.BinaryReader): EditUserReply;
}

export namespace EditUserReply {
  export type AsObject = {
    user?: types_pb.User.AsObject,
  }
}

export class GenerateEnrollmentLinkRequest extends jspb.Message {
  getTarget(): types_pb.User | undefined;
  setTarget(value?: types_pb.User): GenerateEnrollmentLinkRequest;
  hasTarget(): boolean;
  clearTarget(): GenerateEnrollmentLinkRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateEnrollmentLinkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateEnrollmentLinkRequest): GenerateEnrollmentLinkRequest.AsObject;
  static serializeBinaryToWriter(message: GenerateEnrollmentLinkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateEnrollmentLinkRequest;
  static deserializeBinaryFromReader(message: GenerateEnrollmentLinkRequest, reader: jspb.BinaryReader): GenerateEnrollmentLinkRequest;
}

export namespace GenerateEnrollmentLinkRequest {
  export type AsObject = {
    target?: types_pb.User.AsObject,
  }
}

export class GenerateEnrollmentLinkReply extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): GenerateEnrollmentLinkReply;

  getToken(): string;
  setToken(value: string): GenerateEnrollmentLinkReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateEnrollmentLinkReply.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateEnrollmentLinkReply): GenerateEnrollmentLinkReply.AsObject;
  static serializeBinaryToWriter(message: GenerateEnrollmentLinkReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateEnrollmentLinkReply;
  static deserializeBinaryFromReader(message: GenerateEnrollmentLinkReply, reader: jspb.BinaryReader): GenerateEnrollmentLinkReply;
}

export namespace GenerateEnrollmentLinkReply {
  export type AsObject = {
    url: string,
    token: string,
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
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartEnrollmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartEnrollmentRequest): StartEnrollmentRequest.AsObject;
  static serializeBinaryToWriter(message: StartEnrollmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartEnrollmentRequest;
  static deserializeBinaryFromReader(message: StartEnrollmentRequest, reader: jspb.BinaryReader): StartEnrollmentRequest;
}

export namespace StartEnrollmentRequest {
  export type AsObject = {
  }
}

export class StartEnrollmentReply extends jspb.Message {
  getUser(): types_pb.User | undefined;
  setUser(value?: types_pb.User): StartEnrollmentReply;
  hasUser(): boolean;
  clearUser(): StartEnrollmentReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartEnrollmentReply.AsObject;
  static toObject(includeInstance: boolean, msg: StartEnrollmentReply): StartEnrollmentReply.AsObject;
  static serializeBinaryToWriter(message: StartEnrollmentReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartEnrollmentReply;
  static deserializeBinaryFromReader(message: StartEnrollmentReply, reader: jspb.BinaryReader): StartEnrollmentReply;
}

export namespace StartEnrollmentReply {
  export type AsObject = {
    user?: types_pb.User.AsObject,
  }
}

