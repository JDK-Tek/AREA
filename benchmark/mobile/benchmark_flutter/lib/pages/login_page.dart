import 'package:flutter/material.dart';
import 'package:benshmark/tools/screen_scale.dart';
import 'package:benshmark/tools/user_ouput.dart';

class UserBox extends StatefulWidget {
  const UserBox(
      {super.key,
      required this.nameController,
      required this.icon,
      required this.obscureText,
      required this.title});

  final TextEditingController nameController;
  final bool obscureText;
  final IconData icon;
  final String title;

  @override
  State<UserBox> createState() => _UserBox();
}

class _UserBox extends State<UserBox> {
  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: widget.nameController,
      obscureText: widget.obscureText,
      decoration: InputDecoration(
          prefixIcon: Icon(widget.icon),
          suffixIcon: const Icon(Icons.clear),
          border: const OutlineInputBorder(),
          labelText: widget.title),
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPage();
}

class _LoginPage extends State<LoginPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: const Color(0xffe1e4ed),
        appBar: AppBar(
          toolbarHeight: screenScale(context, 0.1).height,
          elevation: 1,
          backgroundColor: const Color(0xfffefffe),
          title: const Text("Login Page"),
        ),
        body: Stack(children: [
          UserOuput(
              title: "email:",
              icon: Icons.email,
              obscureText: false,
              u: "/api/login"),
        ]));
  }
}
