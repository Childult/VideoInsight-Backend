import base64
import json
import time

from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.common.exception.tencent_cloud_sdk_exception import TencentCloudSDKException
from tencentcloud.asr.v20190614 import asr_client, models


class RequestAPI:
    cred = credential.Credential("AKIDCF5h4DJrLcgW90p9UgeOvFLEsIUUao9G", "1oSt54pNhnfPwFKBpFRQBj4D9rMdhX9X")
    httpProfile = HttpProfile()
    httpProfile.endpoint = "asr.tencentcloudapi.com"

    def __init__(self):
        self.clientProfile = ClientProfile()
        self.clientProfile.httpProfile = RequestAPI.httpProfile
        self.client = asr_client.AsrClient(RequestAPI.cred, "", self.clientProfile)

    def create_task(self, url):
        try:
            req = models.CreateRecTaskRequest()
            params = {
                "EngineModelType": "16k_zh_video",  # 非电话场景，音视频领域
                "ChannelNum": 1,  # 单声道
                "SpeakerDiarization": 0,  # 不开启说话人分离
                "ResTextFormat": 0,
                "SourceType": 0,
                "Url": url,
                "FilterModal": 2  # 严格过滤语气词
            }
            req.from_json_string(json.dumps(params))

            resp = self.client.CreateRecTask(req)
            return json.loads(resp.to_json_string())['Data']['TaskId']

        except TencentCloudSDKException as err:
            raise err

    def query_result(self, task_id: str):
        try:
            req = models.DescribeTaskStatusRequest()
            params = {
                "TaskId": task_id
            }
            req.from_json_string(json.dumps(params))

            resp = self.client.DescribeTaskStatus(req)
            return json.loads(resp.to_json_string())['Data']
        except TencentCloudSDKException as err:
            raise err

    def audio_to_text(self, file: str):
        try:
            with open(file, 'rb') as f:
                base64_data = base64.b64encode(f.read()).decode('utf-8')
                task_id = self.create_task(base64_data)

                while True:
                    result_data = self.query_result(task_id)
                    if result_data['Status'] == 2:  # 任务成功
                        return result_data['Result']
                    elif result_data['Status'] == 3:  # 任务失败
                        return result_data['ErrorMsg']

                    time.sleep(3)
        except Exception as e:
            raise e

    def url_to_text(self, url: str):
        try:
            task_id = self.create_task(url)
            while True:
                result_data = self.query_result(task_id)
                if result_data['Status'] == 2:  # 任务成功
                    return result_data['Result']
                elif result_data['Status'] == 3:  # 任务失败
                    return result_data['ErrorMsg']

                time.sleep(3)
        except Exception as e:
            print(e)


if __name__ == '__main__':
    api = RequestAPI()
    try:
        print(api.url_to_text('https://files-1253366698.cos.ap-beijing.myqcloud.com/yydp.mp3'))
        # print(api.audio_to_text('../test.mp3'))
    except Exception as e:
        print('error: ', e)
