import 'package:flutter/material.dart';
import 'package:benshmark/tools/screen_scale.dart';
import 'package:benshmark/tools/log_button.dart';
import 'package:benshmark/pages/user_register.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

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
      body: SizedBox(
        height: screenScale(context, 0.9).height,
        child: Column(
          //mainAxisAlignment: MainAxisAlignment.center,
          children: [
            UserRegister(
                title: "email",
                icon: Icons.email,
                obscureText: false,
                u: "/api/register"),
            LogoutButton(
                width: screenScale(context, 0.07).width, height: screenScale(context, 0.03).height, title: "Already log ?")
          ],
        ),
      ),
    );
  }
}
