import 'package:flutter/material.dart';
import 'package:mobile/tools/space.dart';
import 'package:mobile/pages/login_page.dart';
import 'package:go_router/go_router.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class UserRegister extends StatefulWidget {
  UserRegister(
      {super.key,
      required this.title,
      required this.icon,
      required this.obscureText,
      required this.u});

  final String title;
  final bool obscureText;
  final IconData icon;
  final String u;
  final email = TextEditingController();
  final password = TextEditingController();

  @override
  State<UserRegister> createState() => _UserRegister();
}

class _UserRegister extends State<UserRegister> {
  final FocusNode emailFocusNode = FocusNode();
  final FocusNode passwordFocusNode = FocusNode();

  String? _token;

  @override
  void dispose() {
    emailFocusNode.dispose();
    passwordFocusNode.dispose();
    widget.email.dispose();
    widget.password.dispose();
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
    final String body = "{ \"email\": \"$a\", \"password\": \"$b\" }";
    // print("uuuuuuu = ${u}");
    final Uri uri = Uri.http("localhost:1234", u);
    late final http.Response rep;
    late Map<String, dynamic> content;
    late String? str;

    try {
      rep = await http.post(uri, body: body);
    } catch (e) {
      return _errorMessage("$e");
    }
    content = jsonDecode(rep.body) as Map<String, dynamic>;
    switch ((rep.statusCode / 100) as int) {
      case 2:
        str = content['token']?.toString();
        if (str != null) {
          _token = str;
          if (mounted) {
            context.go("/");
          }
        } else {
          _errorMessage("Enter a valid email and password !");
        }
        break;
      case 4:
        str = content['message']?.toString();
        if (str != null) {
          _errorMessage(str);
        }
        break;
      case 5:
        _errorMessage("Enter a valid email and password !");
      default:
        break;
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
                  //color: const Color(0xff222222),
                  decoration: BoxDecoration(
                    color: const Color(0xff222222),
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Column(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      //crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        //const Space(height: 15),
                        const Text(
                          "REGISTER",
                          style: TextStyle(color: Colors.white, fontSize: 20),
                        ),
                        const Text(
                            "Welcome, we are delighted to have you among us",
                            style: TextStyle(
                                color: Color(0xff8c52ff), fontSize: 12)),
                        //const Space(height: 50),
                        SizedBox(
                          height: 50,
                          width: 300,
                          child: UserBox(
                            nameController: widget.email,
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
                            nameController: widget.password,
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
                              _makeRequest(widget.email.text,
                                  widget.password.text, "api/register");
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
