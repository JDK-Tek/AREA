import 'dart:convert';
import 'package:area/tools/providers.dart';
import 'package:flutter/material.dart';
import 'package:area/tools/log_button.dart';
import 'package:area/tools/space.dart';
import 'package:area/pages/login_page.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'package:http/http.dart' as https;

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
  String _token = "";

  @override
  void initState() {
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
    Map<String, String> headers = {
      "token": _token,
    };
    return headers;
  }

  void switchPage() {
  }

  void badPassword() {
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

  Future<bool> _makeRequest(String a, String b, String u) async {
    final String body = "{ \"email\": \"$a\", \"password\": \"$b\" }";
    final Uri uri =
        Uri.https(Provider.of<IPState>(context, listen: false).ip, u);
    late final https.Response rep;
    late Map<String, dynamic> content;
    late String? str;

    try {
      rep = await https.post(uri, body: body);
    } catch (e) {
      _errorMessage("$e");
      return false;
    }
    if (rep.statusCode >= 500) {
      _errorMessage(rep.body);
      return false;
    }
    content = jsonDecode(rep.body) as Map<String, dynamic>;
    str = content['token']?.toString();
    if (str != null) {
      _token = str;
      if (mounted) {
        Provider.of<UserState>(context, listen: false).setToken(_token);
      }
    }
    return true;
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
              Space(height: MediaQuery.of(context).size.height * 0.06),
              Center(
                child: Container(
                  height: MediaQuery.of(context).size.height < 600
                      ? MediaQuery.of(context).size.height
                      : MediaQuery.of(context).size.height * 0.5,
                  width: MediaQuery.of(context).size.width * 0.85,
                  decoration: BoxDecoration(
                    color: const Color(0xff222222),
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Column(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        const Text(
                          "LOGIN",
                          style: TextStyle(color: Colors.white, fontSize: 20),
                        ),
                        const Text("Nice to see you again",
                            style: TextStyle(
                                color: Color(0xff8c52ff), fontSize: 15)),
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
                                      secondController.text, "api/login")
                                  .then((key) {
                                if (!context.mounted || key == false) return;
                                Provider.of<UserState>(context, listen: false)
                                    .setToken(_token);
                                context.go("/");
                              });
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
