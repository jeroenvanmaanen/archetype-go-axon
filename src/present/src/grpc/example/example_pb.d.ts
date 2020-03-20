import * as jspb from "google-protobuf"

export class Greeting extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Greeting.AsObject;
  static toObject(includeInstance: boolean, msg: Greeting): Greeting.AsObject;
  static serializeBinaryToWriter(message: Greeting, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Greeting;
  static deserializeBinaryFromReader(message: Greeting, reader: jspb.BinaryReader): Greeting;
}

export namespace Greeting {
  export type AsObject = {
    message: string,
  }
}

export class Acknowledgement extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Acknowledgement.AsObject;
  static toObject(includeInstance: boolean, msg: Acknowledgement): Acknowledgement.AsObject;
  static serializeBinaryToWriter(message: Acknowledgement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Acknowledgement;
  static deserializeBinaryFromReader(message: Acknowledgement, reader: jspb.BinaryReader): Acknowledgement;
}

export namespace Acknowledgement {
  export type AsObject = {
    message: string,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

export class GreetCommand extends jspb.Message {
  getAggregateidentifier(): string;
  setAggregateidentifier(value: string): void;

  getMessage(): Greeting | undefined;
  setMessage(value?: Greeting): void;
  hasMessage(): boolean;
  clearMessage(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GreetCommand.AsObject;
  static toObject(includeInstance: boolean, msg: GreetCommand): GreetCommand.AsObject;
  static serializeBinaryToWriter(message: GreetCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GreetCommand;
  static deserializeBinaryFromReader(message: GreetCommand, reader: jspb.BinaryReader): GreetCommand;
}

export namespace GreetCommand {
  export type AsObject = {
    aggregateidentifier: string,
    message?: Greeting.AsObject,
  }
}

export class GreetedEvent extends jspb.Message {
  getMessage(): Greeting | undefined;
  setMessage(value?: Greeting): void;
  hasMessage(): boolean;
  clearMessage(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GreetedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: GreetedEvent): GreetedEvent.AsObject;
  static serializeBinaryToWriter(message: GreetedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GreetedEvent;
  static deserializeBinaryFromReader(message: GreetedEvent, reader: jspb.BinaryReader): GreetedEvent;
}

export namespace GreetedEvent {
  export type AsObject = {
    message?: Greeting.AsObject,
  }
}

