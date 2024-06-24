import 'dart:convert';
import 'dart:typed_data';

import 'package:crypton/crypton.dart';

// import 'package:crypton/crypton.dart';
// import 'package:flutter/foundation.dart';
// import 'package:pointycastle/asn1/primitives/asn1_bit_string.dart';
// import 'package:pointycastle/asn1/primitives/asn1_integer.dart';
// import 'package:pointycastle/asn1/primitives/asn1_null.dart';
// import 'package:pointycastle/asn1/primitives/asn1_object_identifier.dart';
// import 'package:pointycastle/asn1/primitives/asn1_sequence.dart';
// import 'package:pointycastle/export.dart';

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

class KeyPairString {
  String type;
  String publicKey;
  String privateKey;

  KeyPairString(this.type, this.publicKey, this.privateKey);

  KeyPairRSA toRSAPair() {
    return KeyPairRSA(
      RSAPublicKey.fromPEM(publicKey),
      RSAPrivateKey.fromPEM(privateKey),
    );
  }
}

class KeyPairRSA {
  RSAPublicKey publicKey;
  RSAPrivateKey privateKey;

  KeyPairRSA(this.publicKey, this.privateKey);

  KeyPairString toKeyPairString() {
    return KeyPairString(
      'RSA',
      publicKey.toPEM(),
      privateKey.toPEM(),
    );
  }
}

KeyPairRSA _generateRSAKeyPair(_) {
  RSAKeypair rsaKeypair = RSAKeypair.fromRandom();

  // 返回
  return KeyPairRSA(
    rsaKeypair.publicKey,
    rsaKeypair.privateKey,
  );
}

class RSAUtils {
  // 生成 RSA 密鑰對
  static KeyPairRSA generateRSAKeyPair(_) {
    return _generateRSAKeyPair(null);
  }

  static KeyPairRSA generateRSAKeyPairFromPem(
      String privateKey, String publicKey) {
    return KeyPairRSA(
      RSAPublicKey.fromPEM(publicKey),
      RSAPrivateKey.fromPEM(privateKey),
    );
  }

  // 使用 RSA 私鑰對數據進行簽名
  static Uint8List rsaSign(RSAPrivateKey privateKey, String dataToSign) {
    return privateKey.createSHA256Signature(utf8.encode(dataToSign));
  }

  // 使用 RSA 公鑰驗證數據簽名
  static bool rsaVerify(
      RSAPublicKey publicKey, String signedData, Uint8List signature) {
    return publicKey.verifySHA256Signature(utf8.encode(signedData), signature);
  }

  static String encodePublicKeyToPem(key) {
    return key.toPEM() as String;
  }

  static String encodePrivateKeyToPem(key) {
    return key.toPEM() as String;
  }

  static RSAPrivateKey parsePrivateKeyFromPem(String pem) {
    return RSAPrivateKey.fromPEM(pem);
  }

  static RSAPublicKey parsePublicKeyFromPem(String pem) {
    return RSAPublicKey.fromPEM(pem);
  }
}
