# coding: utf-8

from __future__ import absolute_import

from flask import json
from six import BytesIO

from swagger_server.models.job import Job  # noqa: E501
from swagger_server.test import BaseTestCase


class TestJobController(BaseTestCase):
    """JobController integration test stubs"""

    def test_add_job(self):
        """Test case for add_job

        Create a new job
        """
        body = Job()
        response = self.client.open(
            '/v2/job',
            method='POST',
            data=json.dumps(body),
            content_type='application/json')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))

    def test_find_jobs_by_status(self):
        """Test case for find_jobs_by_status

        Finds jobs by status
        """
        query_string = [('status', 'status_example')]
        response = self.client.open(
            '/v2/job/findByStatus',
            method='GET',
            query_string=query_string)
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))

    def test_get_job_by_id(self):
        """Test case for get_job_by_id

        Find job by ID
        """
        response = self.client.open(
            '/v2/job/{jobId}'.format(jobId=789),
            method='GET')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))


if __name__ == '__main__':
    import unittest
    unittest.main()
