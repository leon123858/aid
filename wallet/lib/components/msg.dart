import 'package:flutter/material.dart';
import 'package:getwidget/components/toast/gf_toast.dart';
import 'package:getwidget/position/gf_toast_position.dart';

void showErrorToast(BuildContext context, String message) {
  GFToast.showToast(
    message,
    context,
    toastPosition: GFToastPosition.BOTTOM,
    textStyle: const TextStyle(fontSize: 16, color: Colors.white),
    backgroundColor: Colors.red,
    trailing: const Icon(
      Icons.error,
      color: Colors.white,
    ),
  );
}

void showSuccessToast(BuildContext context, String message) {
  GFToast.showToast(
    message,
    context,
    toastPosition: GFToastPosition.BOTTOM,
    textStyle: const TextStyle(fontSize: 16, color: Colors.white),
    backgroundColor: Colors.green,
    trailing: const Icon(
      Icons.check,
      color: Colors.white,
    ),
  );
}
