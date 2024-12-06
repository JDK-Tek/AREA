// ignore_for_file: avoid_print

import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:mobile/tools/log_button.dart';
import 'package:mobile/tools/space.dart';
import 'package:mobile/pages/login_page.dart';
import 'package:go_router/go_router.dart';
import 'package:http/http.dart' as http;

class UserOuput extends StatefulWidget {
  const UserOuput(
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
  State<UserOuput> createState() => _UserOuput();
}

class _UserOuput extends State<UserOuput> {
  late TextEditingController nameController;
  late TextEditingController secondController;
  late FocusNode emailFocusNode;
  late FocusNode passwordFocusNode;
  String? _token;

  @override
  void initState() {
    // Initialisation des contrôleurs dans le State
    emailFocusNode = FocusNode();
    passwordFocusNode = FocusNode();
    nameController = TextEditingController();
    secondController = TextEditingController();
    super.initState();
  }

  @override
  void dispose() {
    nameController.dispose();
    secondController.dispose();
    emailFocusNode.dispose();
    passwordFocusNode.dispose();
    super.dispose();
  }

  Map<String, String> createHeader() {
    // need to be cancel
    _token ?? "";

    // if (_token == null) {
    //   //throw Exception("Error: missing Token");
    // }
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
    //final Uri uri = Uri.http("172.20.10.3:1234", u);
    late final http.Response rep;
    late Map<String, dynamic> content;
    late String? str;

    try {
      rep = await http.post(uri, body: body);
    } catch (e) {
      print("error in post req");
      print("$e");
      return _errorMessage("$e");
    }
    print(rep.body);
    print(rep.statusCode);
    content = jsonDecode(rep.body) as Map<String, dynamic>;
    switch ((rep.statusCode / 100) as int) {
      case 2:
        print("success");
        str = content['token']?.toString();
        if (str != null) {
          _token = str;
          context.go("/");
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
                          "LOGIN",
                          style: TextStyle(color: Colors.white, fontSize: 20),
                        ),
                        const Text("Nice to see you again",
                            style: TextStyle(
                                color: Color(0xff8c52ff), fontSize: 15)),
                        //const Space(height: 50),
                        SizedBox(
                          height: 50,
                          width: 300,
                          child: UserBox(
                            nameController: nameController,
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
                            nameController: secondController,
                            icon: Icons.password,
                            obscureText: true,
                            title: "password:",
                          ),
                        ),
                        const Center(
                          child: Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                    style: TextStyle(
                                        color: Color(0xffa6a6a6), fontSize: 12),
                                    "You already have an account ?"),
                                RegisterButton(
                                    width: 50, height: 10, title: "register ?")
                              ]),
                        ),
                        FloatingActionButton.extended(
                            label: const Text(
                                style: TextStyle(color: Colors.black), "Login"),
                            backgroundColor: Colors.white,
                            extendedPadding: const EdgeInsets.symmetric(
                                vertical: 8.0, horizontal: 32.0),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(15),
                            ),
                            onPressed: () {
                              _makeRequest(nameController.text,
                                  secondController.text, "api/login");
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
