import json
import unittest
from unittest.mock import MagicMock, patch
from whatsfly.whatsapp import WhatsApp

class TestWhatsAppHandleMessage(unittest.TestCase):
    @patch('whatsfly.whatsapp.new_whatsapp_client_wrapper')
    @patch('whatsfly.whatsapp.threading.Thread')
    def setUp(self, mock_thread, mock_new_client):
        self.whatsapp = WhatsApp()
        # Mocking necessary attributes to avoid errors during tests
        self.whatsapp.c_WhatsAppClientId = 1
        self.whatsapp._userEventHandlers = [MagicMock()]

    def test_handle_message_with_dict_content(self):
        message_data = {
            "eventType": "testEvent",
            "content": {"key": "value"}
        }
        encoded_message = json.dumps(message_data).encode()

        self.whatsapp._handleMessage(encoded_message)

        expected_thandler = {"eventType": "testEvent", "content": {"key": "value"}, "key": "value"}
        self.whatsapp._userEventHandlers[0].assert_called_with(self.whatsapp, expected_thandler)

    def test_handle_message_with_str_content(self):
        message_data = {
            "eventType": "qrCode",
            "content": "some_qr_code_string"
        }
        encoded_message = json.dumps(message_data).encode()

        # Mocking qrcode to avoid actual printing/processing
        with patch('qrcode.QRCode') as mock_qr:
            self.whatsapp._handleMessage(encoded_message)

        expected_thandler = {"eventType": "qrCode", "content": "some_qr_code_string"}
        self.whatsapp._userEventHandlers[0].assert_called_with(self.whatsapp, expected_thandler)

    def test_handle_message_invalid_json(self):
        with self.assertRaises(json.JSONDecodeError):
            self.whatsapp._handleMessage(b"invalid json")

if __name__ == '__main__':
    unittest.main()
