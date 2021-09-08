import easyocr

OCR_READER = easyocr.Reader(['ch_sim', 'en'])  # need to run only once to load model into memory


def ocr(filename: str) -> [str]:
    return OCR_READER.readtext(filename, detail=0)
