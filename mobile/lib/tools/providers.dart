import 'package:flutter/material.dart';

class UserState extends ChangeNotifier {
  String? _token;

  String? get token => _token;

  void unsetToken(Null n) {
    _token = n;
    notifyListeners();
  }

  void setToken(String token) {
    _token = token;
    notifyListeners();
  }
}

class IPState extends ChangeNotifier {
  String _ip = "api.area.jepgo.root.sx";

  String get ip => _ip;

  void setIP(String ip) {
    _ip = ip;
    notifyListeners();
  }
}