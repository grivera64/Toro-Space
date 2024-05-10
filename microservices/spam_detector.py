import grpc
import logging
import grpc._server
import numpy as np
import tensorflow as tf
import traceback
import pandas as pd

from proto import spam_detector_pb2_grpc, spam_detector_pb2

from collections import defaultdict
from concurrent import futures
from grpc_reflection.v1alpha import reflection
# from typing import Any


class SpamDetectorServicer(spam_detector_pb2_grpc.SpamDetectorServicer):
    THRESHOLD = 0.8961

    def __init__(self, model_path):
        self.logger = logging.getLogger('SpamDetectorServicer')
        self.logger.setLevel(logging.DEBUG)

        self.logger.info('Loading Model...')
        self.model = tf.keras.models.load_model(model_path)
        self.logger.info('Model loaded successfully')

        self.logger.info('Preprocessing data...')
        data_df = pd.read_csv('emails.csv')
        self.words = list(data_df.columns[1:-1]) + ['#OOV#']

        self.word_to_index = defaultdict(None)
        for index, word in enumerate(self.words):
            self.word_to_index[word] = index
        self.logger.info('Preprocessing completed')

    def _prompt_to_np(self, prompt):
        word_tracker = defaultdict(int)
        for word in prompt.split(' '):
            if word in self.word_to_index:
                word_tracker[word] += 1
            else:
                word_tracker['#OOV#'] += 1

        res = np.zeros((len(self.words)),)

        for word, count in word_tracker.items():
            res[self.word_to_index[word]] = count
        return res

    def Scan(self, request, context):
        """Scans a post's content to verify if is ham or spam"""
        self.logger.info('Scan initiated')

        # Get content from request
        post_content = request.content
        try:
            self.logger.info(f'Passing in "{request.content}"')
            prompt_np = self._prompt_to_np(post_content)
            prompt_np = prompt_np.reshape((1, -1))
        except Exception as e:
            self.logger.error(f'Convert to np failed: {e}')
            return spam_detector_pb2.ScanResponse(result=spam_detector_pb2.ScanResponse.Result.UNKNOWN)
        
        # Ask Spam Model in tensorflow
        try:
            prediction = self.model.predict(prompt_np)
        except Exception as e:
            self.logger.error(f'Scan failed: {e}')
            return spam_detector_pb2.ScanResponse(result=spam_detector_pb2.ScanResponse.Result.UNKNOWN)

        print(f'isSpam confidence (treshold={SpamDetectorServicer.THRESHOLD}): {prediction}')
        if prediction < SpamDetectorServicer.THRESHOLD:
            result: spam_detector_pb2.ScanResponse.Result = spam_detector_pb2.ScanResponse.Result.HAM
        else:
            result: spam_detector_pb2.ScanResponse.Result = spam_detector_pb2.ScanResponse.Result.SPAM
        
        self.logger.info('Scan succeeded')
        return spam_detector_pb2.ScanResponse(result=result)


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
    server.add_insecure_port('127.0.0.1:3060')

    # Start server
    server.start()
    print('Server started at "127.0.0.1:3060"')
    server.wait_for_termination()
    pass

if __name__ == '__main__':
    main()
