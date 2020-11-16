import * as jspb from 'google-protobuf'

import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb';


export class PublicKeyCredentialCreationOptions extends jspb.Message {
  getAttestation(): PublicKeyCredentialCreationOptions.AttestationConveyancePreference;
  setAttestation(value: PublicKeyCredentialCreationOptions.AttestationConveyancePreference): PublicKeyCredentialCreationOptions;

  getAuthenticatorSelection(): AuthenticatorSelectionCriteria | undefined;
  setAuthenticatorSelection(value?: AuthenticatorSelectionCriteria): PublicKeyCredentialCreationOptions;
  hasAuthenticatorSelection(): boolean;
  clearAuthenticatorSelection(): PublicKeyCredentialCreationOptions;

  getChallenge(): Uint8Array | string;
  getChallenge_asU8(): Uint8Array;
  getChallenge_asB64(): string;
  setChallenge(value: Uint8Array | string): PublicKeyCredentialCreationOptions;

  getExcludeCredentialsList(): Array<PublicKeyCredentialDescriptor>;
  setExcludeCredentialsList(value: Array<PublicKeyCredentialDescriptor>): PublicKeyCredentialCreationOptions;
  clearExcludeCredentialsList(): PublicKeyCredentialCreationOptions;
  addExcludeCredentials(value?: PublicKeyCredentialDescriptor, index?: number): PublicKeyCredentialDescriptor;

  getPubKeyCredParamsList(): Array<PublicKeyCredentialParameters>;
  setPubKeyCredParamsList(value: Array<PublicKeyCredentialParameters>): PublicKeyCredentialCreationOptions;
  clearPubKeyCredParamsList(): PublicKeyCredentialCreationOptions;
  addPubKeyCredParams(value?: PublicKeyCredentialParameters, index?: number): PublicKeyCredentialParameters;

  getRp(): PublicKeyCredentialRpEntity | undefined;
  setRp(value?: PublicKeyCredentialRpEntity): PublicKeyCredentialCreationOptions;
  hasRp(): boolean;
  clearRp(): PublicKeyCredentialCreationOptions;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): PublicKeyCredentialCreationOptions;
  hasTimeout(): boolean;
  clearTimeout(): PublicKeyCredentialCreationOptions;

  getUser(): PublicKeyCredentialUserEntity | undefined;
  setUser(value?: PublicKeyCredentialUserEntity): PublicKeyCredentialCreationOptions;
  hasUser(): boolean;
  clearUser(): PublicKeyCredentialCreationOptions;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKeyCredentialCreationOptions.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKeyCredentialCreationOptions): PublicKeyCredentialCreationOptions.AsObject;
  static serializeBinaryToWriter(message: PublicKeyCredentialCreationOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKeyCredentialCreationOptions;
  static deserializeBinaryFromReader(message: PublicKeyCredentialCreationOptions, reader: jspb.BinaryReader): PublicKeyCredentialCreationOptions;
}

export namespace PublicKeyCredentialCreationOptions {
  export type AsObject = {
    attestation: PublicKeyCredentialCreationOptions.AttestationConveyancePreference,
    authenticatorSelection?: AuthenticatorSelectionCriteria.AsObject,
    challenge: Uint8Array | string,
    excludeCredentialsList: Array<PublicKeyCredentialDescriptor.AsObject>,
    pubKeyCredParamsList: Array<PublicKeyCredentialParameters.AsObject>,
    rp?: PublicKeyCredentialRpEntity.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    user?: PublicKeyCredentialUserEntity.AsObject,
  }

  export enum AttestationConveyancePreference { 
    NONE = 0,
    DIRECT = 1,
    INDIRECT = 2,
  }
}

export class PublicKeyCredentialRequestOptions extends jspb.Message {
  getChallenge(): Uint8Array | string;
  getChallenge_asU8(): Uint8Array;
  getChallenge_asB64(): string;
  setChallenge(value: Uint8Array | string): PublicKeyCredentialRequestOptions;

