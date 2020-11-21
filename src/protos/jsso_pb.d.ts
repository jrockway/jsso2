import * as jspb from 'google-protobuf'

import * as types_pb from './types_pb';
import * as webauthn_pb from './webauthn_pb';


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
  getUsername(): string;
  setUsername(value: string): StartLoginRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartLoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartLoginRequest): StartLoginRequest.AsObject;
  static serializeBinaryToWriter(message: StartLoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartLoginRequest;
  static deserializeBinaryFromReader(message: StartLoginRequest, reader: jspb.BinaryReader): StartLoginRequest;
}

export namespace StartLoginRequest {
  export type AsObject = {
    username: string,
  }
}

export class StartLoginReply extends jspb.Message {
  getCredentialRequestOptions(): webauthn_pb.PublicKeyCredentialRequestOptions | undefined;
  setCredentialRequestOptions(value?: webauthn_pb.PublicKeyCredentialRequestOptions): StartLoginReply;
  hasCredentialRequestOptions(): boolean;
  clearCredentialRequestOptions(): StartLoginReply;

  getToken(): string;
  setToken(value: string): StartLoginReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartLoginReply.AsObject;
  static toObject(includeInstance: boolean, msg: StartLoginReply): StartLoginReply.AsObject;
  static serializeBinaryToWriter(message: StartLoginReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartLoginReply;
  static deserializeBinaryFromReader(message: StartLoginReply, reader: jspb.BinaryReader): StartLoginReply;
}

export namespace StartLoginReply {
  export type AsObject = {
    credentialRequestOptions?: webauthn_pb.PublicKeyCredentialRequestOptions.AsObject,
    token: string,
  }
}

export class FinishLoginRequest extends jspb.Message {
  getCredential(): webauthn_pb.PublicKeyCredential | undefined;
  setCredential(value?: webauthn_pb.PublicKeyCredential): FinishLoginRequest;
  hasCredential(): boolean;
  clearCredential(): FinishLoginRequest;

  getError(): string;
  setError(value: string): FinishLoginRequest;

  getRedirectTo(): string;
  setRedirectTo(value: string): FinishLoginRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinishLoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FinishLoginRequest): FinishLoginRequest.AsObject;
  static serializeBinaryToWriter(message: FinishLoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinishLoginRequest;
  static deserializeBinaryFromReader(message: FinishLoginRequest, reader: jspb.BinaryReader): FinishLoginRequest;
}

export namespace FinishLoginRequest {
  export type AsObject = {
    credential?: webauthn_pb.PublicKeyCredential.AsObject,
    error: string,
    redirectTo: string,
  }
}

export class FinishLoginReply extends jspb.Message {
  getRedirectUrl(): string;
  setRedirectUrl(value: string): FinishLoginReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinishLoginReply.AsObject;
  static toObject(includeInstance: boolean, msg: FinishLoginReply): FinishLoginReply.AsObject;
  static serializeBinaryToWriter(message: FinishLoginReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinishLoginReply;
  static deserializeBinaryFromReader(message: FinishLoginReply, reader: jspb.BinaryReader): FinishLoginReply;
}

export namespace FinishLoginReply {
  export type AsObject = {
    redirectUrl: string,
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

  getCredentialCreationOptions(): webauthn_pb.PublicKeyCredentialCreationOptions | undefined;
  setCredentialCreationOptions(value?: webauthn_pb.PublicKeyCredentialCreationOptions): StartEnrollmentReply;
  hasCredentialCreationOptions(): boolean;
  clearCredentialCreationOptions(): StartEnrollmentReply;

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
    credentialCreationOptions?: webauthn_pb.PublicKeyCredentialCreationOptions.AsObject,
  }
}

export class FinishEnrollmentRequest extends jspb.Message {
  getCredential(): webauthn_pb.PublicKeyCredential | undefined;
  setCredential(value?: webauthn_pb.PublicKeyCredential): FinishEnrollmentRequest;
  hasCredential(): boolean;
  clearCredential(): FinishEnrollmentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinishEnrollmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FinishEnrollmentRequest): FinishEnrollmentRequest.AsObject;
  static serializeBinaryToWriter(message: FinishEnrollmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinishEnrollmentRequest;
  static deserializeBinaryFromReader(message: FinishEnrollmentRequest, reader: jspb.BinaryReader): FinishEnrollmentRequest;
}

export namespace FinishEnrollmentRequest {
  export type AsObject = {
    credential?: webauthn_pb.PublicKeyCredential.AsObject,
  }
}

export class FinishEnrollmentReply extends jspb.Message {
  getLoginUrl(): string;
  setLoginUrl(value: string): FinishEnrollmentReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinishEnrollmentReply.AsObject;
  static toObject(includeInstance: boolean, msg: FinishEnrollmentReply): FinishEnrollmentReply.AsObject;
  static serializeBinaryToWriter(message: FinishEnrollmentReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinishEnrollmentReply;
  static deserializeBinaryFromReader(message: FinishEnrollmentReply, reader: jspb.BinaryReader): FinishEnrollmentReply;
}

export namespace FinishEnrollmentReply {
  export type AsObject = {
    loginUrl: string,
  }
}

export class WhoAmIRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WhoAmIRequest.AsObject;
  static toObject(includeInstance: boolean, msg: WhoAmIRequest): WhoAmIRequest.AsObject;
  static serializeBinaryToWriter(message: WhoAmIRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WhoAmIRequest;
  static deserializeBinaryFromReader(message: WhoAmIRequest, reader: jspb.BinaryReader): WhoAmIRequest;
}

export namespace WhoAmIRequest {
  export type AsObject = {
  }
}

export class WhoAmIReply extends jspb.Message {
  getUser(): types_pb.User | undefined;
  setUser(value?: types_pb.User): WhoAmIReply;
  hasUser(): boolean;
  clearUser(): WhoAmIReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WhoAmIReply.AsObject;
  static toObject(includeInstance: boolean, msg: WhoAmIReply): WhoAmIReply.AsObject;
  static serializeBinaryToWriter(message: WhoAmIReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WhoAmIReply;
  static deserializeBinaryFromReader(message: WhoAmIReply, reader: jspb.BinaryReader): WhoAmIReply;
}

export namespace WhoAmIReply {
  export type AsObject = {
    user?: types_pb.User.AsObject,
  }
}

