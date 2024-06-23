import 'dart:convert';
import 'dart:math';
import 'dart:typed_data';

import 'package:flutter/foundation.dart';
import 'package:pointycastle/asn1/primitives/asn1_bit_string.dart';
import 'package:pointycastle/asn1/primitives/asn1_integer.dart';
import 'package:pointycastle/asn1/primitives/asn1_null.dart';
import 'package:pointycastle/asn1/primitives/asn1_object_identifier.dart';
import 'package:pointycastle/asn1/primitives/asn1_sequence.dart';
import 'package:pointycastle/export.dart';

// How to use:
// final keyPair = RSAUtils.generateRSAKeyPair();
// final publicKey = RSAUtils.encodePublicKeyToPem(keyPair.publicKey);
// final privateKey = RSAUtils.encodePrivateKeyToPem(keyPair.privateKey);
// print(publicKey);
// print(privateKey);
// // sign
// final data = utf8.encode('Hello, World!');
// final signature = RSAUtils.rsaSign(keyPair.privateKey, data);
// print(base64.encode(signature));
// // verify
// final result = RSAUtils.rsaVerify(keyPair.publicKey, data, signature);
// print(result);

AsymmetricKeyPair<RSAPublicKey, RSAPrivateKey> _generateRSAKeyPair(_) {
  // 創建一個安全隨機數生成器
  final secureRandom = SecureRandom('Fortuna')
    ..seed(KeyParameter(Uint8List.fromList(
        List.generate(32, (_) => Random.secure().nextInt(255)))));

  // 創建 RSA 密鑰生成器
  final keyGen = RSAKeyGenerator()
    ..init(ParametersWithRandom(
      RSAKeyGeneratorParameters(BigInt.parse('65537'), 2048, 64),
      secureRandom,
    ));

  // 生成密鑰對
  final pair = keyGen.generateKeyPair();

  // 返回 AsymmetricKeyPair<RSAPublicKey, RSAPrivateKey>
  return AsymmetricKeyPair<RSAPublicKey, RSAPrivateKey>(
    pair.publicKey as RSAPublicKey,
    pair.privateKey as RSAPrivateKey,
  );
}

String _encodePublicKeyToPem(RSAPublicKey publicKey) {
  final algorithmIdentifier = ASN1Sequence();
  algorithmIdentifier.add(ASN1ObjectIdentifier([1, 2, 840, 113549, 1, 1, 1]));
  algorithmIdentifier.add(ASN1Null());

  final publicKeySequence = ASN1Sequence();
  publicKeySequence.add(ASN1Integer(publicKey.modulus));
  publicKeySequence.add(ASN1Integer(publicKey.exponent));

  final topLevelSequence = ASN1Sequence();
  topLevelSequence.add(algorithmIdentifier);
  topLevelSequence.add(ASN1BitString(stringValues: publicKeySequence.encode()));

  final dataBase64 = base64.encode(topLevelSequence.encode());
  return '-----BEGIN PUBLIC KEY-----\n$dataBase64\n-----END PUBLIC KEY-----';
}

String _encodePrivateKeyToPem(RSAPrivateKey privateKey) {
  final version = ASN1Integer(BigInt.from(0));

  final sequence = ASN1Sequence();
  sequence.add(version);
  sequence.add(ASN1Integer(privateKey.n));
  sequence.add(ASN1Integer(privateKey.exponent));
  sequence.add(ASN1Integer(privateKey.privateExponent!));
  sequence.add(ASN1Integer(privateKey.p));
  sequence.add(ASN1Integer(privateKey.q));
  sequence.add(ASN1Integer(
      privateKey.privateExponent! % (privateKey.p! - BigInt.from(1))));
  sequence.add(ASN1Integer(
      privateKey.privateExponent! % (privateKey.q! - BigInt.from(1))));
  sequence.add(ASN1Integer(privateKey.q?.modInverse(privateKey.p!)));

  final dataBase64 = base64.encode(sequence.encode());
  return '-----BEGIN PRIVATE KEY-----\n$dataBase64\n-----END PRIVATE KEY-----';
}

Uint8List _rsaSign(RSAPrivateKey privateKey, Uint8List dataToSign) {
  //final signer = Signer('SHA-256/RSA'); // Get using registry
  final signer = RSASigner(SHA256Digest(), '0609608648016503040201');

  // initialize with true, which means sign
  signer.init(true, PrivateKeyParameter<RSAPrivateKey>(privateKey));

  final sig = signer.generateSignature(dataToSign);

  return sig.bytes;
}

bool _rsaVerify(
    RSAPublicKey publicKey, Uint8List signedData, Uint8List signature) {
  //final signer = Signer('SHA-256/RSA'); // Get using registry
  final sig = RSASignature(signature);

  final verifier = RSASigner(SHA256Digest(), '0609608648016503040201');

  // initialize with false, which means verify
  verifier.init(false, PublicKeyParameter<RSAPublicKey>(publicKey));

  try {
    return verifier.verifySignature(signedData, sig);
  } on ArgumentError {
    return false; // for Pointy Castle 1.0.2 when signature has been modified
  }
}

class RSAUtils {
  // 生成 RSA 密鑰對
  static AsymmetricKeyPair<RSAPublicKey, RSAPrivateKey> generateRSAKeyPair(_) {
    return _generateRSAKeyPair(null);
  }

  // 將 RSAPublicKey 轉換為 PEM 格式
  static String encodePublicKeyToPem(RSAPublicKey publicKey) {
    return _encodePublicKeyToPem(publicKey);
  }

  // 將 RSAPrivateKey 轉換為 PEM 格式
  static String encodePrivateKeyToPem(RSAPrivateKey privateKey) {
    return _encodePrivateKeyToPem(privateKey);
  }

  // 使用 RSA 私鑰對數據進行簽名
  static Uint8List rsaSign(RSAPrivateKey privateKey, Uint8List dataToSign) {
    return _rsaSign(privateKey, dataToSign);
  }

  // 使用 RSA 公鑰驗證數據簽名
  static bool rsaVerify(
      RSAPublicKey publicKey, Uint8List signedData, Uint8List signature) {
    return _rsaVerify(publicKey, signedData, signature);
  }
}
