from audio_analysis.api import audio_to_text
from swagger_server.models.job import Job
from text_analysis.api import text_summarize
from video_analysis.api import extract_audio, extract_key_frame, video_summarize
from video_getter.api import download_video


def add_job(url):
    """Create a new job

    :param url: The url of the video
    :type url: str

    :rtype: None
    """
    video_path = download_video(url)
    audio_path = extract_audio(video_path)
    audio_text = audio_to_text(audio_path)
    text_sum = text_summarize(audio_text)
    video_sum = extract_key_frame(video_summarize(video_path))
    return 'Not implemented.'


def find_jobs_by_status(status):
    """Finds jobs by status

    Multiple status values can be provided with comma separated strings

    :param status: Status values that need to be considered for filter
    :type status: List[str]

    :rtype: List[Job]
    """
    return 'Not implemented.'


def get_job_by_id(jobId):
    """Find job by ID

    Returns a single job

    :param jobId: ID of job to return
    :type jobId: int

    :rtype: Job
    """
    return 'Not implemented.'
