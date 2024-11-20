import 'package:flutter/material.dart';
import 'package:benshmark/tools/screen_scale.dart';

const ColorScheme mainColorScheme = ColorScheme(
  brightness: Brightness.light,
  primary: Colors.white,
  onPrimary: Color.fromARGB(255, 50, 60, 72),
  secondary: Color(0xff9baaed),
  onSecondary: Colors.black,
  error: Color.fromARGB(255, 144, 18, 18),
  onError: Colors.white,
  surface: Color(0xfff1f1fb),
  onSurface: Color(0xff435062),
);

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title, required this.headers});

  final String title;
  final Map<String, String> headers;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  String i = "";
  void test(int index) {
    setState(() {});
  }

  @override
  void initState() {
    super.initState();
  }

  PreferredSizeWidget myOwnAppbar(BuildContext context) {
    return PreferredSize(
      preferredSize: screenScale(context, .1),
      child: const Text("coucou")
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: myOwnAppbar(context),
        body: Container(
          color: Colors.white,
        ),
    );
  }
}
