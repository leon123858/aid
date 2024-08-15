import 'dart:convert';

import 'package:crypton/crypton.dart';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:uuid/uuid.dart';
import 'package:wallet/constants/config.dart';
import 'package:wallet/utils/device.dart';
import 'package:wallet/utils/rsa.dart';

class AIDApiClient {
  final String baseUrl;
  final uuid = const Uuid();
  final http.Client _httpClient = http.Client();
  final _deviceInfo = DeviceInfo();
  var _deviceIP = '';

  var _isInit = false;

  String get deviceInfoHash => _deviceInfo.deviceHash;
  String get deviceIP => _deviceIP;

  AIDApiClient({this.baseUrl = 'http://127.0.0.1:8080'});

  Future<void> init() async {
    if (_isInit) {
      return;
    }
    await _deviceInfo.initPlatformState();
    final response = await http.get(Uri.parse('https://api.ipify.org?format=json'));
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      _deviceIP = data['ip'];
      if (kDebugMode) {
        print('IP address: $_deviceIP');
      }
    } else {
      // warn user
      if (kDebugMode) {
        print('Failed to get IP address');
      }
      _deviceIP = '';
    }
    _isInit = true;
  }

  String generateAID() {
    return uuid.v4();
  }

  Future<Map<String, dynamic>> login(
      String aid, RSAPrivateKey privateKey) async {
    final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
    final sign = RSAUtils.rsaSign(privateKey, timestamp);
    final b64Sign = base64.encode(sign);
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/api/login'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        'aid': aid,
        'browser': _deviceInfo.deviceHash,
        'ip': _deviceIP,
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
        'ip': _deviceIP,
        'publicKey': RSAUtils.encodePublicKeyToPem(publicKey),
      }),
    );
    return _handleResponse(response);
  }

  Future<Map<String, dynamic>> signup(String alias, String pin) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/usage/register'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        "fingerprint": _deviceInfo.deviceHash,
        "ip": _deviceIP,
        "password": pin,
        "token": "",
        "username": alias
      }),
    );
    return _handleResponse(response);
  }

  Future<Map<String, dynamic>> signIn(String alias, String pin, String? token) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/usage/login'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        "fingerprint": _deviceInfo.deviceHash,
        "ip": _deviceIP,
        "password": pin,
        "token": token ?? "",
        "username": alias
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
