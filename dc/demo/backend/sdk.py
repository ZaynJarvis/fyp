import grpc
import json
import queue
import threading
import collection_pb2
import collection_pb2_grpc

class SDK():
    def gen(self):
        for item in iter(self.q.get, None):
            yield item

    def run(self):
        self.stub.Image(self.gen())

    def __init__(self, addr):
        self.q = queue.Queue()
        ch = grpc.insecure_channel(addr)
        self.stub = collection_pb2_grpc.LocalStub(ch)
        download_thread = threading.Thread(target=self.run, name="runner", args=[])
        download_thread.start()

    def Image(self, id, image, result):
        img = collection_pb2.ImageReport(
            id=id, 
            img=image,
            result=json.dumps(result).encode('utf-8'))
        self.q.put(img)

if __name__ == '__main__':
    sdk = SDK("host.docker.internal:7000")
    sdk.Image("id", None,{})
