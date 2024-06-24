import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:wallet/utils/apiWrapper.dart';
import 'package:wallet/utils/rsa.dart';

import '../../components/msg.dart';
import '../../models/aid.dart';

Consumer<AIDListModel> aidListView() {
  return Consumer<AIDListModel>(
    builder: (context, aidListModel, child) {
      return Container(
        color: Colors.grey[100],
        child: ListView.separated(
          itemCount: aidListModel.aidCount,
          separatorBuilder: (context, index) => Divider(
            color: Colors.grey[300],
            height: 1.0,
          ),
          itemBuilder: (context, index) {
            final aid = aidListModel.aidList[index];
            return Container(
              color: Colors.white,
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Row(
                  children: [
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            aid.name,
                            style: const TextStyle(
                              color: Colors.black,
                              fontSize: 18.0,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 8.0),
                          Text(
                            aid.description,
                            style: const TextStyle(color: Colors.grey),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(width: 16.0),
                    ElevatedButton(
                      onPressed: () async {
                        // 實現登錄邏輯
                        try {
                          final key =
                          RSAUtils.parsePrivateKeyFromPem(aid.privateKey);
                          final response = await apiWrapper.login(aid.aid, key);
                          if (context.mounted) {
                            if (response['result']) {
                              showSuccessToast(context, 'Login success');
                            } else {
                              throw Exception(response['content']);
                            }
                          }
                        } catch (e) {
                          if (context.mounted) {
                            showErrorToast(context, "throw error msg:$e");
                          }
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(20.0),
                        ),
                      ),
                      child: const Text('Login'),
                    ),
                    const SizedBox(width: 16.0),
                    IconButton(
                      icon: const Icon(Icons.copy),
                      color: Colors.grey,
                      onPressed: () {
                        // 實現複製邏輯
                      },
                    ),
                  ],
                ),
              ),
            );
          },
        ),
      );
    },
  );
}
