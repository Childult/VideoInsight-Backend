#!/usr/bin/python
# -*- coding:utf-8 -*-
from __future__ import print_function

import logging
import os
import os.path as osp
import argparse
import sys
import h5py
from tqdm import tqdm
import torch
import torch.nn as nn
import torch.backends.cudnn as cudnn

from networkss.DSN import *
from utils import vsum_tool
from utils.generate_dataset import Generate_Dataset
import cv2

logger = logging.getLogger('main.vs')
logger.setLevel(level=logging.INFO)

parser = argparse.ArgumentParser("Pytorch code for unsupervised video summarization with REINFORCE")
# Dataset options
parser.add_argument('-i', '--input', type=str, default='', help="input video")
parser.add_argument('-o', '--output', type=str, default='output/', help="output video")
# Misc
parser.add_argument('--seed', type=int, default=1, help="random seed (default: 1)")
parser.add_argument('--gpu', type=str, default='0', help="which gpu devices to use")
# Model options
parser.add_argument('--input-dim', type=int, default=2048, help="input dimension (default: 1024)")
parser.add_argument('--hidden-dim', type=int, default=256, help="hidden unit dimension of DSN (default: 256)")
parser.add_argument('--num-layers', type=int, default=1, help="number of RNN layers (default: 1)")
parser.add_argument('--rnn-cell', type=str, default='lstm', help="RNN cell type (default: lstm)")

parser.add_argument('-d', '--dataset', type=str, help="path to h5 dataset (required)")

parser.add_argument('--model', type=str, default='/swc/code/video_analysis/model/best_model_epoch60.pth.tar',
                    help="path to model file")
parser.add_argument('--save-dir', type=str, default='output/', help="path to save output (default: 'output/')")
parser.add_argument('--use-cpu', action='store_true', help="use cpu device")

parser.add_argument('--save-name', default='', help="'generate video '")
parser.add_argument('--fps', type=int, default=30, help="frames per second")
parser.add_argument('--width', type=int, default=640, help="frame width")
parser.add_argument('--height', type=int, default=480, help="frame height")

parser.add_argument('--train-data', action='store_true', help="")

args = parser.parse_args()

torch.manual_seed(args.seed)
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
os.environ['CUDA_VISIBLE_DEVICES'] = args.gpu
use_gpu = torch.cuda.is_available()
if args.use_cpu:
    use_gpu = False


def main(video_path=''):
    print("==========\nArgs:{}\n==========".format(args))

    if use_gpu:
        print("Currently using GPU {}".format(args.gpu))
        cudnn.benchmark = True
        torch.cuda.manual_seed_all(args.seed)
    else:
        print("Currently using CPU")

    print("Initialize dataset {}".format(args.dataset))
    dataset = h5py.File(args.dataset, 'r')
    test_keys = []

    for key in dataset.keys():
        test_keys.append(key)

    print("Load model")
    model = DSN(in_dim=args.input_dim, hid_dim=args.hidden_dim, num_layers=args.num_layers, cell=args.rnn_cell)
    print("Model size: {:.5f}M".format(sum(p.numel() for p in model.parameters()) / 1000000.0))

    if args.model:
        print("Loading checkpoint from '{}'".format(args.model))
        checkpoint = torch.load(args.model, map_location=device)
        model.load_state_dict(checkpoint)

    if use_gpu:
        model = nn.DataParallel(model).cuda()
    evaluate(model, dataset, test_keys)
    logger.info('Begin to do video summary')
    if video_path != '':
        video2summary(os.path.join(args.save_dir, 'result.h5'), video_path, args.save_dir)
    else:
        video2summary(os.path.join(args.save_dir, 'result.h5'), args.input, args.save_dir)


def evaluate(model, dataset, test_keys):
    with torch.no_grad():
        model.eval()

        if not os.path.isdir(args.save_dir):
            os.mkdir(args.save_dir)

        h5_res = h5py.File(os.path.join(args.save_dir, 'result.h5'), 'w')
        for key in dataset.keys():
            print(dataset[key].name)
            print(dataset[key])
        for key_idx, key in enumerate(test_keys):
            seq = dataset[key]['features'][...]
            seq = torch.from_numpy(seq).unsqueeze(0)
            if use_gpu:
                seq = seq.cuda()
            probs = model(seq)
            probs = probs.data.cpu().squeeze().numpy()

            cps = dataset[key]['change_points'][...]
            num_frames = dataset[key]['n_frames'][()]
            nfps = dataset[key]['n_frame_per_seg'][...].tolist()
            positions = dataset[key]['picks'][...]
            video_name = dataset[key]['video_name'][()]
            fps = dataset[key]['fps'][()]

            sum = 0
            for i in range(len(nfps)):
                sum += nfps[i]

            machine_summary = vsum_tool.generate_summary(probs, cps, num_frames, nfps, positions)
            h5_res.create_dataset(key + '/score', data=probs)
            h5_res.create_dataset(key + '/machine_summary', data=machine_summary)
            h5_res.create_dataset(key + '/video_name', data=video_name)
            h5_res.create_dataset(key + '/fps', data=fps)

    h5_res.close()


def frm2video(video_dir, summary, vid_writer):
    print('[INFO] Video Summary')
    video_capture = cv2.VideoCapture(video_dir)
    count = 0
    for idx, val in tqdm(enumerate(summary)):
        ret, frame = video_capture.read()
        if val == 1 and ret:
            frm = cv2.resize(frame, (args.width, args.height))
            vid_writer.write(frm)
        else:
            count += 1
    print('[OUTPUT] total {} frame, ignore {} frame'.format(len(summary) - count, count))


def video2summary(h5_dir, video_dir, output_dir):
    if not osp.exists(output_dir):
        os.mkdir(output_dir)

    h5_res = h5py.File(h5_dir, 'r')
    print(len(list(h5_res.keys())))

    for idx1 in range(len(list(h5_res.keys()))):
        key = list(h5_res.keys())[idx1]
        summary = h5_res[key]['machine_summary'][...]
        print(h5_res[str.encode(key)][str.encode('video_name')][()])
        video_name = bytes.decode(h5_res[key]['video_name'][()]).split('/')[-1]
        fps = h5_res[key]['fps'][()]
        if not os.path.isdir(osp.join(output_dir, video_name)):
            os.mkdir(osp.join(output_dir, video_name))
        vid_writer = cv2.VideoWriter(
            osp.join(output_dir, video_name, args.save_name),
            cv2.VideoWriter_fourcc('M', 'P', '4', 'V'),
            fps,
            (args.width, args.height),
        )
        frm2video(video_dir, summary, vid_writer)
        vid_writer.release()
    h5_res.close()


def video_summarize_api(video_path, save_dir='/swc/resource/compressed/'):
    video_name = video_path.split('/')[-1].split('.')[0]
    args.dataset = os.path.join(save_dir, video_name + '.h5')
    args.save_name = os.path.join(save_dir, video_name + '-compressed.mp4')
    args.save_dir = save_dir
    logger.info('Begin to generate dataset')
    gen = Generate_Dataset(video_path, args.dataset)
    gen.generate_dataset()
    gen.h5_file.close()
    logger.info('Done: generate dataset')
    main(video_path)
    return args.save_name


if __name__ == '__main__':
    print(video_summarize_api("/swc/code/video_analysis/dataset/test.mp4"))
