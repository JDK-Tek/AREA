import 'package:flutter/material.dart';

class UserState extends ChangeNotifier {
  String? _token;

  String? get token => _token;

  void setToken(String token) {
    _token = token;
    notifyListeners();
  }
}

class IPState extends ChangeNotifier {
  String _ip = "api.area.jepgo.root.sx";

  String get ip => _ip;

  void setIP(String token) {
    _ip = ip;
    notifyListeners();
  }
}