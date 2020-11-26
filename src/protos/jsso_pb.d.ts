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

  getRedirectToken(): string;
  setRedirectToken(value: string): FinishLoginRequest;

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
    redirectToken: string,
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

export class AuthorizeHTTPRequest extends jspb.Message {
  getRequestMethod(): string;
  setRequestMethod(value: string): AuthorizeHTTPRequest;

  getRequestUri(): string;
  setRequestUri(value: string): AuthorizeHTTPRequest;

  getRequestId(): string;
  setRequestId(value: string): AuthorizeHTTPRequest;

  getAuthorizationHeadersList(): Array<string>;
  setAuthorizationHeadersList(value: Array<string>): AuthorizeHTTPRequest;
  clearAuthorizationHeadersList(): AuthorizeHTTPRequest;
  addAuthorizationHeaders(value: string, index?: number): AuthorizeHTTPRequest;

  getCookiesList(): Array<string>;
  setCookiesList(value: Array<string>): AuthorizeHTTPRequest;
  clearCookiesList(): AuthorizeHTTPRequest;
  addCookies(value: string, index?: number): AuthorizeHTTPRequest;

  getIpAddress(): string;
  setIpAddress(value: string): AuthorizeHTTPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthorizeHTTPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AuthorizeHTTPRequest): AuthorizeHTTPRequest.AsObject;
  static serializeBinaryToWriter(message: AuthorizeHTTPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthorizeHTTPRequest;
  static deserializeBinaryFromReader(message: AuthorizeHTTPRequest, reader: jspb.BinaryReader): AuthorizeHTTPRequest;
}

export namespace AuthorizeHTTPRequest {
  export type AsObject = {
    requestMethod: string,
    requestUri: string,
    requestId: string,
    authorizationHeadersList: Array<string>,
    cookiesList: Array<string>,
    ipAddress: string,
  }
}

export class Allow extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): Allow;

  getGroupsList(): Array<string>;
  setGroupsList(value: Array<string>): Allow;
  clearGroupsList(): Allow;
  addGroups(value: string, index?: number): Allow;

  getBearerToken(): string;
  setBearerToken(value: string): Allow;

  getAddHeadersList(): Array<types_pb.Header>;
  setAddHeadersList(value: Array<types_pb.Header>): Allow;
  clearAddHeadersList(): Allow;
  addAddHeaders(value?: types_pb.Header, index?: number): types_pb.Header;

  getAppendHeadersList(): Array<types_pb.Header>;
  setAppendHeadersList(value: Array<types_pb.Header>): Allow;
  clearAppendHeadersList(): Allow;
  addAppendHeaders(value?: types_pb.Header, index?: number): types_pb.Header;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Allow.AsObject;
  static toObject(includeInstance: boolean, msg: Allow): Allow.AsObject;
  static serializeBinaryToWriter(message: Allow, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Allow;
  static deserializeBinaryFromReader(message: Allow, reader: jspb.BinaryReader): Allow;
}

export namespace Allow {
  export type AsObject = {
    username: string,
    groupsList: Array<string>,
    bearerToken: string,
    addHeadersList: Array<types_pb.Header.AsObject>,
    appendHeadersList: Array<types_pb.Header.AsObject>,
  }
}

export class Deny extends jspb.Message {
  getReason(): string;
  setReason(value: string): Deny;

  getRedirect(): Deny.Redirect | undefined;
  setRedirect(value?: Deny.Redirect): Deny;
  hasRedirect(): boolean;
  clearRedirect(): Deny;

  getResponse(): Deny.Response | undefined;
  setResponse(value?: Deny.Response): Deny;
  hasResponse(): boolean;
  clearResponse(): Deny;

  getDestinationCase(): Deny.DestinationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Deny.AsObject;
  static toObject(includeInstance: boolean, msg: Deny): Deny.AsObject;
  static serializeBinaryToWriter(message: Deny, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Deny;
  static deserializeBinaryFromReader(message: Deny, reader: jspb.BinaryReader): Deny;
}

export namespace Deny {
  export type AsObject = {
    reason: string,
    redirect?: Deny.Redirect.AsObject,
    response?: Deny.Response.AsObject,
  }

  export class Redirect extends jspb.Message {
    getRedirectUrl(): string;
    setRedirectUrl(value: string): Redirect;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Redirect.AsObject;
    static toObject(includeInstance: boolean, msg: Redirect): Redirect.AsObject;
    static serializeBinaryToWriter(message: Redirect, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Redirect;
    static deserializeBinaryFromReader(message: Redirect, reader: jspb.BinaryReader): Redirect;
  }

  export namespace Redirect {
    export type AsObject = {
      redirectUrl: string,
    }
  }


  export class Response extends jspb.Message {
    getContentType(): string;
    setContentType(value: string): Response;

    getBody(): string;
    setBody(value: string): Response;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Response.AsObject;
    static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
    static serializeBinaryToWriter(message: Response, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Response;
    static deserializeBinaryFromReader(message: Response, reader: jspb.BinaryReader): Response;
  }

  export namespace Response {
    export type AsObject = {
      contentType: string,
      body: string,
    }
  }


  export enum DestinationCase { 
    DESTINATION_NOT_SET = 0,
    REDIRECT = 2,
    RESPONSE = 3,
  }
}

export class AuthorizeHTTPReply extends jspb.Message {
  getAllow(): Allow | undefined;
  setAllow(value?: Allow): AuthorizeHTTPReply;
  hasAllow(): boolean;
  clearAllow(): AuthorizeHTTPReply;

  getDeny(): Deny | undefined;
  setDeny(value?: Deny): AuthorizeHTTPReply;
  hasDeny(): boolean;
  clearDeny(): AuthorizeHTTPReply;

  getDecisionCase(): AuthorizeHTTPReply.DecisionCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthorizeHTTPReply.AsObject;
  static toObject(includeInstance: boolean, msg: AuthorizeHTTPReply): AuthorizeHTTPReply.AsObject;
  static serializeBinaryToWriter(message: AuthorizeHTTPReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthorizeHTTPReply;
  static deserializeBinaryFromReader(message: AuthorizeHTTPReply, reader: jspb.BinaryReader): AuthorizeHTTPReply;
}

export namespace AuthorizeHTTPReply {
  export type AsObject = {
    allow?: Allow.AsObject,
    deny?: Deny.AsObject,
  }

  export enum DecisionCase { 
    DECISION_NOT_SET = 0,
    ALLOW = 1,
    DENY = 2,
  }
}

