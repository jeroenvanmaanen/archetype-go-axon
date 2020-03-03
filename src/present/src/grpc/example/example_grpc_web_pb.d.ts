import * as grpcWeb from 'grpc-web';

import {
  Acknowledgement,
  Greeting} from './example_pb';

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

}

export class GreeterServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  greet(
    request: Greeting,
    metadata?: grpcWeb.Metadata
  ): Promise<Acknowledgement>;

}

