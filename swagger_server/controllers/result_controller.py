import connexion

from swagger_server.models.result import Result  # noqa: E501
from swagger_server import util


def get_result_by_id(resultId):  # noqa: E501
    """Get result by ID

    Returns a single result of a job # noqa: E501

    :param resultId: ID of result that needs to be fetched
    :type resultId: int

    :rtype: Result
    """
    return 'do some magic!'
