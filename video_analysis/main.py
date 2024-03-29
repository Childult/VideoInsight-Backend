import asyncio
import logging

import grpc

from api import generate_abstract_from_video
from utils.logger import init_logger
from video_pb2 import VideoInfo, Result
from video_pb2_grpc import VideoAnalysisServicer, add_VideoAnalysisServicer_to_server


logger = logging.getLogger('main.grpc')
logger.setLevel(level=logging.INFO)


class VideoService(VideoAnalysisServicer):
    async def GetStaticVideoAbstract(self, request: VideoInfo, context: grpc.aio.ServicerContext) -> Result:
        ret = generate_abstract_from_video(request.file, request.save_dir)
        return Result(job_id=request.job_id, error=ret['Error'], pic_name=ret['VAbstract'])


async def serve() -> None:
    server = grpc.aio.server()
    add_VideoAnalysisServicer_to_server(VideoService(), server)
    listen_addr = '[::]:50051'
    server.add_insecure_port(listen_addr)
    logger.info("Starting server on %s", listen_addr)
    await server.start()
    try:
        await server.wait_for_termination()
    except KeyboardInterrupt:
        # Shuts down the server with 0 seconds of grace period. During the
        # grace period, the server won't accept new connections and allow
        # existing RPCs to continue within the grace period.
        await server.stop(0)


if __name__ == '__main__':
    init_logger()
    asyncio.run(serve())
