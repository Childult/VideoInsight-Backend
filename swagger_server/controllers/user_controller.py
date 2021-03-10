import connexion
import six

from swagger_server.models.user import User
from swagger_server import util


def create_user(body):
    """Create user

    This can only be done by the logged in user.

    :param body: Created user object
    :type body: dict | bytes

    :rtype: None
    """
    if connexion.request.is_json:
        body = User.from_dict(connexion.request.get_json())
    return 'do some magic!'


def login_user(username, password):
    """Logs user into the system

    :param username: The user name for login
    :type username: str
    :param password: The password for login in clear text
    :type password: str

    :rtype: str
    """
    return 'do some magic!'


def logout_user():
    """Logs out current logged in user session

    :rtype: None
    """
    return 'do some magic!'
