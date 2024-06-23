// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'aid.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class AIDAdapter extends TypeAdapter<AID> {
  @override
  final int typeId = 0;

  @override
  AID read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return AID(
      aid: fields[0] as String,
      name: fields[1] as String,
      description: fields[2] as String,
      publicKey: fields[3] as String,
      privateKey: fields[4] as String,
    );
  }

  @override
  void write(BinaryWriter writer, AID obj) {
    writer
      ..writeByte(5)
      ..writeByte(0)
      ..write(obj.aid)
      ..writeByte(1)
      ..write(obj.name)
      ..writeByte(2)
      ..write(obj.description)
      ..writeByte(3)
      ..write(obj.publicKey)
      ..writeByte(4)
      ..write(obj.privateKey);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is AIDAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
