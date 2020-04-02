import * as grpcWeb from 'grpc-web';

import {
  AccessToken,
  Acknowledgement,
  Credentials,
  Empty,
  Greeting,
  SearchQuery} from './example_pb';

export class GreeterServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  greet(
    request: Greeting,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Acknowledgement) => void
  ): grpcWeb.ClientReadableStream<Acknowledgement>;

  greetings(
    request: Empty,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<Greeting>;

  record(
    request: Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Empty) => void
  ): grpcWeb.ClientReadableStream<Empty>;

  stop(
    request: Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Empty) => void
  ): grpcWeb.ClientReadableStream<Empty>;

  search(
    request: SearchQuery,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<Greeting>;

  authorize(
    request: Credentials,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AccessToken) => void
  ): grpcWeb.ClientReadableStream<AccessToken>;

}

export class GreeterServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  greet(
    request: Greeting,
    metadata?: grpcWeb.Metadata
  ): Promise<Acknowledgement>;

  greetings(
    request: Empty,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<Greeting>;

  record(
    request: Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<Empty>;

  stop(
    request: Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<Empty>;

  search(
    request: SearchQuery,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<Greeting>;

  authorize(
    request: Credentials,
    metadata?: grpcWeb.Metadata
  ): Promise<AccessToken>;

}

