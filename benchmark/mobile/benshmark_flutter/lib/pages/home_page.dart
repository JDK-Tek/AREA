import 'package:flutter/material.dart';
import 'package:benshmark/tools/screen_scale.dart';
//import 'package:benshmark/tools/space.dart';

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
  late String button;
  List<String> _list = ["test1", "test2"];
  String i = "";
  void test(int index) {
    setState(() {});
  }

  @override
  void initState() {
    button = _list[0];
    super.initState();
  }

  PreferredSizeWidget myOwnAppbar(BuildContext context) {
    return PreferredSize(
        preferredSize: screenScale(context, .1), child: const Text("coucou"));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: myOwnAppbar(context),
      body: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              //const SpaceW(width: 10,),
              Center(
                child: Container(
                  width: 100,
                  height: 100,
                  decoration: const BoxDecoration(
                    image: DecorationImage(
                      image: NetworkImage(
                          "https://blog.realogs.in/content/images/2021/09/kotlin.png"),
                      fit: BoxFit.contain,
                    ),
                  ),
                ),
              ),
              //const SpaceW(width: 250,),
              Center(
                child: Container(
                  width: 100,
                  height: 100,
                  decoration: const BoxDecoration(
                    image: DecorationImage(
                      image: NetworkImage(
                          "https://i2.wp.com/softwareengineeringdaily.com/wp-content/uploads/2018/10/flutter.jpg?fit=1570,1500&ssl=1"),
                      fit: BoxFit.contain,
                    ),
                  ),
                ),
              ),
            ],
          ),
          DropdownButton<String>(
            value: button,
            icon: const Icon(Icons.arrow_downward),
            elevation: 16,
            onChanged: (val) {
              setState(() {
                button = val!;
              });
            },
            items: _list.map((val) {
              return DropdownMenuItem<String>(
                value: val,
                child: Text(val),
              );
            }).toList(),
          ),
        ],
      ),
    );
  }
}
