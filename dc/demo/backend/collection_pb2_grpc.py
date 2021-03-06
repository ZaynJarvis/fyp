# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import collection_pb2 as collection__pb2


class CloudPullServiceStub(object):
    """RPC-wise, both pull and push model are not for high concurrency scenario due to the constraint on file handlers and
    HTTP connection. (since both require long connection, otherwise there will be delay issue).
    However, the combination of Cloud push config + agent push event model is okay, but still not suitable for high
    concurrency.
    In the case of many agents, like city security camera, using light-weight message queue like mqtt will be better.

    use pull model by assumption that config changes will be less frequent than collection notification
    here, open an HTTP2 long connection will be more efficient than client setting up short term TCP connections.
    pull model allow cloud to enable active health check, while need to setup registry for service discovery.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.SendConfig = channel.unary_unary(
                '/CloudPullService/SendConfig',
                request_serializer=collection__pb2.CollectionConfig.SerializeToString,
                response_deserializer=collection__pb2.Result.FromString,
                )
        self.HealthCheck = channel.unary_unary(
                '/CloudPullService/HealthCheck',
                request_serializer=collection__pb2.CloudInfo.SerializeToString,
                response_deserializer=collection__pb2.Result.FromString,
                )
        self.Listen = channel.unary_stream(
                '/CloudPullService/Listen',
                request_serializer=collection__pb2.CloudInfo.SerializeToString,
                response_deserializer=collection__pb2.CollectionEvent.FromString,
                )


class CloudPullServiceServicer(object):
    """RPC-wise, both pull and push model are not for high concurrency scenario due to the constraint on file handlers and
    HTTP connection. (since both require long connection, otherwise there will be delay issue).
    However, the combination of Cloud push config + agent push event model is okay, but still not suitable for high
    concurrency.
    In the case of many agents, like city security camera, using light-weight message queue like mqtt will be better.

    use pull model by assumption that config changes will be less frequent than collection notification
    here, open an HTTP2 long connection will be more efficient than client setting up short term TCP connections.
    pull model allow cloud to enable active health check, while need to setup registry for service discovery.
    """

    def SendConfig(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def HealthCheck(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Listen(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_CloudPullServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'SendConfig': grpc.unary_unary_rpc_method_handler(
                    servicer.SendConfig,
                    request_deserializer=collection__pb2.CollectionConfig.FromString,
                    response_serializer=collection__pb2.Result.SerializeToString,
            ),
            'HealthCheck': grpc.unary_unary_rpc_method_handler(
                    servicer.HealthCheck,
                    request_deserializer=collection__pb2.CloudInfo.FromString,
                    response_serializer=collection__pb2.Result.SerializeToString,
            ),
            'Listen': grpc.unary_stream_rpc_method_handler(
                    servicer.Listen,
                    request_deserializer=collection__pb2.CloudInfo.FromString,
                    response_serializer=collection__pb2.CollectionEvent.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'CloudPullService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class CloudPullService(object):
    """RPC-wise, both pull and push model are not for high concurrency scenario due to the constraint on file handlers and
    HTTP connection. (since both require long connection, otherwise there will be delay issue).
    However, the combination of Cloud push config + agent push event model is okay, but still not suitable for high
    concurrency.
    In the case of many agents, like city security camera, using light-weight message queue like mqtt will be better.

    use pull model by assumption that config changes will be less frequent than collection notification
    here, open an HTTP2 long connection will be more efficient than client setting up short term TCP connections.
    pull model allow cloud to enable active health check, while need to setup registry for service discovery.
    """

    @staticmethod
    def SendConfig(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/CloudPullService/SendConfig',
            collection__pb2.CollectionConfig.SerializeToString,
            collection__pb2.Result.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def HealthCheck(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/CloudPullService/HealthCheck',
            collection__pb2.CloudInfo.SerializeToString,
            collection__pb2.Result.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Listen(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(request, target, '/CloudPullService/Listen',
            collection__pb2.CloudInfo.SerializeToString,
            collection__pb2.CollectionEvent.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)


class AgentPushServiceStub(object):
    """use push model so registry for sd is not needed.
    but push model introduce a tradeoff between a unnecessary long connection (ListenConfig) and delay on polling
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.ListenConfig = channel.unary_stream(
                '/AgentPushService/ListenConfig',
                request_serializer=collection__pb2.AgentInfo.SerializeToString,
                response_deserializer=collection__pb2.CollectionConfig.FromString,
                )
        self.GetConfig = channel.unary_unary(
                '/AgentPushService/GetConfig',
                request_serializer=collection__pb2.AgentInfo.SerializeToString,
                response_deserializer=collection__pb2.CollectionConfig.FromString,
                )
        self.SendNotification = channel.stream_unary(
                '/AgentPushService/SendNotification',
                request_serializer=collection__pb2.CollectionEvent.SerializeToString,
                response_deserializer=collection__pb2.Result.FromString,
                )


class AgentPushServiceServicer(object):
    """use push model so registry for sd is not needed.
    but push model introduce a tradeoff between a unnecessary long connection (ListenConfig) and delay on polling
    """

    def ListenConfig(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetConfig(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SendNotification(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AgentPushServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'ListenConfig': grpc.unary_stream_rpc_method_handler(
                    servicer.ListenConfig,
                    request_deserializer=collection__pb2.AgentInfo.FromString,
                    response_serializer=collection__pb2.CollectionConfig.SerializeToString,
            ),
            'GetConfig': grpc.unary_unary_rpc_method_handler(
                    servicer.GetConfig,
                    request_deserializer=collection__pb2.AgentInfo.FromString,
                    response_serializer=collection__pb2.CollectionConfig.SerializeToString,
            ),
            'SendNotification': grpc.stream_unary_rpc_method_handler(
                    servicer.SendNotification,
                    request_deserializer=collection__pb2.CollectionEvent.FromString,
                    response_serializer=collection__pb2.Result.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'AgentPushService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class AgentPushService(object):
    """use push model so registry for sd is not needed.
    but push model introduce a tradeoff between a unnecessary long connection (ListenConfig) and delay on polling
    """

    @staticmethod
    def ListenConfig(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(request, target, '/AgentPushService/ListenConfig',
            collection__pb2.AgentInfo.SerializeToString,
            collection__pb2.CollectionConfig.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetConfig(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/AgentPushService/GetConfig',
            collection__pb2.AgentInfo.SerializeToString,
            collection__pb2.CollectionConfig.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SendNotification(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_unary(request_iterator, target, '/AgentPushService/SendNotification',
            collection__pb2.CollectionEvent.SerializeToString,
            collection__pb2.Result.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)


class LocalStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Image = channel.stream_unary(
                '/Local/Image',
                request_serializer=collection__pb2.ImageReport.SerializeToString,
                response_deserializer=collection__pb2.Result.FromString,
                )


class LocalServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Image(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_LocalServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Image': grpc.stream_unary_rpc_method_handler(
                    servicer.Image,
                    request_deserializer=collection__pb2.ImageReport.FromString,
                    response_serializer=collection__pb2.Result.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'Local', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Local(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Image(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_unary(request_iterator, target, '/Local/Image',
            collection__pb2.ImageReport.SerializeToString,
            collection__pb2.Result.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
