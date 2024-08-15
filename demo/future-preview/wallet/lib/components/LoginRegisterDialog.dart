import 'package:flutter/material.dart';
import 'package:getwidget/components/loader/gf_loader.dart';
import 'package:getwidget/types/gf_loader_type.dart';
import 'package:wallet/utils/apiWrapper.dart';

import 'msg.dart';

class LoginRegisterDialog extends StatefulWidget {
  const LoginRegisterDialog({super.key});

  @override
  LoginRegisterDialogState createState() => LoginRegisterDialogState();
}

class LoginRegisterDialogState extends State<LoginRegisterDialog> {
  final TextEditingController _aliasController = TextEditingController();
  final TextEditingController _pinController = TextEditingController();
  final TextEditingController _authKeyController = TextEditingController();
  bool _isLoading = false;
  bool _isLogin = true; // 預設為登入模式
  bool _useAuthKey = false;

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(_isLogin ? '登入' : '註冊'),
      content: SizedBox(
        width: 300, // 固定寬度
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              TextField(
                controller: _aliasController,
                decoration: const InputDecoration(
                  labelText: 'Alias',
                ),
              ),
              const SizedBox(height: 16.0),
              TextField(
                controller: _pinController,
                decoration: const InputDecoration(
                  labelText: 'PIN',
                ),
                keyboardType: TextInputType.number,
                obscureText: true,
              ),
              const SizedBox(height: 16.0),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(_isLogin ? '新用戶？註冊' : '已有帳號？登入'),
                  Switch(
                    value: !_isLogin,
                    onChanged: (value) {
                      setState(() {
                        _isLogin = !value;
                        if (!_isLogin) {
                          _useAuthKey = false; // 註冊模式下重置認證金鑰選項
                        }
                      });
                    },
                  ),
                ],
              ),
              if (_isLogin) // 只在登入模式下顯示認證金鑰選項
                CheckboxListTile(
                  title: const Text('使用認證金鑰'),
                  value: _useAuthKey,
                  onChanged: (bool? value) {
                    setState(() {
                      _useAuthKey = value ?? false;
                    });
                  },
                ),
              if (_isLogin && _useAuthKey)
                TextField(
                  controller: _authKeyController,
                  decoration: const InputDecoration(
                    labelText: '認證金鑰',
                  ),
                ),
              if (_isLoading)
                const Padding(
                  padding: EdgeInsets.only(top: 16.0),
                  child: GFLoader(
                    type: GFLoaderType.circle,
                    loaderColorOne: Colors.blue,
                    loaderColorTwo: Colors.red,
                    loaderColorThree: Colors.green,
                  ),
                ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: _isLoading ? null : () => Navigator.of(context).pop(),
          child: const Text('取消'),
        ),
        ElevatedButton(
          onPressed: _isLoading ? null : _handleSubmit,
          child: Text(_isLogin ? '登入' : '註冊'),
        ),
      ],
    );
  }

  Future<void> _handleSubmit() async {
    if (_aliasController.text.isEmpty || _pinController.text.isEmpty) {
      showErrorToast(context, 'Alias 和 PIN 為必填項');
      return;
    }
    if (_pinController.text.length < 4) {
      showErrorToast(context, 'PIN 最少 4 位數');
      return;
    }
    if (_isLogin && _useAuthKey && _authKeyController.text.isEmpty) {
      showErrorToast(context, '選擇使用認證金鑰時，認證金鑰為必填項');
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      // 這裡應該實現實際的登入或註冊邏輯
      await Future.delayed(const Duration(seconds: 2)); // 模擬網絡請求

      if (!_isLogin) {
        // 註冊邏輯
        var result = await apiWrapper.signup(_aliasController.text, _pinController.text);
        if (!result['result']) {
          throw Exception(result['message']);
        }
        // show success toast with uuid
        if (mounted) {
          // showSuccessToast(context, '註冊成功，New AID: ${result['uuid']}');
          showSuccessToast(context, '註冊成功');
        }
      } else {
        // 登入邏輯
        var result = await apiWrapper.signIn(_aliasController.text, _pinController.text, _useAuthKey? _authKeyController.text : null);
        if (!result['result']) {
          throw Exception(result['message']);
        }
        if (mounted) {
          // showSuccessToast(context, '註冊成功，New AID: ${result['uuid']}');
          showSuccessToast(context, '登入成功');
        }
      }
      if (mounted) {
        Navigator.of(context).pop(true); // 返回 true 表示操作成功
      }
    } catch (e) {
      if (mounted) {
        showErrorToast(context, "${_isLogin ? '登入' : '註冊'}失敗：$e");
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }
}
