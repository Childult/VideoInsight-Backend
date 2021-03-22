import logging

LOG_DIR = '/swc/log/'
MODULE_NAME = 'video_analysis'


def init_logger():
    logger = logging.getLogger('main')
    logger.setLevel(level=logging.INFO)

    formatter = logging.Formatter('%(levelname)s - %(asctime)s - %(name)s - %(message)s', datefmt='%m/%d %H:%M')

    # Handler
    levels = {
        'info': logging.INFO,
        'error': logging.ERROR
    }

    for level in levels.keys():
        handler = logging.FileHandler(f'{LOG_DIR}{MODULE_NAME}_{level}.log')
        handler.setLevel(levels.get(level))
        handler.setFormatter(formatter)
        logger.addHandler(handler)

    stream_handler = logging.StreamHandler()
    stream_handler.setFormatter(formatter)
    logger.addHandler(stream_handler)
