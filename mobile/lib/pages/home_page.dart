import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:area/pages/servicepage.dart';
import 'package:area/pages/appletspage.dart';

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
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const HeaderSection(),
            AppletSection(),
            const ServiceSection()
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
    double screenHeight = MediaQuery.of(context).size.height;
    double screenWidth = MediaQuery.of(context).size.width;

    return Container(
      padding: const EdgeInsets.all(0.0),
      margin: const EdgeInsets.all(0.0),
      height: screenWidth < screenHeight
          ? screenHeight * 0.30
          : screenHeight * 0.80,
      width: screenWidth,
      decoration:
          BoxDecoration(gradient: LinearGradient(colors: predefinedColors, begin: Alignment.topRight, end: Alignment.bottomLeft)),
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
                  height: screenHeight * 0.08,
                  width: screenWidth < screenHeight
                      ? screenWidth * 0.2
                      : screenWidth * 0.1,
                  fit: BoxFit.contain,
                ),
              ),
              ElevatedButton(
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
                  style: GoogleFonts.nunito(
                    fontWeight: FontWeight.w900,
                    fontSize: screenWidth < screenHeight
                        ? screenWidth * 0.049
                        : screenWidth * 0.025,
                    color: Colors.black,
                  ),
                ),
              ),
            ],
          ),
          SizedBox(
            height: screenWidth < screenHeight
                ? screenHeight * 0.095
                : screenHeight * 0.3,
            width: screenWidth * 0.8,
            child: Text(
              "AUTOMATION FOR BUSINESS AND HOME",
              textAlign: TextAlign.center,
              style: GoogleFonts.nunito(
                fontSize: screenWidth < screenHeight
                    ? screenWidth * 0.07
                    : screenWidth * 0.05,
                fontWeight: FontWeight.w900,
                color: Colors.white,
              ),
            ),
          ),
          Text(
            "Save time and get more done",
            textAlign: TextAlign.center,
            style: GoogleFonts.nunito(
              color: const Color.fromARGB(255, 186, 151, 255),
              fontSize: screenWidth < screenHeight
                  ? screenWidth * 0.04
                  : screenWidth * 0.03,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(
            height: screenWidth < screenHeight
                ? screenHeight * 0.01
                : screenHeight * 0.03,
          ),
          SizedBox(
            height: screenWidth < screenHeight
                ? screenHeight * 0.05
                : screenHeight * 0.08,
            width: screenWidth < screenHeight
                ? screenWidth * 0.38
                : screenWidth * 0.17,
            child: ElevatedButton(
              onPressed: () {
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
                "Get started →",
                textAlign: TextAlign.center,
                style: GoogleFonts.nunito(
                  fontWeight: FontWeight.w900,
                  fontSize: screenWidth < screenHeight
                      ? screenWidth * 0.049
                      : screenWidth * 0.025,
                  color: Colors.black,
                ),
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
    double screenHeight = MediaQuery.of(context).size.height;
    double screenWidth = MediaQuery.of(context).size.width;

    return Container(
      padding: const EdgeInsets.all(0.0),
      margin: const EdgeInsets.all(0.0),
      width: screenWidth,
      decoration:
          BoxDecoration(gradient: LinearGradient(colors: predefinedColors, begin: Alignment.topRight, end: Alignment.bottomLeft)),
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
                  height: screenHeight * 0.08,
                  width: screenWidth < screenHeight
                      ? screenWidth * 0.2
                      : screenWidth * 0.1,
                  fit: BoxFit.contain,
                ),
              ),
              ElevatedButton(
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
                  style: GoogleFonts.nunito(
                    fontWeight: FontWeight.w900,
                    fontSize: screenWidth < screenHeight
                        ? screenWidth * 0.049
                        : screenWidth * 0.025,
                    color: Colors.black,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
