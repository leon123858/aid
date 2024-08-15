import 'package:flutter/material.dart';
import 'package:hive/hive.dart';

part 'aid.g.dart'; // should run `flutter pub run build_runner build` when can not find this file

@HiveType(typeId: 0)
class AID extends HiveObject {
  @HiveField(0)
  final String aid;

  @HiveField(1)
  final String name;

  @HiveField(2)
  final String description;

  @HiveField(3)
  final String publicKey;

  @HiveField(4)
  final String privateKey;

  AID({
    required this.aid,
    required this.name,
    required this.description,
    required this.publicKey,
    required this.privateKey,
  });

  static fromJson(Map item) {
    return AID(
      aid: item['aid'],
      name: item['name'],
      description: item['description'],
      publicKey: item['publicKey'],
      privateKey: item['privateKey'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'aid': aid,
      'name': name,
      'description': description,
      'publicKey': publicKey,
      'privateKey': privateKey,
    };
  }
}

class AIDListModel extends ChangeNotifier {
  final List<AID> _aidList = [];
  late Box<AID> _aidBox;

  // getter all values in db
  List<AID> get aidList => _aidList;

  int get aidCount => _aidList.length;

  List<AID> get exportAIDList => _aidBox.values.toList();

  Future<void> importAIDList(List<AID> value) async {
    // clear UI List
    _aidList.clear();
    // add to UI List
    _aidList.addAll(value);
    // clear db
    await _aidBox.clear();
    // add to db
    for (var aid in value) {
      await _aidBox.put(aid.aid, aid);
    }
    notifyListeners();
  }

  Future<void> addAID(AID aid) async {
    // add to UI List
    _aidList.add(aid);
    // add to db
    await _aidBox.put(aid.aid, aid);
    notifyListeners();
  }

  Future<void> clearAIDList() async {
    // clear UI List
    _aidList.clear();
    // clear db
    await _aidBox.clear();
    notifyListeners();
  }

  Future<void> getAIDByKeyword(String keyword) async {
    // clear UI List
    _aidList.clear();
    // get from db
    _aidBox.values.toList().forEach((aid) {
      if (aid.name.toLowerCase().contains(keyword) ||
          aid.description.toLowerCase().contains(keyword)) {
        _aidList.add(aid);
      }
    });
    notifyListeners();
  }

  Future<void> initAIDList() async {
    // clear UI List
    _aidList.clear();
    // link db
    if (!Hive.isBoxOpen('aidBox')) {
      _aidBox = await Hive.openBox<AID>('aidBox');
    } else {
      _aidBox = Hive.box<AID>('aidBox');
    }
    // full UI List
    _aidList.addAll(_aidBox.values);
    notifyListeners();
  }
}
