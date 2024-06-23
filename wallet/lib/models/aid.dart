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
}

class AIDListModel extends ChangeNotifier {
  final List<AID> _aidList = [];
  late Box<AID> _aidBox;

  bool isInitialized = false;

  // getter all values in db
  List<AID> get aidList => _aidBox.values.toList();

  int get aidCount => _aidList.length;

  Future<void> addAID(AID aid) async {
    // add to UI List
    _aidList.add(aid);
    // add to db
    await _aidBox.put(aid.aid, aid);
    notifyListeners();
  }

  Future<void> removeAID(String localId) async {
    // remove from UI List
    _aidList.removeWhere((aid) => aid.aid == localId);
    // remove from db
    await _aidBox.delete(localId);
    notifyListeners();
  }

  Future<void> updateAID(AID updatedAID) async {
    final index = _aidList.indexWhere((aid) => aid.aid == updatedAID.aid);
    if (index != -1) {
      // update UI List
      _aidList[index] = updatedAID;
      // update db
      await _aidBox.put(updatedAID.aid, updatedAID);
      notifyListeners();
    }
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
    _aidBox.toMap().forEach((key, value) {
      if (value.name.contains(keyword) || value.description.contains(keyword)) {
        _aidList.add(value);
      }
    });
    notifyListeners();
  }

  Future<void> initAIDList() async {
    if (isInitialized) {
      return;
    }
    // clear UI List
    _aidList.clear();
    // link db
    if (!Hive.isBoxOpen('aidBox')) {
      _aidBox = await Hive.openBox<AID>('aidBox');
    } else {
      _aidBox = Hive.box<AID>('aidBox');
    }
    // full UI List
    _aidBox.toMap().forEach((key, value) {
      _aidList.add(value);
    });
    isInitialized = true;
    notifyListeners();
  }
}
