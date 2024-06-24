import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:pointycastle/asymmetric/api.dart';
import 'package:uuid/uuid.dart';
import 'package:wallet/constants/config.dart';
import 'package:wallet/utils/device.dart';
import 'package:wallet/utils/rsa.dart';

// How to use:
// final keyPair = RSAUtils.generateRSAKeyPair();
// // http request to 20.2.209.109
// final apiWrapper = AIDApiClient(baseUrl: 'http://20.2.209.109');
// final aid = apiWrapper.generateAID();
// final response = await apiWrapper.register(aid, keyPair.publicKey);
// print(response);
// final response2 = await apiWrapper.login(aid, keyPair.privateKey);
// print(response2);

class AIDApiClient {
  final String baseUrl;
  final uuid = const Uuid();
  final http.Client _httpClient = http.Client();
  final _deviceInfo = DeviceInfo();

  var _isInit = false;

  AIDApiClient({this.baseUrl = 'http://127.0.0.1:8080'});

  Future<void> init() async {
    if (_isInit) {
      return;
    }
    await _deviceInfo.initPlatformState();
    _isInit = true;
  }

  String generateAID() {
    return uuid.v4();
  }

  Future<Map<String, dynamic>> login(
      String aid, RSAPrivateKey privateKey) async {
    final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
    final sign = RSAUtils.rsaSign(privateKey, utf8.encode(timestamp));
    final b64Sign = base64.encode(sign);
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/api/login'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        'aid': aid,
        'browser': _deviceInfo.deviceHash,
        'ip': "",
        'sign': b64Sign,
        'timestamp': timestamp,
      }),
    );
    return _handleResponse(response);
  }

  Future<Map<String, dynamic>> register(
      String aid, RSAPublicKey publicKey) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/api/register'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        'aid': aid,
        'browser': _deviceInfo.deviceHash,
        'ip': "",
        'publicKey': RSAUtils.encodePublicKeyToPem(publicKey),
      }),
    );
    return _handleResponse(response);
  }

  Map<String, dynamic> _handleResponse(http.Response response) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      return json.decode(response.body);
    } else {
      throw AIDApiException(response.statusCode, response.body);
    }
  }

  void dispose() {
    _httpClient.close();
  }
}

class AIDApiException implements Exception {
  final int statusCode;
  final String body;

  AIDApiException(this.statusCode, this.body);

  @override
  String toString() => 'AIDApiException: $statusCode\n$body';
}

final apiWrapper = AIDApiClient(baseUrl: serverUrl);
