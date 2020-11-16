import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class User extends jspb.Message {
  getId(): number;
  setId(value: number): User;

  getUsername(): string;
  setUsername(value: string): User;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): User;
  hasCreatedAt(): boolean;
  clearCreatedAt(): User;

  getDisabledAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDisabledAt(value?: google_protobuf_timestamp_pb.Timestamp): User;
  hasDisabledAt(): boolean;
  clearDisabledAt(): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: number,
    username: string,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    disabledAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class SessionMetadata extends jspb.Message {
  getIpAddress(): string;
  setIpAddress(value: string): SessionMetadata;

  getUserAgent(): string;
  setUserAgent(value: string): SessionMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: SessionMetadata): SessionMetadata.AsObject;
  static serializeBinaryToWriter(message: SessionMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionMetadata;
  static deserializeBinaryFromReader(message: SessionMetadata, reader: jspb.BinaryReader): SessionMetadata;
}

export namespace SessionMetadata {
  export type AsObject = {
    ipAddress: string,
    userAgent: string,
  }
}

export class Session extends jspb.Message {
  getId(): Uint8Array | string;
  getId_asU8(): Uint8Array;
  getId_asB64(): string;
  setId(value: Uint8Array | string): Session;

  getUser(): User | undefined;
  setUser(value?: User): Session;
  hasUser(): boolean;
  clearUser(): Session;

  getMetadata(): SessionMetadata | undefined;
  setMetadata(value?: SessionMetadata): Session;
  hasMetadata(): boolean;
  clearMetadata(): Session;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Session;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Session;

  getExpiresAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpiresAt(value?: google_protobuf_timestamp_pb.Timestamp): Session;
  hasExpiresAt(): boolean;
  clearExpiresAt(): Session;

  getTaintsList(): Array<string>;
  setTaintsList(value: Array<string>): Session;
  clearTaintsList(): Session;
  addTaints(value: string, index?: number): Session;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Session.AsObject;
  static toObject(includeInstance: boolean, msg: Session): Session.AsObject;
  static serializeBinaryToWriter(message: Session, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Session;
  static deserializeBinaryFromReader(message: Session, reader: jspb.BinaryReader): Session;
}

export namespace Session {
  export type AsObject = {
    id: Uint8Array | string,
    user?: User.AsObject,
    metadata?: SessionMetadata.AsObject,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    expiresAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    taintsList: Array<string>,
  }
}

export class Credential extends jspb.Message {
  getId(): number;
  setId(value: number): Credential;

  getCredentialId(): Uint8Array | string;
  getCredentialId_asU8(): Uint8Array;
  getCredentialId_asB64(): string;
  setCredentialId(value: Uint8Array | string): Credential;

  getPublicKey(): Uint8Array | string;
  getPublicKey_asU8(): Uint8Array;
  getPublicKey_asB64(): string;
  setPublicKey(value: Uint8Array | string): Credential;

  getUser(): User | undefined;
  setUser(value?: User): Credential;
  hasUser(): boolean;
  clearUser(): Credential;

  getName(): string;
  setName(value: string): Credential;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Credential;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Credential;

  getDeletedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDeletedAt(value?: google_protobuf_timestamp_pb.Timestamp): Credential;
  hasDeletedAt(): boolean;
  clearDeletedAt(): Credential;

  getCreatedBySessionId(): Uint8Array | string;
  getCreatedBySessionId_asU8(): Uint8Array;
  getCreatedBySessionId_asB64(): string;
  setCreatedBySessionId(value: Uint8Array | string): Credential;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Credential.AsObject;
  static toObject(includeInstance: boolean, msg: Credential): Credential.AsObject;
  static serializeBinaryToWriter(message: Credential, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Credential;
  static deserializeBinaryFromReader(message: Credential, reader: jspb.BinaryReader): Credential;
}

export namespace Credential {
  export type AsObject = {
    id: number,
    credentialId: Uint8Array | string,
    publicKey: Uint8Array | string,
    user?: User.AsObject,
    name: string,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    deletedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    createdBySessionId: Uint8Array | string,
  }
}

