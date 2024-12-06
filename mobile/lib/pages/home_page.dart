import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:footer/footer_view.dart';
import 'package:mobile/tools/footerarea.dart';

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
        children: const [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [HeaderSection(), AppletSection()],
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
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Column(
            children: [
              Align(
                alignment: Alignment.topRight,
                child: ElevatedButton(
                onPressed: () {
                  context.go("/login");
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
              )
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
                  ? screenWidth * 0.035
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
          // Align(
          //   alignment: Alignment.bottomCenter,
          //   child: Image.asset(
          //     "assets/deco.png",
          //     fit: BoxFit.fitWidth,
          //   ),
          // )
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
  const AppletSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        color: Color.fromARGB(255, 255, 255, 255),
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
                color: const Color.fromARGB(255, 0, 0, 0),
              ),
            ),
          ),
          Wrap(
            spacing: 5.0,
            alignment: WrapAlignment.center,
            children: [
              Applet(
                icon1: Icons.audiotrack,
                icon2: Icons.favorite,
                nameAREA: "Exemple",
                nameService: "Google",
                press: () {
                  print("mon code");
                },
                color: const Color(0xff05b348),
              ),
              Applet(
                icon1: Icons.audiotrack,
                icon2: Icons.favorite,
                nameAREA: "Exemple",
                nameService: "Google",
                press: () {
                  print("mon code");
                },
                color: const Color(0xff222222),
              ),
              Applet(
                icon1: Icons.audiotrack,
                icon2: Icons.favorite,
                nameAREA: "Exemple",
                nameService: "Google",
                press: () {
                  print("mon code");
                },
                color: const Color(0xff341d4f),
              ),
              Applet(
                icon1: Icons.audiotrack,
                icon2: Icons.favorite,
                nameAREA: "Exemple",
                nameService: "Google",
                press: () {
                  print("mon code");
                },
                color: const Color(0xffff1970),
              ),
            ],
          ),
          Container(
            alignment: const Alignment(0, 0),
            margin: const EdgeInsets.only(top: 20),
            child: Text(
              "or choose from 900+ services",
              textAlign: TextAlign.center,
              style: GoogleFonts.nunito(
                fontSize: 16.8,
                fontWeight: FontWeight.w900,
                color: const Color(0xff410cab),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
