// ignore_for_file: avoid_print

import 'package:flutter/material.dart';
import 'package:benshmark/pages/home_page.dart';
import 'package:benshmark/tools/log_button.dart';
import 'package:benshmark/tools/space.dart';
import 'package:benshmark/pages/login_page.dart';
//import 'package:http/http.dart' as http;

class UserOuput extends StatefulWidget {
  UserOuput(
      {super.key,
      required this.title,
      required this.icon,
      required this.obscureText, required this.u});

  final String title;
  final bool obscureText;
  final IconData icon;
  final String u;
  final nameController = TextEditingController();
  final secondController = TextEditingController();

  @override
  State<UserOuput> createState() => _UserOuput();
}

class _UserOuput extends State<UserOuput> {
  String? _token;

  @override
  void dispose() {
    widget.nameController.dispose();
    widget.secondController.dispose();
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

  void switchPage() {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => MyHomePage(
            title: widget.secondController.text, headers: createHeader()),
      ),
    );
  }

  void badPassword() {
    Navigator.pop(context);
    Navigator.push(
        context, MaterialPageRoute(builder: (context) => const LoginPage()));
  }

  // Future<T?> _errorMessage<T>(String message) async {
  //   return showDialog(
  //     context: context,
  //     builder: (context) {
  //       return Center(
  //         child: Text(
  //           message,
  //           style: const TextStyle(
  //             fontSize: 30,
  //             fontWeight: FontWeight.bold,
  //             color: Colors.red,
  //           ),
  //         ),
  //       );
  //     },
  //   );
  // }

  // Future<void> _makeRequest(String a, String b, String u) async {
  //   final String body = "{ \"email\": \"$a\", \"password\": \"$b\" }";
  //   // print("uuuuuuu = ${u}");
  //   final Uri uri = Uri.http("127.0.0.1:8000", u);
  //   late final http.Response rep;
  //   late Map<String, dynamic> content;
  //   late String? str;

  //   try {
  //     rep = await http.post(uri, body: body);
  //   } catch (e) {
  //     print("error in post req");
  //     return _errorMessage("$e");
  //   }
  //   print(rep.body);
  //   content = jsonDecode(rep.body) as Map<String, dynamic>;
  //   switch ((rep.statusCode / 100) as int) {
  //     case 2:
  //       print("success");
  //       str = content['token']?.toString();
  //       if (str != null) {
  //         _token = str;
  //         switchPage();
  //       }
  //       break;
  //     case 4:
  //       str = content['message']?.toString();
  //       if (str != null) {
  //         _errorMessage(str);
  //       }
  //       break;
  //     default:
  //       break;
  //   }
  //}

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Center(
          child: FloatingActionButton(
            backgroundColor: Colors.green,
            onPressed: () {
              switchPage();
              // _makeRequest(
              //   widget.nameController.text,
              //   widget.secondController.text, widget.u
              // );
            },
            tooltip: "Show me",
            child: const Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SizedBox(
                  child: DecoratedBox(
                    decoration: BoxDecoration(color: Colors.green),
                    child: Text(
                      "Valider",
                      style: TextStyle(
                        color: Colors.black,
                        fontWeight: FontWeight.w200,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
        Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            children: [
              const Space(height: 100),
              UserBox(
                nameController: widget.nameController,
                icon: Icons.email,
                obscureText: false,
                title: "email:",
              ),
              const Space(height: 10),
              UserBox(
                nameController: widget.secondController,
                icon: Icons.password,
                obscureText: true,
                title: "password:",
              ),
              const Space(height: 30),
              const RegisterButton(
                width: 100,
                height: 30,
                title: "not register ?",
              )
            ],
          ),
        ),
      ],
    );
  }
}
