/* eslint-disable */
/**
 * @fileoverview gRPC-Web generated client stub for org.leialearns.grpc.example
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.org = {};
proto.org.leialearns = {};
proto.org.leialearns.grpc = {};
proto.org.leialearns.grpc.example = require('./example_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.org.leialearns.grpc.example.GreeterServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.org.leialearns.grpc.example.GreeterServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.org.leialearns.grpc.example.Greeting,
 *   !proto.org.leialearns.grpc.example.Acknowledgement>}
 */
const methodDescriptor_GreeterService_Greet = new grpc.web.MethodDescriptor(
  '/org.leialearns.grpc.example.GreeterService/Greet',
  grpc.web.MethodType.UNARY,
  proto.org.leialearns.grpc.example.Greeting,
  proto.org.leialearns.grpc.example.Acknowledgement,
  /**
   * @param {!proto.org.leialearns.grpc.example.Greeting} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Acknowledgement.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.org.leialearns.grpc.example.Greeting,
 *   !proto.org.leialearns.grpc.example.Acknowledgement>}
 */
const methodInfo_GreeterService_Greet = new grpc.web.AbstractClientBase.MethodInfo(
  proto.org.leialearns.grpc.example.Acknowledgement,
  /**
   * @param {!proto.org.leialearns.grpc.example.Greeting} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Acknowledgement.deserializeBinary
);


/**
 * @param {!proto.org.leialearns.grpc.example.Greeting} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.org.leialearns.grpc.example.Acknowledgement)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.org.leialearns.grpc.example.Acknowledgement>|undefined}
 *     The XHR Node Readable Stream
 */
proto.org.leialearns.grpc.example.GreeterServiceClient.prototype.greet =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Greet',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Greet,
      callback);
};


/**
 * @param {!proto.org.leialearns.grpc.example.Greeting} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.org.leialearns.grpc.example.Acknowledgement>}
 *     A native promise that resolves to the response
 */
proto.org.leialearns.grpc.example.GreeterServicePromiseClient.prototype.greet =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Greet',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Greet);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.org.leialearns.grpc.example.Empty,
 *   !proto.org.leialearns.grpc.example.Greeting>}
 */
const methodDescriptor_GreeterService_Greetings = new grpc.web.MethodDescriptor(
  '/org.leialearns.grpc.example.GreeterService/Greetings',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.org.leialearns.grpc.example.Empty,
  proto.org.leialearns.grpc.example.Greeting,
  /**
   * @param {!proto.org.leialearns.grpc.example.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Greeting.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.org.leialearns.grpc.example.Empty,
 *   !proto.org.leialearns.grpc.example.Greeting>}
 */
const methodInfo_GreeterService_Greetings = new grpc.web.AbstractClientBase.MethodInfo(
  proto.org.leialearns.grpc.example.Greeting,
  /**
   * @param {!proto.org.leialearns.grpc.example.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Greeting.deserializeBinary
);


/**
 * @param {!proto.org.leialearns.grpc.example.Empty} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.org.leialearns.grpc.example.Greeting>}
 *     The XHR Node Readable Stream
 */
proto.org.leialearns.grpc.example.GreeterServiceClient.prototype.greetings =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Greetings',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Greetings);
};


/**
 * @param {!proto.org.leialearns.grpc.example.Empty} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.org.leialearns.grpc.example.Greeting>}
 *     The XHR Node Readable Stream
 */
proto.org.leialearns.grpc.example.GreeterServicePromiseClient.prototype.greetings =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Greetings',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Greetings);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.org.leialearns.grpc.example.Empty,
 *   !proto.org.leialearns.grpc.example.Empty>}
 */
