import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:footer/footer_view.dart';
import 'package:area/tools/footerarea.dart';

List<Color> predefinedColors = [
  const Color(0xff410cab),
  const Color(0xff222222),
  const Color(0xffa6a6a6),
  const Color(0xff5e17eb),
];

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return SafeArea(
        child: Scaffold(
      backgroundColor: Colors.white,
      body: FooterView(
        footer: const Footerarea().build(context),
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const HeaderSection(),
              AppletSection(),
              const ServiceSection()
            ],
          ),
        ],
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
      decoration: const BoxDecoration(
        image: DecorationImage(
          image: AssetImage("assets/background_home.png"),
          fit: BoxFit.cover,
        ),
      ),
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
                        ? screenWidth * 0.2 : screenWidth * 0.1,
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
              color: const Color(0xff5f18eb),
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
                context.go("/");
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

class Applet extends StatelessWidget {
  final String nameService;
  final IconData icon1;
  final IconData icon2;
  final String nameAREA;
  final VoidCallback press;
  final Color color;

  const Applet(
      {super.key,
      required this.icon1,
      required this.icon2,
      required this.nameService,
      required this.nameAREA,
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
            screenWidth < screenHeight ? screenHeight * 0.2 : screenWidth * 0.3,
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
                              size: 30.0,
                            ),
                            Icon(
                              icon2,
                              color: Colors.white,
                              size: 30.0,
                            )
                          ],
                        )),
                    Align(
                      alignment: Alignment.topLeft,
                      child: Text(
                        nameAREA,
                        textAlign: TextAlign.start,
                        style: GoogleFonts.nunito(
                          fontSize: 16,
                          fontWeight: FontWeight.w900,
                          color: const Color.fromARGB(255, 255, 255, 255),
                        ),
                      ),
                    )
                  ],
                ),
                Align(
                  alignment: Alignment.bottomLeft,
                  child: Padding(
                    padding: const EdgeInsets.all(7.0),
                    child: Text(
                      nameService,
                      style: GoogleFonts.nunito(
                        fontSize: 16,
                        fontWeight: FontWeight.w900,
                        color: const Color.fromARGB(255, 255, 255, 255),
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
      nameAREA: "Every 10sec receive message on Discord",
      icon1: Icons.discord,
      icon2: Icons.timer,
      color: const Color(0xff7289da),
      press: () {
        print("Applet ${"Every 10sec receive message on Discord"} clicked");
      },
    ),
  ];

  @override
  Widget build(BuildContext context) {
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
              style: GoogleFonts.nunito(
                fontSize: 16.8,
                fontWeight: FontWeight.w900,
                color: Colors.black,
              ),
            ),
          ),
          Wrap(
            spacing: 10.0,
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
    return applet;
  }
}

class Service extends StatelessWidget {
  final String serviceName;
  final IconData icon;
  final VoidCallback onPress;
  final Color color;

  const Service({
    super.key,
    required this.serviceName,
    required this.icon,
    required this.onPress,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    double screenHeight = MediaQuery.of(context).size.height;
    double screenWidth = MediaQuery.of(context).size.width;

    return Container(
      width:
          screenWidth < screenHeight ? screenWidth * 0.25 : screenWidth * 0.2,
      height:
          screenWidth < screenHeight ? screenHeight * 0.15 : screenWidth * 0.15,
      margin: const EdgeInsets.only(top: 20),
      child: ElevatedButton(
        onPressed: onPress,
        style: ElevatedButton.styleFrom(
          backgroundColor: color,
          padding: const EdgeInsets.all(10),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15),
          ),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              color: Colors.white,
              size: screenWidth < screenHeight
                  ? screenWidth * 0.08
                  : screenWidth * 0.05,
            ),
            SizedBox(height: screenWidth < screenHeight ? 10 : 5),
            Text(
              serviceName,
              textAlign: TextAlign.center,
              style: GoogleFonts.nunito(
                fontSize: screenWidth < screenHeight
                    ? screenWidth * 0.04
                    : screenWidth * 0.025,
                fontWeight: FontWeight.w700,
                color: Colors.white,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class ServiceSection extends StatelessWidget {
  const ServiceSection({super.key});

  @override
  Widget build(BuildContext context) {
    final List<Map<String, dynamic>> services = [
      {
        'name': 'Discord',
        'icon': Icons.discord,
        'color': Colors.blueAccent,
      },
      {
        'name': 'Weather',
        'icon': Icons.cloud,
        'color': Colors.lightBlue,
      },
      {
        'name': 'Time',
        'icon': Icons.timer,
        'color': Colors.green,
      },
    ];

    return Container(
      alignment: Alignment.center,
      padding: const EdgeInsets.symmetric(vertical: 20),
      child: Column(
        children: [
          Text(
            "Services Available",
            style: GoogleFonts.nunito(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: Colors.black,
            ),
          ),
          const SizedBox(height: 20),
          Wrap(
            spacing: 10.0,
            alignment: WrapAlignment.center,
            children: services
                .map(
                  (service) => Service(
                    serviceName: service['name'],
                    icon: service['icon'],
                    color: service['color'],
                    onPress: () {
                      print("Service ${service['name']} pressed");
                    },
                  ),
                )
                .toList(),
          ),
        ],
      ),
    );
  }
}
