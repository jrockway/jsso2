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

  getUserId(): number;
  setUserId(value: number): Session;

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
    userId: number,
    metadata?: SessionMetadata.AsObject,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    expiresAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

