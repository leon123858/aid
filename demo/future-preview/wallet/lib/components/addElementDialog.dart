import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:getwidget/components/loader/gf_loader.dart';
import 'package:getwidget/types/gf_loader_type.dart';
import 'package:provider/provider.dart';

import '../models/aid.dart';
import '../utils/apiWrapper.dart';
import '../utils/rsa.dart';
import 'msg.dart';

class AddElementDialog extends StatefulWidget {
  const AddElementDialog({super.key});

  @override
  AddElementDialogState createState() => AddElementDialogState();
}

class AddElementDialogState extends State<AddElementDialog> {
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _descriptionController = TextEditingController();
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Add Element'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          TextField(
            controller: _nameController,
            decoration: const InputDecoration(
              labelText: 'Name',
            ),
          ),
          const SizedBox(height: 16.0),
          TextField(
            controller: _descriptionController,
            decoration: const InputDecoration(
              labelText: 'Description',
            ),
          ),
          if (_isLoading)
            const GFLoader(
              type: GFLoaderType.circle,
              loaderColorOne: Colors.blue,
              loaderColorTwo: Colors.red,
              loaderColorThree: Colors.green,
            ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: _isLoading
              ? null
              : () {
                  Navigator.of(context).pop();
                },
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: _isLoading ? null : _handleAddElement,
          child: const Text('Add'),
        ),
      ],
    );
  }

  Future<void> _handleAddElement() async {
    if (_nameController.text.isEmpty || _descriptionController.text.isEmpty) {
      showErrorToast(context, 'Name and Description are required');
      return;
    }
    final provider = Provider.of<AIDListModel>(context, listen: false);
    final navigator = Navigator.of(context);
    try {
      setState(() {
        _isLoading = true;
      });
      final aid = apiWrapper.generateAID();
      final rsaKeyPair = await compute(RSAUtils.generateRSAKeyPair, null);
      final response = await apiWrapper.register(aid, rsaKeyPair.publicKey);
      if (!response['result']) {
        throw Exception(response['content']);
      }
      await provider.addAID(AID(
        aid: aid,
        name: _nameController.text,
        description: _descriptionController.text,
        publicKey: RSAUtils.encodePublicKeyToPem(rsaKeyPair.publicKey),
        privateKey: RSAUtils.encodePrivateKeyToPem(rsaKeyPair.privateKey),
      ));
      navigator.pop();
    } catch (e) {
      if (mounted) {
        showErrorToast(context, "Failed to add element: $e");
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
