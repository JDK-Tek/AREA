import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class RegisterButton extends StatefulWidget {
  const RegisterButton(
      {super.key,
      required this.width,
      required this.height,
      required this.title});

  final String title;
  final double width;
  final double height;

  @override
  State<RegisterButton> createState() => _RegisterButton();
}

class _RegisterButton extends State<RegisterButton> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      child: TextButton(
        onPressed: () {
          context.go("/register");
        },
        child: Text(
            style: const TextStyle(
                fontSize: 12,
                color: Colors.white,
                decoration: TextDecoration.underline,
                decorationThickness: 2.0,
                decorationColor: Colors.white),
            widget.title),
      ),
    );
  }
}

class LogoutButton extends StatefulWidget {
  const LogoutButton(
      {super.key,
      required this.width,
      required this.height,
      required this.title});

  final String title;
  final double width;
  final double height;

  @override
  State<LogoutButton> createState() => _LogoutButton();
}

class _LogoutButton extends State<LogoutButton> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: widget.width,
      height: widget.height,
      child: FloatingActionButton(
        heroTag: "Logout",
        onPressed: () {
          context.go("/login");
        },
        tooltip: 'Switch Page',
        backgroundColor: const Color(0xff6175ff),
        child: Text(
          widget.title,
          style:
              const TextStyle(fontWeight: FontWeight.bold, color: Colors.black),
        ),
      ),
    );
  }
}