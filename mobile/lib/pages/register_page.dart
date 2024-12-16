import 'package:flutter/material.dart';
import 'package:area/tools/screen_scale.dart';
import 'package:area/pages/login_page.dart';
import 'package:area/pages/user_register.dart';
import 'package:go_router/go_router.dart';

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
    return SafeArea(
        child: Scaffold(
            backgroundColor: const Color(0xff222222),
            body: SingleChildScrollView(
              child: Stack(
                children: [
                  const AnimatedBackground(),
                  Align(
                    alignment: Alignment.topLeft,
                    child: ElevatedButton(
                        style: ElevatedButton.styleFrom(
                            backgroundColor: const Color.fromARGB(0, 0, 0, 0),
                            foregroundColor: const Color.fromARGB(0, 0, 0, 0),
                            shadowColor: const Color.fromARGB(0, 0, 0, 0),),
                        onPressed: () {
                          context.go("/applets");
                        },
                        child: Icon(Icons.arrow_back,
                            color: Colors.white,
                            size: screenScale(context, 0.05).height)),
                  ),
                  SizedBox(
                    height: screenScale(context, 0.9).height,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
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
            )));
  }
}
