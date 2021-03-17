import logging

import cv2
from func_timeout import func_timeout, FunctionTimedOut

import video_summary as vsumm
import key_frame_extractor as keyframe_extractor

logger = logging.getLogger('main.api')
logger.setLevel(level=logging.INFO)

# 处理时间不应超过视频时间（分钟）的倍数
TIME_OUT_RATIO = 6


def generate_abstract_from_video(file: str, save_dir: str) -> dict:
    """
    根据视频文件，生成静态视频摘要（图片）。
    :param file: 视频文件路径
    :param save_dir: 静态视频摘要（图片）的保存目录，需确保该目录已经存在
    :return: 静态视频摘要（图片）文件名列表（不含目录名）
    """
    ret = {
        "VAbstract": [],  # 静态视频摘要文件名列表
        "Error": ""  # 错误信息
    }
    video_min = get_video_duration(file)
    try:
        # 对视频进行摘要，得到浓缩版视频
        logger.info('Begin to do video summarize: %s, %s', file, save_dir)
        compressed_video = func_timeout(video_min * TIME_OUT_RATIO * 60, video_summarize, args=(file, save_dir))
        logger.info('Done: video summarize: %s, %s', file, save_dir)
    except FunctionTimedOut:
        logger.error("summarize the video time out!")
        ret['Error'] = '视频摘要超时'
    except Exception as e:
        logger.error(e, exc_info=True)
        ret['Error'] = 'failed to summarize the video: ' + str(e)
    else:
        # 若无异常则进行下一步操作
        try:
            logger.info('Begin to extract key frames: %s, %s', file, save_dir)
            ret['VAbstract'] = extract_key_frame(compressed_video, save_dir)
            logger.info('Done: extract key frames: %s, %s', file, save_dir)
        except Exception as e:
            logger.error(e, exc_info=True)
            ret['Error'] = 'failed to extract key frames: ' + str(e)

    return ret


def get_video_duration(file: str) -> int:
    """
    获取视频文件时长，以分钟为单位
    :param file: 视频文件
    :return: 视频时长，单位：分钟
    """
    video = cv2.VideoCapture(file)
    fps = int(round(video.get(cv2.CAP_PROP_FPS)))  # 帧率
    frame_counter = int(video.get(cv2.CAP_PROP_FRAME_COUNT))  # 总帧数
    video.release()
    return round(frame_counter / fps / 60)


def extract_key_frame(file: str, save_dir: str) -> [str]:
    """
    提取视频中的关键帧
    :param file: 视频所在路径
    :param save_dir: 关键帧图片保存目录
    :return: 关键帧图片名称列表（不包含目录名），如 ['keyframe_1.jpg', 'keyframe_2.jpg']
    """
    return keyframe_extractor.extract_keyframes(file, save_dir)


def video_summarize(file: str, save_dir: str) -> str:
    """
    视频摘要
    :param file: 视频文件路径 dataset/taichi.mp4
    :param save_dir:
    :return: 生成的视频路径 output/taichi/taichi.mp4
    """
    return vsumm.video_summarize_api(file, save_dir)


if __name__ == '__main__':
    print(generate_abstract_from_video('/swc/code/video_analysis/dataset/test.mp4', '/swc/resource/'))
