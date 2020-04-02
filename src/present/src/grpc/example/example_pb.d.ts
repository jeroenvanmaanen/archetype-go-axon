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

export class Credentials extends jspb.Message {
  getApikey(): string;
  setApikey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Credentials.AsObject;
  static toObject(includeInstance: boolean, msg: Credentials): Credentials.AsObject;
  static serializeBinaryToWriter(message: Credentials, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Credentials;
  static deserializeBinaryFromReader(message: Credentials, reader: jspb.BinaryReader): Credentials;
}

export namespace Credentials {
  export type AsObject = {
    apikey: string,
  }
}

export class AccessToken extends jspb.Message {
  getJwt(): string;
  setJwt(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessToken.AsObject;
  static toObject(includeInstance: boolean, msg: AccessToken): AccessToken.AsObject;
  static serializeBinaryToWriter(message: AccessToken, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessToken;
  static deserializeBinaryFromReader(message: AccessToken, reader: jspb.BinaryReader): AccessToken;
}

export namespace AccessToken {
  export type AsObject = {
    jwt: string,
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

export class RecordCommand extends jspb.Message {
  getAggregateidentifier(): string;
  setAggregateidentifier(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RecordCommand.AsObject;
  static toObject(includeInstance: boolean, msg: RecordCommand): RecordCommand.AsObject;
  static serializeBinaryToWriter(message: RecordCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RecordCommand;
  static deserializeBinaryFromReader(message: RecordCommand, reader: jspb.BinaryReader): RecordCommand;
}

export namespace RecordCommand {
  export type AsObject = {
    aggregateidentifier: string,
  }
}

export class StopCommand extends jspb.Message {
  getAggregateidentifier(): string;
  setAggregateidentifier(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopCommand.AsObject;
  static toObject(includeInstance: boolean, msg: StopCommand): StopCommand.AsObject;
  static serializeBinaryToWriter(message: StopCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopCommand;
  static deserializeBinaryFromReader(message: StopCommand, reader: jspb.BinaryReader): StopCommand;
}

export namespace StopCommand {
  export type AsObject = {
    aggregateidentifier: string,
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

export class StartedRecordingEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartedRecordingEvent.AsObject;
  static toObject(includeInstance: boolean, msg: StartedRecordingEvent): StartedRecordingEvent.AsObject;
  static serializeBinaryToWriter(message: StartedRecordingEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartedRecordingEvent;
  static deserializeBinaryFromReader(message: StartedRecordingEvent, reader: jspb.BinaryReader): StartedRecordingEvent;
}

export namespace StartedRecordingEvent {
  export type AsObject = {
  }
}

export class StoppedRecordingEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StoppedRecordingEvent.AsObject;
  static toObject(includeInstance: boolean, msg: StoppedRecordingEvent): StoppedRecordingEvent.AsObject;
  static serializeBinaryToWriter(message: StoppedRecordingEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StoppedRecordingEvent;
  static deserializeBinaryFromReader(message: StoppedRecordingEvent, reader: jspb.BinaryReader): StoppedRecordingEvent;
}

export namespace StoppedRecordingEvent {
  export type AsObject = {
  }
}

export class SearchQuery extends jspb.Message {
  getQuery(): string;
  setQuery(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SearchQuery): SearchQuery.AsObject;
  static serializeBinaryToWriter(message: SearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchQuery;
  static deserializeBinaryFromReader(message: SearchQuery, reader: jspb.BinaryReader): SearchQuery;
}

export namespace SearchQuery {
  export type AsObject = {
    query: string,
  }
}

