import sys, os, you_get, json


def download_video(url: str, path: str) -> str:
    cmd = 'you-get -o ' + path + ' ' + url + ' --json'
    r = os.popen(cmd)
    text = r.buffer.read().decode(encoding='utf8')
    parsed_text = json.loads(text)

    # get title
    title = parsed_text['title']

    # get accessible qualities (480p default), format and container
    qualities = ['480', '720', '360', '240', '1080']
    for quality in qualities:
        formats = [k for k, v in parsed_text['streams'].items() if quality in v['quality']]
        if formats:
            break
    if not formats:
        formats = [k for k, v in parsed_text['streams'].items()]
    formats.sort()
    fmt = formats[0]
    container = parsed_text['streams'][fmt]['container']
    
    sys.argv = ['you-get', '-o', path, '-O', 'output', url, '--format=%s' % fmt]
    try:
        you_get.main()
    except Exception as e:
        return "Error: {0}".format(str(e))

    return title + '.' + container