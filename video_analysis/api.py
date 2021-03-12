import video_summary as vsumm
import key_frame_extractor as keyframe_extractor


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
    try:
        # 对视频进行摘要，得到浓缩版视频
        compressed_video = video_summarize(file, save_dir)
    except Exception as e:
        ret['Error'] = 'failed to summarize the video: ' + str(e)
    else:
        # 若无异常则进行下一步操作
        try:
            ret['VAbstract'] = extract_key_frame(compressed_video, save_dir)
        except Exception as e:
            ret['Error'] = 'failed to extract key frames: ' + str(e)

    return ret


def extract_key_frame(file: str, save_dir: str) -> [str]:
    """
    提取视频中的关键帧
    :param file: 视频所在路径
    :param save_dir: 关键帧图片保存目录
    :return: 关键帧图片名称列表（不包含目录名），如 ['keyframe_1.jpg', 'keyframe_2.jpg']
    """
    return keyframe_extractor.extract_keyframes(file, save_dir)


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


def video_summarize(file: str, save_dir: str) -> str:
    """
    视频摘要
    :param file: 视频文件路径 dataset/taichi.mp4
    :param save_dir:
    :return: 生成的视频路径 output/taichi/taichi.mp4
    """
    return vsumm.video_summarize_api(file, save_dir)

# if __name__ == '__main__':
#     print(generate_abstract_from_video('/swc/code/video_analysis/dataset/3.mp4', '/swc/resource/'))
