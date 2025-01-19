import 'package:area/tools/login_button.dart';
import 'package:area/tools/providers.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:area/pages/appletspage.dart';
import 'package:provider/provider.dart';

List<Color> predefinedColors = [
  const Color(0xff222222),
  const Color(0xff410cab),
  const Color(0xff5e17eb),
];

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => HomePageState();
}

class HomePageState extends State<HomePage> {
  int currentPageIndex = 0;

  @override
  Widget build(BuildContext context) {
    final List<String> dest = ["/applets", "/create", "/services", "/plus"];

    return SafeArea(
        child: Scaffold(
      bottomNavigationBar: NavigationBar(
        backgroundColor: Colors.black,
        indicatorColor: Colors.grey,
        selectedIndex: 0,
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
      body: const SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            HeaderSection(),
            AppletSection(),
          ],
        ),
      ),
    ));
  }
}

class HeaderSection extends StatelessWidget {
  const HeaderSection({super.key});

  @override
  Widget build(BuildContext context) {
    var token = Provider.of<UserState>(context, listen: false).token;
    return Container(
      padding: const EdgeInsets.all(0.0),
      margin: const EdgeInsets.all(0.0),
      height:
          MediaQuery.of(context).size.width < MediaQuery.of(context).size.height
              ? MediaQuery.of(context).size.height * 0.35
              : MediaQuery.of(context).size.height * 0.80,
      width: MediaQuery.of(context).size.width,
      decoration: BoxDecoration(
          gradient: LinearGradient(
              colors: predefinedColors,
              begin: Alignment.topRight,
              end: Alignment.bottomLeft)),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              ElevatedButton(
                onPressed: () {
                  context.go('/');
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.transparent,
                  shadowColor: Colors.transparent,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(20),
                  ),
                ),
                child: Image.asset(
                  'assets/fullLogo.png',
                  height: MediaQuery.of(context).size.height * 0.08,
                  width: MediaQuery.of(context).size.width <
                          MediaQuery.of(context).size.height
                      ? MediaQuery.of(context).size.width * 0.2
                      : MediaQuery.of(context).size.width * 0.1,
                  fit: BoxFit.contain,
                ),
              ),
              const LoginButton(),
            ],
          ),
          SizedBox(
            height: MediaQuery.of(context).size.width <
                    MediaQuery.of(context).size.height
                ? MediaQuery.of(context).size.height * 0.095
                : MediaQuery.of(context).size.height * 0.3,
            width: MediaQuery.of(context).size.width * 0.8,
            child: Text(
              "AUTOMATION FOR BUSINESS AND HOME",
              textAlign: TextAlign.center,
              style: TextStyle(
                  fontSize: MediaQuery.of(context).size.width <
                          MediaQuery.of(context).size.height
                      ? MediaQuery.of(context).size.width * 0.055
                      : MediaQuery.of(context).size.width * 0.05,
                  fontWeight: FontWeight.w900,
                  color: Colors.white,
                  fontFamily: 'Nunito-Black'),
            ),
          ),
          SizedBox(
            height: MediaQuery.of(context).size.width <
                    MediaQuery.of(context).size.height
                ? MediaQuery.of(context).size.height * 0.01
                : 0,
            width: 0,
          ),
          Text(
            "Save time and get more done",
            textAlign: TextAlign.center,
            style: TextStyle(
                color: const Color.fromARGB(255, 186, 151, 255),
                fontSize: MediaQuery.of(context).size.width <
                        MediaQuery.of(context).size.height
                    ? MediaQuery.of(context).size.width * 0.04
                    : MediaQuery.of(context).size.width * 0.03,
                fontWeight: FontWeight.w900,
                fontFamily: 'Nunito-Bold'),
          ),
          SizedBox(
            height: MediaQuery.of(context).size.width <
                    MediaQuery.of(context).size.height
                ? MediaQuery.of(context).size.height * 0.01
                : MediaQuery.of(context).size.height * 0.03,
          ),
          SizedBox(
            height: MediaQuery.of(context).size.width <
                    MediaQuery.of(context).size.height
                ? MediaQuery.of(context).size.height * 0.05
                : MediaQuery.of(context).size.height * 0.08,
            width: MediaQuery.of(context).size.width <
                    MediaQuery.of(context).size.height
                ? MediaQuery.of(context).size.width * 0.38
                : MediaQuery.of(context).size.width * 0.2,
            child: ElevatedButton(
              onPressed: () {
                if (token != null) {
                }
                context.go("/register");
              },
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.all(0.0),
                backgroundColor: Colors.white,
                foregroundColor: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(20),
                ),
              ),
              child: Text(
                token != null ? "My Applets" : "Get started →",
                textAlign: TextAlign.center,
                style: TextStyle(
                    fontWeight: FontWeight.w900,
                    fontSize: MediaQuery.of(context).size.width <
                            MediaQuery.of(context).size.height
                        ? MediaQuery.of(context).size.width * 0.046
                        : MediaQuery.of(context).size.width * 0.025,
                    color: Colors.black,
                    fontFamily: 'Nunito-Bold'),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class MiniHeaderSection extends StatelessWidget {
  const MiniHeaderSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(0.0),
      margin: const EdgeInsets.all(0.0),
      width: MediaQuery.of(context).size.width,
      decoration: BoxDecoration(
          gradient: LinearGradient(
              colors: predefinedColors,
              begin: Alignment.topRight,
              end: Alignment.bottomLeft)),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              ElevatedButton(
                onPressed: () {
                  context.go('/');
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.transparent,
                  shadowColor: Colors.transparent,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(20),
                  ),
                ),
                child: Image.asset(
                  'assets/fullLogo.png',
                  height: MediaQuery.of(context).size.height * 0.08,
                  width: MediaQuery.of(context).size.width <
                          MediaQuery.of(context).size.height
                      ? MediaQuery.of(context).size.width * 0.2
                      : MediaQuery.of(context).size.width * 0.1,
                  fit: BoxFit.contain,
                ),
              ),
              const LoginButton(),
            ],
          ),
        ],
      ),
    );
  }
}
