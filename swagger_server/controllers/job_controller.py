import connexion

from swagger_server.models.job import Job  # noqa: E501


def add_job(body):  # noqa: E501
    """Create a new job

     # noqa: E501

    :param body: Pet object that needs to be added to the store
    :type body: dict | bytes

    :rtype: None
    """
    if connexion.request.is_json:
        body = Job.from_dict(connexion.request.get_json())  # noqa: E501
    return 'do some magic!'


def find_jobs_by_status(status):  # noqa: E501
    """Finds jobs by status

    Multiple status values can be provided with comma separated strings # noqa: E501

    :param status: Status values that need to be considered for filter
    :type status: List[str]

    :rtype: List[Job]
    """
    return 'do some magic!'


def get_job_by_id(jobId):  # noqa: E501
    """Find job by ID

    Returns a single job # noqa: E501

    :param jobId: ID of job to return
    :type jobId: int

    :rtype: Job
    """
    return 'do some magic!'
