import 'package:flutter/material.dart';
import 'package:mobile/pages/login_page.dart';
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
