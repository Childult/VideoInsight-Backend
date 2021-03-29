from unittest import TestCase

from text_analysis.utils import full_to_half, sentence_split


class TestUtils(TestCase):
    def test_full_to_half(self):
        test_data = [
            (
                '你好,世界。',
                '你好，世界。'
            ),
            (
                ',。.?!',
                '，。.？！'
            ),
        ]

        for data in test_data:
            self.assertEqual(data[0], full_to_half(data[1]))

    def test_sentence_split(self):
        test_data = [
            (
                [
                    '第一句。',
                    '第二句!',
                    '第三句?',
                ],
                '第一句。第二句！第三句？'
            ),
            (
                [
                    '句子。',
                    '句子,句子!',
                    '句子?',
                    '句子。',
                ],
                '句子。句子，句子！句子？句子。'
            ),
        ]
        for data in test_data:
            self.assertEqual(data[0], sentence_split(data[1]))

