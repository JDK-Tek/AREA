import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:area/pages/home_page.dart';
import 'package:flutter/cupertino.dart';

class DevelopersPage extends StatefulWidget {
  const DevelopersPage({super.key});

  @override
  State<DevelopersPage> createState() => DevelopersPageState();
}

class DevelopersPageState extends State<DevelopersPage> {
  int currentPageIndex = 3;

  @override
  Widget build(BuildContext context) {
    final List<String> dest = [
      "/applets",
      "/create",
      "/services",
      "/developers"
    ];
    return SafeArea(
        child: Scaffold(
      bottomNavigationBar: NavigationBar(
        backgroundColor: Colors.black,
        indicatorColor: Colors.grey,
        selectedIndex: 3,
        onDestinationSelected: (int index) {
          setState(() {
            currentPageIndex = index;
            context.go(dest[index]);
          });
        },
        destinations: const [
          NavigationDestination(
              icon: Icon(Icons.folder, color: Colors.white), label: 'Applets'),
          NavigationDestination(
              icon: Icon(Icons.add_circle_outline, color: Colors.white),
              label: 'Create'),
          NavigationDestination(
              icon: Icon(Icons.cloud, color: Colors.white), label: 'Services'),
          NavigationDestination(
              icon: Icon(CupertinoIcons.ellipsis, color: Colors.white),
              label: 'Developers'),
        ],
      ),
      backgroundColor: Colors.white,
      body: const SingleChildScrollView( child:
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              MiniHeaderSection(),
              Text(
                "Elise Esteban Greg Paul John",
                textAlign: TextAlign.center,
              )
            ],
          ),
      ),
    ));
  }
}
