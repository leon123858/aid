import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/aid.dart';

class AIDListView extends StatelessWidget {
  const AIDListView({super.key});

  @override
  Widget build(BuildContext context) {
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
                        onPressed: () {
                          // 實現登錄邏輯
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
}
