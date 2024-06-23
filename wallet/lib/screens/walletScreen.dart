import 'package:flutter/material.dart';
import 'package:getwidget/components/search_bar/gf_search_bar.dart';
import 'package:wallet/screens/appBar.dart';
import 'package:wallet/screens/drawer.dart';

import 'aidList.dart';

class AIDWalletScreen extends StatelessWidget {
  const AIDWalletScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: appBar(context, "AID Wallet"),
      drawer: drawer(context),
      body: Column(
        children: [
          Container(
            color: Colors.white,
            padding: const EdgeInsets.all(16.0),
            child: GFSearchBar(
              searchList: const ['search'],
              hideSearchBoxWhenItemSelected: false,
              overlaySearchListHeight: 0,
              searchQueryBuilder: (query, list) {
                return list
                    .where((item) =>
                        item.toLowerCase().contains(query.toLowerCase()))
                    .toList();
              },
              overlaySearchListItemBuilder: (String item) {
                return Container();
              },
              searchBoxInputDecoration: InputDecoration(
                hintText: 'Search AID',
                hintStyle: const TextStyle(color: Colors.grey),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10.0),
                  borderSide: BorderSide.none,
                ),
                filled: true,
                fillColor: Colors.grey[200],
              ),
            ),
          ),
          const Expanded(
            child: AIDListView(),
          ),
        ],
      ),
    );
  }
}