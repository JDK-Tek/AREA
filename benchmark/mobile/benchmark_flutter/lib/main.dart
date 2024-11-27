import 'package:flutter/material.dart';
import 'package:benshmark/tools/screen_scale.dart';
import 'package:benshmark/pages/login_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'For Benshmark',
      theme: ThemeData(
        appBarTheme: AppBarTheme(
          titleTextStyle: TextStyle(
            fontSize: (screenScale(context, .035).width * .8).clamp(20, 30),
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
          backgroundColor: Colors.white,
        ),
        //colorScheme: Colors.white,
        primaryColor: Colors.white,
        useMaterial3: true,
      ),
      home: const LoginPage(),
    );
  }
}
