import asyncio

import grpc

from video_analysis.video_pb2 import VideoInfo
from video_analysis.video_pb2_grpc import VideoAnalysisStub


async def run() -> None:
    async with grpc.aio.insecure_channel('192.168.2.80:50051') as channel:
        stub = VideoAnalysisStub(channel)
        response = await stub.GetStaticVideoAbstract(
            VideoInfo(job_id='1', file='/swc/resource/1617010399/MTYxNzAxMDQxMC4yNDM2OTJodHRwczovL3d3dy5iaWxpYmlsaS5jb20vdmlkZW8vQlYxM1o0eTFHN1Ax.mp4', save_dir='/swc/resource/tests/'))
    print(str(response))


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    result = loop.run_until_complete(run())
