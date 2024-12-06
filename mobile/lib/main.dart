import 'package:flutter/material.dart';
import 'package:mobile/home_page.dart';
import 'package:mobile/pages/login_page.dart';
import 'package:mobile/pages/register_page.dart';
import 'package:mobile/pages/user_register.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/developers.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  MyApp({super.key});

  final GoRouter _router = GoRouter(
    initialLocation: '/login',
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => RegisterPage(),
      ),
      GoRoute(
        path: '/developers',
        builder: (context, state) {
          return const DevelopersPage();
        },
      ),
      GoRoute(
        path: '/aboutus',
        builder: (context, state) {
          return const DevelopersPage();
        },
      ),
      GoRoute(
        path: '/contact',
        builder: (context, state) {
          return const DevelopersPage();
        },
      ),
    ],
  );

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      title: 'AREA',
      routerConfig: _router,
    );
  }
}
