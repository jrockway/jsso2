import * as jspb from 'google-protobuf'



export class TransformNameRequest extends jspb.Message {
  getName(): string;
  setName(value: string): TransformNameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransformNameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TransformNameRequest): TransformNameRequest.AsObject;
  static serializeBinaryToWriter(message: TransformNameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransformNameRequest;
  static deserializeBinaryFromReader(message: TransformNameRequest, reader: jspb.BinaryReader): TransformNameRequest;
}

export namespace TransformNameRequest {
  export type AsObject = {
    name: string,
  }
}

export class TransformNameReply extends jspb.Message {
  getResult(): string;
  setResult(value: string): TransformNameReply;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TransformNameReply.AsObject;
  static toObject(includeInstance: boolean, msg: TransformNameReply): TransformNameReply.AsObject;
  static serializeBinaryToWriter(message: TransformNameReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TransformNameReply;
  static deserializeBinaryFromReader(message: TransformNameReply, reader: jspb.BinaryReader): TransformNameReply;
}

export namespace TransformNameReply {
  export type AsObject = {
    result: string,
  }
}

