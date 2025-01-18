import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'package:area/tools/providers.dart';

class LoginButton extends StatefulWidget {
  const LoginButton({super.key});

  @override
  State<LoginButton> createState() => LoginButtonState();
}

class LoginButtonState extends State<LoginButton> {
  @override
  Widget build(BuildContext context) {
    final token = Provider.of<UserState>(context).token;
    if (token != null) {
      return ElevatedButton(
        onPressed: () {
          Provider.of<UserState>(context, listen: false).unsetToken(null);
          context.go("/");
        },
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.all(8),
          backgroundColor: Colors.white,
          foregroundColor: Colors.black,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
        ),
        child: Text(
          "Logout",
          textAlign: TextAlign.center,
          style: TextStyle(
              fontWeight: FontWeight.w900,
              fontSize: MediaQuery.of(context).size.width <
                      MediaQuery.of(context).size.height
                  ? MediaQuery.of(context).size.width * 0.049
                  : MediaQuery.of(context).size.width * 0.025,
              color: Colors.black,
              fontFamily: 'Nunito-Bold'),
        ),
      );
    } else {
      return ElevatedButton(
        onPressed: () {
          context.go("/login");
        },
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.all(8),
          backgroundColor: Colors.white,
          foregroundColor: Colors.black,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
        ),
        child: Text(
          "Login",
          textAlign: TextAlign.center,
          style: TextStyle(
              fontWeight: FontWeight.w900,
              fontSize: MediaQuery.of(context).size.width <
                      MediaQuery.of(context).size.height
                  ? MediaQuery.of(context).size.width * 0.049
                  : MediaQuery.of(context).size.width * 0.025,
              color: Colors.black,
              fontFamily: 'Nunito-Bold'),
        ),
      );
    }
  }
}
