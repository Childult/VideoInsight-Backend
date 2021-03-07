import video_summary as vsumm
import key_frame_extractor as keyframe_extractor

def extract_key_frame(file: str) -> [str]:
    """
    提取视频中的关键帧
    :param file: 视频所在路径 dataset/taichi.mp4
    :return: 关键帧图片路径 keyframes/taichi (taichi是视频文件名，关键帧在该路径下，如keyframes/taichi/keyframe_1.jpg)
    """
    return keyframe_extractor.extract_keyframes(file)


def extract_audio(file: str) -> str:
    """
    提取视频文件中的音频
    :param file: 视频所在路径
    :return: 音频文件名
    """
    return 'Not implemented.'


def scene_text_spotting(file: str) -> [str]:
    """
    场景文本识别
    :param file:
    :return: 识别的文本列表
    """
    return ['Not implemented.']


def video_summarize(file: str) -> str:
    """
    视频摘要
    :param file: 视频文件路径 dataset/taichi.mp4
    :return: 生成的视频路径 output/taichi/taichi.mp4
    """
    return vsumm.video_summarize_api(file)
