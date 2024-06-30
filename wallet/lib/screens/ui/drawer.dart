import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:wallet/components/LoginRegisterDialog.dart';
import 'package:wallet/components/msg.dart';
import 'package:wallet/constants/config.dart';
import 'package:wallet/utils/apiWrapper.dart';
import 'package:wallet/utils/clipboard.dart';

import '../../models/aid.dart';

Drawer walletDrawer(BuildContext context) {
  return Drawer(
    child: Container(
      color: Colors.black26,
      child: ListView(
        children: <Widget>[
          const DrawerHeader(
              decoration: BoxDecoration(
                color: Colors.white,
              ),
              child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      "Lab408",
                      style: TextStyle(
                        color: Colors.black,
                        fontSize: 24.0,
                        fontWeight: FontWeight.bold,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    SizedBox(height: 8.0),
                    Text(
                      'AID Wallet',
                      style: TextStyle(
                        color: Colors.black,
                        fontSize: 18.0,
                        fontStyle: FontStyle.italic,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ])),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton(
              onPressed: () async {
                final provider =
                    Provider.of<AIDListModel>(context, listen: false);
                final list = await copyRead();
                await provider.importAIDList(list);
                if (context.mounted) {
                  showSuccessToast(context, 'Wallet copied to clipboard');
                }
              },
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 40),
              ),
              child: const Text('Import Wallet'),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton(
              onPressed: () async {
                final provider =
                    Provider.of<AIDListModel>(context, listen: false);
                await provider.initAIDList();
                await copyWrite(provider.exportAIDList);
                if (context.mounted) {
                  showSuccessToast(context, 'Get copy from clipboard');
                }
              },
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 40),
              ),
              child: const Text('Export Wallet'),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton(
              onPressed: () async {
                // Handle button press
                final provider =
                    Provider.of<AIDListModel>(context, listen: false);
                await provider.initAIDList();
                await provider.clearAIDList();
              },
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 40),
              ),
              child: const Text('Clear Wallet'),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton(
              onPressed: () async {
                try {
                  await copyWriteString(apiWrapper.deviceInfoHash);
                  if (context.mounted) {
                    showSuccessToast(context, 'Get copy from clipboard');
                  }
                } catch (e) {
                  if (context.mounted) {
                    showErrorToast(context, "$e");
                  }
                }
              },
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 40),
              ),
              child: const Text('Copy Device fingerprint'),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton(
              onPressed: () async {
                // Handle button press
                const url = documentUrl;
                Uri uri = Uri.parse(url);
                if (await canLaunchUrl(uri)) {
                  await launchUrl(uri);
                } else {
                  throw 'Could not launch $url';
                }
              },
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 40),
              ),
              child: const Text('Documentation',
                  style: TextStyle(color: Colors.grey)),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: ElevatedButton(
              onPressed: () async {
                showDialog(
                    context: context,
                    builder: (context) => const LoginRegisterDialog());
              },
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 40),
              ),
              child:
                  const Text('Demo Page', style: TextStyle(color: Colors.grey)),
            ),
          ),
        ],
      ),
    ),
  );
}