const methodDescriptor_GreeterService_Record = new grpc.web.MethodDescriptor(
  '/org.leialearns.grpc.example.GreeterService/Record',
  grpc.web.MethodType.UNARY,
  proto.org.leialearns.grpc.example.Empty,
  proto.org.leialearns.grpc.example.Empty,
  /**
   * @param {!proto.org.leialearns.grpc.example.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.org.leialearns.grpc.example.Empty,
 *   !proto.org.leialearns.grpc.example.Empty>}
 */
const methodInfo_GreeterService_Record = new grpc.web.AbstractClientBase.MethodInfo(
  proto.org.leialearns.grpc.example.Empty,
  /**
   * @param {!proto.org.leialearns.grpc.example.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Empty.deserializeBinary
);


/**
 * @param {!proto.org.leialearns.grpc.example.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.org.leialearns.grpc.example.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.org.leialearns.grpc.example.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.org.leialearns.grpc.example.GreeterServiceClient.prototype.record =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Record',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Record,
      callback);
};


/**
 * @param {!proto.org.leialearns.grpc.example.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.org.leialearns.grpc.example.Empty>}
 *     A native promise that resolves to the response
 */
proto.org.leialearns.grpc.example.GreeterServicePromiseClient.prototype.record =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Record',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Record);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.org.leialearns.grpc.example.Empty,
 *   !proto.org.leialearns.grpc.example.Empty>}
 */
const methodDescriptor_GreeterService_Stop = new grpc.web.MethodDescriptor(
  '/org.leialearns.grpc.example.GreeterService/Stop',
  grpc.web.MethodType.UNARY,
  proto.org.leialearns.grpc.example.Empty,
  proto.org.leialearns.grpc.example.Empty,
  /**
   * @param {!proto.org.leialearns.grpc.example.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.org.leialearns.grpc.example.Empty,
 *   !proto.org.leialearns.grpc.example.Empty>}
 */
const methodInfo_GreeterService_Stop = new grpc.web.AbstractClientBase.MethodInfo(
  proto.org.leialearns.grpc.example.Empty,
  /**
   * @param {!proto.org.leialearns.grpc.example.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Empty.deserializeBinary
);


/**
 * @param {!proto.org.leialearns.grpc.example.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.org.leialearns.grpc.example.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.org.leialearns.grpc.example.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.org.leialearns.grpc.example.GreeterServiceClient.prototype.stop =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Stop',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Stop,
      callback);
};


/**
 * @param {!proto.org.leialearns.grpc.example.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.org.leialearns.grpc.example.Empty>}
 *     A native promise that resolves to the response
 */
proto.org.leialearns.grpc.example.GreeterServicePromiseClient.prototype.stop =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Stop',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Stop);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.org.leialearns.grpc.example.SearchQuery,
 *   !proto.org.leialearns.grpc.example.Greeting>}
 */
const methodDescriptor_GreeterService_Search = new grpc.web.MethodDescriptor(
  '/org.leialearns.grpc.example.GreeterService/Search',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.org.leialearns.grpc.example.SearchQuery,
  proto.org.leialearns.grpc.example.Greeting,
  /**
   * @param {!proto.org.leialearns.grpc.example.SearchQuery} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Greeting.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.org.leialearns.grpc.example.SearchQuery,
 *   !proto.org.leialearns.grpc.example.Greeting>}
 */
const methodInfo_GreeterService_Search = new grpc.web.AbstractClientBase.MethodInfo(
  proto.org.leialearns.grpc.example.Greeting,
  /**
   * @param {!proto.org.leialearns.grpc.example.SearchQuery} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.org.leialearns.grpc.example.Greeting.deserializeBinary
);


/**
 * @param {!proto.org.leialearns.grpc.example.SearchQuery} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.org.leialearns.grpc.example.Greeting>}
 *     The XHR Node Readable Stream
 */
proto.org.leialearns.grpc.example.GreeterServiceClient.prototype.search =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Search',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Search);
};


/**
 * @param {!proto.org.leialearns.grpc.example.SearchQuery} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.org.leialearns.grpc.example.Greeting>}
 *     The XHR Node Readable Stream
 */
proto.org.leialearns.grpc.example.GreeterServicePromiseClient.prototype.search =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/org.leialearns.grpc.example.GreeterService/Search',
      request,
      metadata || {},
      methodDescriptor_GreeterService_Search);
};


module.exports = proto.org.leialearns.grpc.example;

