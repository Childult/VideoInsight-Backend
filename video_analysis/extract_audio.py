import os


def extract_audio(file):
    title = file.split('.', -1)[0]
    try:
        os.system('ffmpeg -i ' + file + ' -f mp3 -ar 16000 ' + title + '.mp3')
    except Exception as e:
        return "Error: {0}".format(str(e))
    return title + '.aac'
