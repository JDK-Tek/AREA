import 'package:flutter/material.dart';
import 'package:area/tools/space.dart';
import 'package:area/pages/login_page.dart';
import 'package:go_router/go_router.dart';
import 'package:http/http.dart' as https;
import 'dart:convert';
import 'package:area/tools/providers.dart';
import 'package:provider/provider.dart';

class UserRegister extends StatefulWidget {
  const UserRegister(
      {super.key,
      required this.title,
      required this.icon,
      required this.obscureText,
      required this.u});

  final String title;
  final bool obscureText;
  final IconData icon;
  final String u;

  @override
  State<UserRegister> createState() => _UserRegister();
}

class _UserRegister extends State<UserRegister> {
  final FocusNode emailFocusNode = FocusNode();
  final FocusNode passwordFocusNode = FocusNode();

  final email = TextEditingController();
  final password = TextEditingController();

  String? _token;

  @override
  void dispose() {
    emailFocusNode.dispose();
    passwordFocusNode.dispose();
    email.dispose();
    password.dispose();
    super.dispose();
  }

  Map<String, String> createHeader() {
    if (_token == null) {
      throw Exception("Error: missing Token");
    }
    Map<String, String> headers = {
      "token": _token ?? "",
    };
    return headers;
  }

  Future<T?> _errorMessage<T>(String message) async {
    return showDialog(
      context: context,
      builder: (context) {
        return Center(
          child: Text(
            message,
            style: const TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.bold,
              color: Colors.red,
            ),
          ),
        );
      },
    );
  }

  Future<void> _makeRequest(String a, String b, String u) async {
    final Map<String, String> body = {
      "email": a,
      "password": b,
    };
    final Uri uri =
        Uri.https(Provider.of<IPState>(context, listen: false).ip, u);
    late final https.Response rep;
    late Map<String, dynamic> content;

    try {
      rep = await https.post(
        uri,
        headers: {
          "Content-Type": "application/json",
        },
        body: jsonEncode(body),
      );
      switch (rep.statusCode) {
        case 200:
          content = jsonDecode(rep.body);
          _token = content['token']?.toString();
          if (_token != null && mounted) {
            Provider.of<UserState>(context, listen: false).setToken(_token!);
            context.go("/");
          } else {
            await _errorMessage("Token not received from server.");
          }
          break;
        case 400:
          await _errorMessage("400 Bad Request: ${rep.body}");
          break;
        case 500:
          await _errorMessage("Server error: ${rep.body}");
          break;
        default:
          await _errorMessage("Unexpected server response.");
      }
    } catch (e) {
      await _errorMessage("Failed to make request: $e");
    }
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              const Space(height: 120),
              Center(
                child: Container(
                  height: MediaQuery.of(context).size.height < 600
                      ? MediaQuery.of(context).size.height
                      : MediaQuery.of(context).size.height * 0.5,
                  width: MediaQuery.of(context).size.width * 0.75,
                  decoration: BoxDecoration(
                    color: const Color(0xff222222),
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Column(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        const Text(
                          "REGISTER",
                          style: TextStyle(color: Colors.white, fontSize: 20),
                        ),
                        const Text(
                            "Welcome, we are delighted to have you among us",
                            style: TextStyle(
                                color: Color(0xff8c52ff), fontSize: 12)),
                        SizedBox(
                          height: 50,
                          width: 300,
                          child: UserBox(
                            nameController: email,
                            icon: Icons.email,
                            obscureText: false,
                            title: "email:",
                          ),
                        ),
                        const Space(height: 10),
                        SizedBox(
                          height: 50,
                          width: 300,
                          child: UserBox(
                            nameController: password,
                            icon: Icons.password,
                            obscureText: true,
                            title: "password:",
                          ),
                        ),
                        Center(
                          child: Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                const Text(
                                    style: TextStyle(
                                        color: Color(0xffa6a6a6), fontSize: 12),
                                    "You already have an account ?"),
                                TextButton(
                                  onPressed: () {
                                    context.go("/login");
                                  },
                                  child: const Text(
                                      style: TextStyle(
                                          fontSize: 12,
                                          color: Colors.white,
                                          decoration: TextDecoration.underline,
                                          decorationThickness: 2.0,
                                          decorationColor: Colors.white),
                                      "Login ?"),
                                )
                              ]),
                        ),
                        FloatingActionButton.extended(
                            label: const Text(
                                style: TextStyle(color: Colors.black),
                                "Register"),
                            backgroundColor: Colors.white,
                            extendedPadding: const EdgeInsets.symmetric(
                                vertical: 8.0, horizontal: 32.0),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(15),
                            ),
                            onPressed: () {
                              _makeRequest(
                                  email.text, password.text, "api/register");
                            }),
                      ]),
                ),
              ),
              const Space(height: 30),
            ],
          ),
        ),
      ),
    );
  }
}
