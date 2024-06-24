import 'package:flutter/material.dart';
import 'package:getwidget/components/appbar/gf_appbar.dart';

import '../../components/addElementDialog.dart';

GFAppBar appBar(BuildContext context, String title) {
  return GFAppBar(
    title: const Text(
      'AID Wallet',
      style: TextStyle(
        color: Colors.black,
        fontSize: 24.0,
        fontWeight: FontWeight.bold,
      ),
    ),
    centerTitle: true,
    backgroundColor: Colors.white,
    elevation: 2.0,
    actions: [
      Padding(
        padding: const EdgeInsets.only(right: 16.0),
        child: IconButton(
          icon: const Icon(
            Icons.add_circle_outline,
            color: Colors.blue,
            size: 28.0,
          ),
          onPressed: () {
            showDialog(
              context: context,
              builder: (context) => const AddElementDialog(),
            );
          },
        ),
      ),
    ],
    leading: Builder(
      builder: (context) => IconButton(
        icon: const Icon(
          Icons.menu,
          color: Colors.black,
          size: 28.0,
        ),
        onPressed: () {
          Scaffold.of(context).openDrawer();
        },
      ),
    ),
  );
}
