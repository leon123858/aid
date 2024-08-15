import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:wallet/components/msg.dart';
import 'package:wallet/screens/ui/aidSearchBar.dart';
import 'package:wallet/screens/ui/appBar.dart';
import 'package:wallet/screens/ui/drawer.dart';

import '../models/aid.dart';
import 'ui/aidList.dart';

class AIDWalletScreen extends StatelessWidget {
  const AIDWalletScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: appBar(context, "AID Wallet"),
      drawer: walletDrawer(context),
      body: Column(
        children: [
          Container(
            color: Colors.white,
            padding: const EdgeInsets.all(16.0),
            child: SimpleAIDSearchBar(
              onSearch: (value) async {
                final queryStr = value.toLowerCase();
                final provider =
                    Provider.of<AIDListModel>(context, listen: false);
                try {
                  await provider.getAIDByKeyword(queryStr);
                  if (context.mounted) {
                    showSuccessToast(context, "Search completed");
                  }
                } catch (e) {
                  if (context.mounted) {
                    showErrorToast(context, e.toString());
                  }
                }
              },
            ),
          ),
          Expanded(
            child: aidListView(),
          ),
        ],
      ),
    );
  }
}
