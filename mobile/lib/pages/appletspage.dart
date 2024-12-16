import 'package:flutter/material.dart';
import 'package:area/pages/home_page.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';

class AppletsPage extends StatefulWidget {
  const AppletsPage({super.key});
  @override
  State<AppletsPage> createState() => AppletPageState();
}

class AppletPageState extends State<AppletsPage> {
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
          children: [const HeaderSection(), AppletSection()],
        ),
      ),
    ));
  }
}

class Applet extends StatelessWidget {
  final String nameService;
  final IconData icon1;
  final IconData icon2;
  final String nameAREA;
  final String route;
  final VoidCallback press;
  final Color color;

  const Applet(
      {super.key,
      required this.icon1,
      required this.icon2,
      required this.nameService,
      required this.nameAREA,
      required this.route,
      required this.press,
      required this.color});

  @override
  Widget build(BuildContext context) {
    double screenHeight = MediaQuery.of(context).size.height;
    double screenWidth = MediaQuery.of(context).size.width;

    return Container(
        width: screenWidth < screenHeight
            ? screenWidth * 0.48
            : screenWidth * 0.38,
        height:
            screenWidth < screenHeight ? screenHeight * 0.25 : screenWidth * 0.3,
        margin: const EdgeInsets.only(top: 20),
        child: ElevatedButton(
            onPressed: () {
              press();
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: color,
              padding: const EdgeInsets.all(9),
              shape: const RoundedRectangleBorder(
                borderRadius: BorderRadius.all(Radius.elliptical(15, 15)),
              ),
            ),
            child: Stack(
              alignment: Alignment.bottomLeft,
              children: [
                Column(
                  children: [
                    Align(
                        alignment: Alignment.topLeft,
                        child: Row(
                          children: [
                            Icon(
                              icon1,
                              color: Colors.white,
                              size: screenWidth < screenHeight
                                  ? screenWidth * 0.1
                                  : screenWidth * 0.06,
                            ),
                            Icon(
                              icon2,
                              color: Colors.white,
                              size: screenWidth < screenHeight
                                  ? screenWidth * 0.1
                                  : screenWidth * 0.06,
                            )
                          ],
                        )),
                    Align(
                      alignment: Alignment.topLeft,
                      child: Text(
                        nameAREA,
                        textAlign: TextAlign.start,
                        style: TextStyle(
                          fontSize: screenWidth < screenHeight
                              ? screenWidth * 0.04
                              : screenWidth * 0.02,
                          fontWeight: FontWeight.w900,
                          color: const Color.fromARGB(255, 255, 255, 255),
                          fontFamily: 'Nunito-Bold'
                        ),
                      ),
                    )
                  ],
                ),
                Align(
                  alignment: Alignment.bottomLeft,
                  child: Padding(
                    padding: const EdgeInsets.all(4.0),
                    child: Text(
                      nameService,
                      style: TextStyle(
                        fontSize: screenWidth < screenHeight
                            ? screenWidth * 0.045
                            : screenWidth * 0.02,
                        fontWeight: FontWeight.w900,
                        color: const Color.fromARGB(255, 255, 255, 255),
                        fontFamily: 'Nunito-Bold'
                      ),
                    ),
                  ),
                )
              ],
            )));
  }
}

class AppletSection extends StatelessWidget {
  AppletSection({super.key});

  final List<Applet> applets = [
    Applet(
      nameService: "Discord",
      nameAREA: "In 10sec receive message on Discord",
      icon1: Icons.discord,
      icon2: Icons.timer,
      color: const Color(0xff7289da),
      route: "/discordarea",
      press: () {},
    ),
  ];

  @override
  Widget build(BuildContext context) {
    double screenHeight = MediaQuery.of(context).size.height;
    double screenWidth = MediaQuery.of(context).size.width;
    return Container(
      alignment: Alignment.center,
      decoration: const BoxDecoration(
        color: Colors.white,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Container(
            margin: const EdgeInsets.only(top: 20),
            child: Text(
              "Get started with any Applet",
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: screenWidth < screenHeight
                    ? screenWidth * 0.045
                    : screenWidth * 0.02,
                fontWeight: FontWeight.w900,
                color: Colors.black,
                fontFamily: 'Nunito-Bold'
              ),
            ),
          ),
          Wrap(
            spacing: screenWidth < screenHeight
                ? screenWidth * 0.040
                : screenWidth * 0.02,
            alignment: WrapAlignment.center,
            children: applets
                .map((applet) => _buildAppletCard(context, applet))
                .toList(),
          ),
        ],
      ),
    );
  }

  Widget _buildAppletCard(BuildContext context, Applet applet) {
    return Applet(
        nameService: applet.nameService,
        nameAREA: applet.nameAREA,
        icon1: applet.icon1,
        icon2: applet.icon2,
        color: applet.color,
        press: () {
          context.go(applet.route);
        },
        route: applet.route);
  }
}
