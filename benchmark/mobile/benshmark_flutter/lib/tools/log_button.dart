import 'package:flutter/material.dart';
import 'package:benshmark/pages/register_page.dart';
import 'package:benshmark/pages/login_page.dart';
//import 'package:flutter_application_1/pages/home_page.dart';

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
      width: widget.width,
      height: widget.height,
      child: FloatingActionButton(
        heroTag: "Register",
        onPressed: () => Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => const RegisterPage(),
          ),
        ),
        tooltip: 'Switch Page',
        backgroundColor: const Color(0xff9daaff),
        child: Text(style: const TextStyle(fontWeight: FontWeight.w200), widget.title),
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
        onPressed: () => Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => const LoginPage(),
          ),
        ),
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

// class LoginButton extends StatefulWidget {
//   const LoginButton(
//       {super.key,
//       required this.width,
//       required this.height,
//       required this.title});

//   final String title;
//   final double width;
//   final double height;

//   @override
//   State<LoginButton> createState() => _LoginButton();
// }

// class _LoginButton extends State<LoginButton> {
//   @override
//   Widget build(BuildContext context) {
//     return SizedBox(
//       width: widget.width,
//       height: widget.height,
//       child: FloatingActionButton(
//         heroTag: "Login",
//         onPressed: () => Navigator.push(
//           context,
//           MaterialPageRoute(
//             builder: (context) => MyHomePage(title: widget.title),
//           ),
//         ),
//         tooltip: 'Switch Page',
//         backgroundColor: const Color(0xffc1cbff),
//         child: Text(
//           widget.title,
//           style: const TextStyle(
//               fontWeight: FontWeight.bold,
//               color: Color.fromARGB(255, 58, 53, 53)),
//         ),
//       ),
//     );
//   }
// }

// class ProfileButton extends StatefulWidget {
//   const ProfileButton(
//       {super.key,
//       required this.width,
//       required this.height,
//       required this.title});

//   final String title;
//   final double width;
//   final double height;

//   @override
//   State<ProfileButton> createState() => _ProfileButton();
// }

// class _ProfileButton extends State<ProfileButton> {
//   @override
//   Widget build(BuildContext context) {
//     return SizedBox(
//       width: widget.width,
//       height: widget.height,
//       child: FloatingActionButton(
//         heroTag: "Profile",
//         onPressed: () => Navigator.push(
//           context,
//           MaterialPageRoute(
//             builder: (context) => const ProfilePage(),
//           ),
//         ),
//         tooltip: 'Switch Page',
//         backgroundColor: const Color(0xffc1cbff),
//         child: Text(
//           widget.title,
//           style: const TextStyle(
//               fontWeight: FontWeight.bold,
//               color: Color.fromARGB(255, 58, 53, 53)),
//         ),
//       ),
//     );
//   }
// }
