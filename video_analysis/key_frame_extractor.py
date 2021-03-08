# -*- coding: utf-8 -*-
import cv2
import operator
import numpy as np
import os
# import matplotlib.pyplot as plt
import sys
from scipy.signal import argrelextrema


# 基于帧间差分的关键帧提取算法：
# 我们知道，将两帧图像进行差分，得到图像的平均像素强度可以用来衡量两帧图像的变化大小。因此，基于帧间差分的平均强度，每当视频中的某一帧与前一帧画面内容产生了大的变化，我们便认为它是关键帧，并将其提取出来。
# 这里一般有三种算法，基于差分强度的顺序，基于差分强度阈值，基于局部最大值。比较推荐使用第三种方法来提取视频的关键帧，我们选择具有平均帧间差分强度局部最大值的帧作为视频的关键帧，这种方法的提取结果在丰富度上表现更好一些，提取结果均匀分散在视频中。


# 对平均帧间差分强度时间序列进行平滑，有效的移除噪声来避免将相似场景下的若干帧均同时提取为关键帧
def smooth(x, window_len=13, window='hanning'):
    s = np.r_[2 * x[0] - x[window_len:1:-1],
              x, 2 * x[-1] - x[-1:-window_len:-1]]

    if window == 'flat':  # moving average
        w = np.ones(window_len, 'd')
    else:
        w = getattr(np, window)(window_len)
    y = np.convolve(w / w.sum(), s, mode='same')
    return y[window_len - 1:-window_len + 1]


class Frame:
    """class to hold information about each frame
    
    """

    def __init__(self, id, diff):
        self.id = id
        self.diff = diff

    def __lt__(self, other):
        if self.id == other.id:
            return self.id < other.id
        return self.id < other.id

    def __gt__(self, other):
        return other.__lt__(self)

    def __eq__(self, other):
        return self.id == other.id and self.id == other.id

    def __ne__(self, other):
        return not self.__eq__(other)


def rel_change(a, b):
    x = (b - a) / max(a, b)
    return x


def extract_keyframes(video_path: str, save_dir: str) -> [str]:
    # 视频文件名
    video_name = video_path.split('/')[-1].split('.')[0]

    # Setting fixed threshold criteria
    USE_THRESH = False
    # fixed threshold value
    THRESH = 0.6
    # Setting fixed threshold criteria
    USE_TOP_ORDER = False
    # Setting local maxima criteria
    USE_LOCAL_MAXIMA = True
    # Number of top sorted frames
    NUM_TOP_FRAMES = 50

    # smoothing window size
    len_window = int(50)

    # load video and compute diff between frames
    cap = cv2.VideoCapture(video_path)
    curr_frame = None
    prev_frame = None
    frame_diffs = []
    frames = []
    success, frame = cap.read()
    i = 0
    while success:
        luv = cv2.cvtColor(frame, cv2.COLOR_BGR2LUV)
        curr_frame = luv
        if curr_frame is not None and prev_frame is not None:
            # logic here
            diff = cv2.absdiff(curr_frame, prev_frame)
            diff_sum = np.sum(diff)
            diff_sum_mean = diff_sum / (diff.shape[0] * diff.shape[1])
            frame_diffs.append(diff_sum_mean)
            frame = Frame(i, diff_sum_mean)
            frames.append(frame)
        prev_frame = curr_frame
        i = i + 1
        success, frame = cap.read()
    cap.release()

    # compute keyframe
    keyframe_id_set = set()
    if USE_TOP_ORDER:
        # sort the list in descending order
        frames.sort(key=operator.attrgetter("diff"), reverse=True)
        for keyframe in frames[:NUM_TOP_FRAMES]:
            keyframe_id_set.add(keyframe.id)
    if USE_THRESH:
        for i in range(1, len(frames)):
            if rel_change(np.float(frames[i - 1].diff), np.float(frames[i].diff)) >= THRESH:
                keyframe_id_set.add(frames[i].id)
    if USE_LOCAL_MAXIMA:
        diff_array = np.array(frame_diffs)
        sm_diff_array = smooth(diff_array, len_window)
        frame_indexes = np.asarray(argrelextrema(sm_diff_array, np.greater))[0]
        for i in frame_indexes:
            keyframe_id_set.add(frames[i - 1].id)

        # 存帧间差分强度图，我们不需要就注释了
        # plt.figure(figsize=(40, 20))
        # plt.locator_params(numticks=100)
        # plt.stem(sm_diff_array)
        # plt.savefig(dir + 'plot.png')

    # save all keyframes as image
    cap = cv2.VideoCapture(str(video_path))
    curr_frame = None
    keyframes = []
    success, frame = cap.read()
    idx = 0
    count = 0

    ret = []
    while success:
        if idx in keyframe_id_set:
            count += 1
            name = "keyframe_" + str(count) + ".jpg"
            ret.append(name)
            cv2.imwrite(save_dir + name, frame)
            keyframe_id_set.remove(idx)
        idx = idx + 1
        success, frame = cap.read()
    cap.release()

    return ret


if __name__ == '__main__':
    print(extract_keyframes('output/1/1.mp4', 'output/'))
