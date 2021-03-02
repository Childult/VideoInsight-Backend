import os


def extract_audio(file):
    title = file.split('.', -1)[0]
    try:
        os.system('ffmpeg -i ' + '\"' + file + '\"' + ' -vn -c:a copy ' + title + '.aac')
    except Exception as e:
        return "Error: {0}".format(str(e))
    return title + '.aac'
