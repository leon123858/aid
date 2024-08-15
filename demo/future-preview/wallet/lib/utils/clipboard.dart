import 'dart:convert';

import 'package:super_clipboard/super_clipboard.dart';

import '../models/aid.dart';

Future<List<AID>> copyRead() async {
  final clipboard = SystemClipboard.instance;
  if (clipboard == null) {
    throw Exception('Clipboard API is not supported on this platform.');
  }
  final reader = await clipboard.read();

  if (reader.canProvide(Formats.plainTextFile)) {
    final jsonStr = await reader.readValue(Formats.plainText);
    if (jsonStr == null) {
      throw Exception('Clipboard does not contain JSON data');
    }
    final json = jsonDecode(jsonStr);
    if (json is! List) {
      throw Exception('Clipboard JSON data is not a list');
    }
    final aids = <AID>[];
    for (final item in json) {
      if (item is! Map) {
        throw Exception('Clipboard JSON data item is not a map');
      }
      final aid = AID.fromJson(item);
      aids.add(aid);
    }
    return aids;
  } else if (reader.canProvide(Formats.plainText)) {
    throw Exception('Clipboard does not contain JSON data');
  }
  throw Exception('Could not get clipboard contents format');
}

Future<void> copyWrite(List<AID> aids) async {
  final clipboard = SystemClipboard.instance;
  if (clipboard == null) {
    throw Exception('Clipboard API is not supported on this platform.');
  }
  final item = DataWriterItem();
  item.add(
      Formats.plainText(jsonEncode(aids.map((aid) => aid.toJson()).toList())));
  await clipboard.write([item]);
}

Future<void> copyWriteString(String str) async {
  final clipboard = SystemClipboard.instance;
  if (clipboard == null) {
    throw Exception('Clipboard API is not supported on this platform.');
  }
  final item = DataWriterItem();
  item.add(Formats.plainText(str));
  await clipboard.write([item]);
}
