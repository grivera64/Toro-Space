import grpc

from concurrent import futures
from grpc_reflection.v1alpha import reflection
from proto import spam_detector_pb2_grpc, spam_detector_pb2
from spam_detector_servicer import SpamDetectorServicer

def main():
    print('Creating SpamDetectorServicer...')
    model_path = 'spam_detector_model'
    servicer = SpamDetectorServicer(model_path)
    print('Creating SpamDetectorServicer...')

    # Add Service and grpc Reflection
    SERVICE_NAMES = (
        spam_detector_pb2.DESCRIPTOR.services_by_name['SpamDetector'].full_name,
        reflection.SERVICE_NAME
    )
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    spam_detector_pb2_grpc.add_SpamDetectorServicer_to_server(servicer, server)
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    server.add_insecure_port('0.0.0.0:3060')

    # Start server
    server.start()
    print('Server started at "0.0.0.0:3060"')
    server.wait_for_termination()
    pass

if __name__ == '__main__':
    main()
