import 'package:flutter/material.dart';
import 'package:mobile/tools/screen_scale.dart';
import 'package:mobile/pages/login_page.dart';
import 'package:mobile/pages/user_register.dart';

class RegisterPage extends StatefulWidget {
  RegisterPage({super.key});

  final nameController = TextEditingController();
  final secondController = TextEditingController();

  @override
  State<RegisterPage> createState() => _RegisterPage();
}

class _RegisterPage extends State<RegisterPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: const Color(0xfff5f6fA),
        appBar: AppBar(
          toolbarHeight: screenScale(context, 0.1).height,
          elevation: 1,
          backgroundColor: const Color(0xfffefffe),
          title: const Text("Register Page"),
        ),
        body: SingleChildScrollView(
          child: Stack(
            children: [
              const AnimatedBackground(),
              SizedBox(
                height: screenScale(context, 0.9).height,
                child: Column(
                  //mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    UserRegister(
                        title: "email",
                        icon: Icons.email,
                        obscureText: false,
                        u: "/api/register"),
                  ],
                ),
              ),
            ],
          ),
        ));
  }
}