  getAllowedCredentialsList(): Array<PublicKeyCredentialDescriptor>;
  setAllowedCredentialsList(value: Array<PublicKeyCredentialDescriptor>): PublicKeyCredentialRequestOptions;
  clearAllowedCredentialsList(): PublicKeyCredentialRequestOptions;
  addAllowedCredentials(value?: PublicKeyCredentialDescriptor, index?: number): PublicKeyCredentialDescriptor;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): PublicKeyCredentialRequestOptions;
  hasTimeout(): boolean;
  clearTimeout(): PublicKeyCredentialRequestOptions;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKeyCredentialRequestOptions.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKeyCredentialRequestOptions): PublicKeyCredentialRequestOptions.AsObject;
  static serializeBinaryToWriter(message: PublicKeyCredentialRequestOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKeyCredentialRequestOptions;
  static deserializeBinaryFromReader(message: PublicKeyCredentialRequestOptions, reader: jspb.BinaryReader): PublicKeyCredentialRequestOptions;
}

export namespace PublicKeyCredentialRequestOptions {
  export type AsObject = {
    challenge: Uint8Array | string,
    allowedCredentialsList: Array<PublicKeyCredentialDescriptor.AsObject>,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class AuthenticatorSelectionCriteria extends jspb.Message {
  getAuthenticatorAttachment(): AuthenticatorSelectionCriteria.AuthenticatorAttachment;
  setAuthenticatorAttachment(value: AuthenticatorSelectionCriteria.AuthenticatorAttachment): AuthenticatorSelectionCriteria;

  getRequireResidentKey(): boolean;
  setRequireResidentKey(value: boolean): AuthenticatorSelectionCriteria;

  getUserVerification(): AuthenticatorSelectionCriteria.UserVerificationRequirement;
  setUserVerification(value: AuthenticatorSelectionCriteria.UserVerificationRequirement): AuthenticatorSelectionCriteria;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticatorSelectionCriteria.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticatorSelectionCriteria): AuthenticatorSelectionCriteria.AsObject;
  static serializeBinaryToWriter(message: AuthenticatorSelectionCriteria, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticatorSelectionCriteria;
  static deserializeBinaryFromReader(message: AuthenticatorSelectionCriteria, reader: jspb.BinaryReader): AuthenticatorSelectionCriteria;
}

export namespace AuthenticatorSelectionCriteria {
  export type AsObject = {
    authenticatorAttachment: AuthenticatorSelectionCriteria.AuthenticatorAttachment,
    requireResidentKey: boolean,
    userVerification: AuthenticatorSelectionCriteria.UserVerificationRequirement,
  }

  export enum AuthenticatorAttachment { 
    MISSING_AUTHENTICATOR_ATTACHMENT = 0,
    CROSS_PLATFORM = 1,
    PLATFORM = 2,
  }

  export enum UserVerificationRequirement { 
    MISSING_USER_VERIFICATION_REQUIREMENT = 0,
    DISCOURAGED = 1,
    PREFERRED = 2,
    REQUIRED = 3,
  }
}

export class PublicKeyCredentialDescriptor extends jspb.Message {
  getId(): Uint8Array | string;
  getId_asU8(): Uint8Array;
  getId_asB64(): string;
  setId(value: Uint8Array | string): PublicKeyCredentialDescriptor;

  getTransportsList(): Array<PublicKeyCredentialDescriptor.AuthenticatorTransport>;
  setTransportsList(value: Array<PublicKeyCredentialDescriptor.AuthenticatorTransport>): PublicKeyCredentialDescriptor;
  clearTransportsList(): PublicKeyCredentialDescriptor;
  addTransports(value: PublicKeyCredentialDescriptor.AuthenticatorTransport, index?: number): PublicKeyCredentialDescriptor;

  getType(): string;
  setType(value: string): PublicKeyCredentialDescriptor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKeyCredentialDescriptor.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKeyCredentialDescriptor): PublicKeyCredentialDescriptor.AsObject;
  static serializeBinaryToWriter(message: PublicKeyCredentialDescriptor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKeyCredentialDescriptor;
  static deserializeBinaryFromReader(message: PublicKeyCredentialDescriptor, reader: jspb.BinaryReader): PublicKeyCredentialDescriptor;
}

export namespace PublicKeyCredentialDescriptor {
  export type AsObject = {
    id: Uint8Array | string,
    transportsList: Array<PublicKeyCredentialDescriptor.AuthenticatorTransport>,
    type: string,
  }

  export enum AuthenticatorTransport { 
    MISSING_AUTHENTICATOR_TRANSPORT = 0,
    BLE = 1,
    INTERNAL = 2,
    NFC = 3,
    USB = 4,
  }
}

export class PublicKeyCredentialParameters extends jspb.Message {
  getAlg(): number;
  setAlg(value: number): PublicKeyCredentialParameters;

  getType(): string;
  setType(value: string): PublicKeyCredentialParameters;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKeyCredentialParameters.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKeyCredentialParameters): PublicKeyCredentialParameters.AsObject;
  static serializeBinaryToWriter(message: PublicKeyCredentialParameters, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKeyCredentialParameters;
  static deserializeBinaryFromReader(message: PublicKeyCredentialParameters, reader: jspb.BinaryReader): PublicKeyCredentialParameters;
}

export namespace PublicKeyCredentialParameters {
  export type AsObject = {
    alg: number,
    type: string,
  }
}

export class PublicKeyCredentialRpEntity extends jspb.Message {
  getName(): string;
  setName(value: string): PublicKeyCredentialRpEntity;

  getIcon(): string;
  setIcon(value: string): PublicKeyCredentialRpEntity;

  getId(): string;
  setId(value: string): PublicKeyCredentialRpEntity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKeyCredentialRpEntity.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKeyCredentialRpEntity): PublicKeyCredentialRpEntity.AsObject;
  static serializeBinaryToWriter(message: PublicKeyCredentialRpEntity, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKeyCredentialRpEntity;
  static deserializeBinaryFromReader(message: PublicKeyCredentialRpEntity, reader: jspb.BinaryReader): PublicKeyCredentialRpEntity;
}

export namespace PublicKeyCredentialRpEntity {
  export type AsObject = {
    name: string,
    icon: string,
    id: string,
  }
}

export class PublicKeyCredentialUserEntity extends jspb.Message {
  getName(): string;
  setName(value: string): PublicKeyCredentialUserEntity;

  getIcon(): string;
  setIcon(value: string): PublicKeyCredentialUserEntity;

  getDisplayName(): string;
  setDisplayName(value: string): PublicKeyCredentialUserEntity;

  getId(): Uint8Array | string;
  getId_asU8(): Uint8Array;
  getId_asB64(): string;
  setId(value: Uint8Array | string): PublicKeyCredentialUserEntity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKeyCredentialUserEntity.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKeyCredentialUserEntity): PublicKeyCredentialUserEntity.AsObject;
  static serializeBinaryToWriter(message: PublicKeyCredentialUserEntity, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKeyCredentialUserEntity;
  static deserializeBinaryFromReader(message: PublicKeyCredentialUserEntity, reader: jspb.BinaryReader): PublicKeyCredentialUserEntity;
}

export namespace PublicKeyCredentialUserEntity {
  export type AsObject = {
    name: string,
    icon: string,
    displayName: string,
    id: Uint8Array | string,
  }
}

export class PublicKeyCredential extends jspb.Message {
  getId(): string;
  setId(value: string): PublicKeyCredential;

  getType(): string;
  setType(value: string): PublicKeyCredential;

  getClientDataJson(): Uint8Array | string;
  getClientDataJson_asU8(): Uint8Array;
  getClientDataJson_asB64(): string;
  setClientDataJson(value: Uint8Array | string): PublicKeyCredential;

  getAttestationObject(): Uint8Array | string;
  getAttestationObject_asU8(): Uint8Array;
  getAttestationObject_asB64(): string;
  setAttestationObject(value: Uint8Array | string): PublicKeyCredential;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKeyCredential.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKeyCredential): PublicKeyCredential.AsObject;
  static serializeBinaryToWriter(message: PublicKeyCredential, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKeyCredential;
  static deserializeBinaryFromReader(message: PublicKeyCredential, reader: jspb.BinaryReader): PublicKeyCredential;
}

export namespace PublicKeyCredential {
  export type AsObject = {
    id: string,
    type: string,
    clientDataJson: Uint8Array | string,
    attestationObject: Uint8Array | string,
  }
}

