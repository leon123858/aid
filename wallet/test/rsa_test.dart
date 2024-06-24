import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';
import 'package:crypton/crypton.dart';

// 假设这个文件包含 RSAUtils 类和相关的类
import 'package:wallet/utils/rsa.dart';

void main() {
  group('RSAUtils Tests', () {
    test('generateRSAKeyPair creates valid key pair', () {
      final keyPair = RSAUtils.generateRSAKeyPair(null);
      expect(keyPair, isA<KeyPairRSA>());
      expect(keyPair.publicKey, isA<RSAPublicKey>());
      expect(keyPair.privateKey, isA<RSAPrivateKey>());
    });

    test('generateRSAKeyPairFromPem creates valid key pair', () {
      final keyPair = RSAUtils.generateRSAKeyPair(null);
      final publicKeyPem = RSAUtils.encodePublicKeyToPem(keyPair.publicKey);
      final privateKeyPem = RSAUtils.encodePrivateKeyToPem(keyPair.privateKey);

      final regeneratedKeyPair = RSAUtils.generateRSAKeyPairFromPem(privateKeyPem, publicKeyPem);
      expect(regeneratedKeyPair, isA<KeyPairRSA>());
      expect(regeneratedKeyPair.publicKey.toPEM(), equals(publicKeyPem));
      expect(regeneratedKeyPair.privateKey.toPEM(), equals(privateKeyPem));
    });

    test('rsaSign and rsaVerify work correctly', () {
      final keyPair = RSAUtils.generateRSAKeyPair(null);
      const dataToSign = 'Hello, RSA!';

      final signature = RSAUtils.rsaSign(keyPair.privateKey, dataToSign);
      expect(signature, isA<Uint8List>());

      final isValid = RSAUtils.rsaVerify(keyPair.publicKey, dataToSign, signature);
      expect(isValid, isTrue);

      // Test with invalid data
      final isInvalid = RSAUtils.rsaVerify(keyPair.publicKey, 'Invalid data', signature);
      expect(isInvalid, isFalse);
    });

    test('encodePublicKeyToPem and parsePublicKeyFromPem work correctly', () {
      final keyPair = RSAUtils.generateRSAKeyPair(null);
      final publicKeyPem = RSAUtils.encodePublicKeyToPem(keyPair.publicKey);
      expect(publicKeyPem, isA<String>());
      expect(publicKeyPem, contains('BEGIN PUBLIC KEY'));

      final parsedPublicKey = RSAUtils.parsePublicKeyFromPem(publicKeyPem);
      expect(parsedPublicKey, isA<RSAPublicKey>());
      expect(parsedPublicKey.toPEM(), equals(publicKeyPem));
    });

    test('encodePrivateKeyToPem and parsePrivateKeyFromPem work correctly', () {
      final keyPair = RSAUtils.generateRSAKeyPair(null);
      final privateKeyPem = RSAUtils.encodePrivateKeyToPem(keyPair.privateKey);
      expect(privateKeyPem, isA<String>());
      expect(privateKeyPem, contains('BEGIN RSA PRIVATE KEY'));

      final parsedPrivateKey = RSAUtils.parsePrivateKeyFromPem(privateKeyPem);
      expect(parsedPrivateKey, isA<RSAPrivateKey>());
      expect(parsedPrivateKey.toPEM(), equals(privateKeyPem));
    });

    test('KeyPairString and KeyPairRSA conversion works correctly', () {
      final rsaKeyPair = RSAUtils.generateRSAKeyPair(null);
      final keyPairString = rsaKeyPair.toKeyPairString();

      expect(keyPairString, isA<KeyPairString>());
      expect(keyPairString.type, equals('RSA'));
      expect(keyPairString.publicKey, contains('BEGIN PUBLIC KEY'));
      expect(keyPairString.privateKey, contains('BEGIN RSA PRIVATE KEY'));

      final convertedRsaKeyPair = keyPairString.toRSAPair();
      expect(convertedRsaKeyPair, isA<KeyPairRSA>());
      expect(convertedRsaKeyPair.publicKey.toPEM(), equals(rsaKeyPair.publicKey.toPEM()));
      expect(convertedRsaKeyPair.privateKey.toPEM(), equals(rsaKeyPair.privateKey.toPEM()));
    });
  });
}