import 'package:flutter/material.dart';
import 'package:benshmark/pages/home_page.dart';
import 'package:benshmark/tools/log_button.dart';
import 'package:benshmark/tools/space.dart';
import 'package:benshmark/pages/login_page.dart';
//import 'package:http/http.dart' as http;
import 'package:benshmark/tools/screen_scale.dart';
//import 'dart:convert';

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
  final name = TextEditingController();
  final surname = TextEditingController();
  final birth = TextEditingController();
  final gender = TextEditingController();

  @override
  State<UserRegister> createState() => _UserRegister();
}

class _UserRegister extends State<UserRegister> {
  String? _token;

  @override
  void dispose() {
    widget.email.dispose();
    widget.password.dispose();
    widget.name.dispose();
    widget.surname.dispose();
    widget.birth.dispose();
    widget.gender.dispose();
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

  void switchPage() {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) =>
            MyHomePage(title: widget.email.text, headers: createHeader()),
      ),
    );
  }

  void badPassword() {
    Navigator.pop(context);
    Navigator.push(
        context, MaterialPageRoute(builder: (context) => const LoginPage()));
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

  // Future<void> _makeRequest(String email, String password, String name,
  //     String surname, String birthdate, String gender, String u) async {
  //   final String body =
  //       "{\"email\": \"$email\", \"password\": \"$password\", \"name\": \"$name\", \"surname\": \"$surname\", \"date\": \"$birthdate\", \"gender\": \"$gender\"}";
  //   //print("uuuuuuu = ${u}");
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
  // }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            children: [
              UserBox(
                nameController: widget.email,
                icon: Icons.email,
                obscureText: false,
                title: "email:",
              ),
              Space(height: screenScale(context, 0.05).height),
              UserBox(
                nameController: widget.password,
                icon: Icons.password,
                obscureText: true,
                title: "password:",
              ),
              Space(height: screenScale(context, 0.05).height),
              UserBox(
                nameController: widget.name,
                icon: Icons.text_fields_outlined,
                obscureText: false,
                title: "name:",
              ),
              Space(height: screenScale(context, 0.05).height),
              UserBox(
                nameController: widget.surname,
                icon: Icons.text_fields_outlined,
                obscureText: false,
                title: "surname:",
              ),
              Space(height: screenScale(context, 0.05).height),
              UserBox(
                nameController: widget.birth,
                icon: Icons.text_fields_outlined,
                obscureText: false,
                title: "date:",
              ),
              Space(height: screenScale(context, 0.05).height),
              UserBox(
                nameController: widget.gender,
                icon: Icons.text_fields_outlined,
                obscureText: false,
                title: "gender:",
              ),
              Space(height: screenScale(context, 0.05).height),
              Center(
                child: FloatingActionButton(
                  backgroundColor: Colors.green,
                  onPressed: () {
                    // _makeRequest(
                    //   widget.email.text,
                    //   widget.password.text,
                    //   widget.name.text,
                    //   widget.surname.text,
                    //   widget.birth.text,
                    //   widget.gender.text,
                    //   widget.u,
                    // );
                    switchPage();
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
              Space(height: screenScale(context, 0.03).height),
              RegisterButton(
                width: screenScale(context, 0.07).width,
                height: screenScale(context, 0.03).height,
                title: "not register ?",
              )
            ],
          ),
        ),
      ],
    );
  }
}
