import 'package:flutter/material.dart';
import 'package:footer/footer.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:go_router/go_router.dart';


class Footerarea extends StatelessWidget {
  const Footerarea({super.key});

  @override
  Footer build(BuildContext context) {
    double screenHeight = MediaQuery.of(context).size.height;
    double screenWidth = MediaQuery.of(context).size.width;
    return Footer(
        backgroundColor: const Color(0xff222222),
        child: Column(
          children: [
            Align(
                alignment: Alignment.topLeft,
                child: Row(
                  children: [
                    const Icon(
                      Icons.logo_dev,
                      color: Colors.white,
                      size: 30.0,
                    ),
                    Text(
                      "AREA",
                      style: GoogleFonts.nunito(
                        fontSize: screenWidth < screenHeight
                            ? screenWidth * 0.05
                            : screenWidth * 0.04,
                        fontWeight: FontWeight.w900,
                        color: Colors.white,
                      ),
                    ),
                  ],
                )),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                ElevatedButton(
                  onPressed: () {
                    context.go("/developers");
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color.fromARGB(0, 0, 0, 0),
                    foregroundColor: const Color.fromARGB(0, 0, 0, 0),
                    shadowColor: const Color.fromARGB(0, 0, 0, 0),
                  ),
                  child: Text(
                    "Developers",
                    style: GoogleFonts.nunito(
                      fontWeight: FontWeight.w900,
                      fontSize: screenWidth < screenHeight
                          ? screenWidth * 0.04
                          : screenWidth * 0.02,
                      color: Colors.white,
                    ),
                  ),
                ),
                ElevatedButton(
                  onPressed: () {
                    context.go("/aboutus");
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color.fromARGB(0, 0, 0, 0),
                    foregroundColor: const Color.fromARGB(0, 0, 0, 0),
                    shadowColor: const Color.fromARGB(0, 0, 0, 0),
                  ),
                  child: Text(
                    "About us",
                    style: GoogleFonts.nunito(
                      fontWeight: FontWeight.w900,
                      fontSize: screenWidth < screenHeight
                          ? screenWidth * 0.04
                          : screenWidth * 0.02,
                      color: Colors.white,
                    ),
                  ),
                ),
                ElevatedButton(
                  onPressed: () {
                    context.go("/contact");
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color.fromARGB(0, 0, 0, 0),
                    foregroundColor: const Color.fromARGB(0, 0, 0, 0),
                    shadowColor: const Color.fromARGB(0, 0, 0, 0),
                  ),
                  child: Text(
                    "Contact",
                    style: GoogleFonts.nunito(
                      fontWeight: FontWeight.w900,
                      fontSize: screenWidth < screenHeight
                          ? screenWidth * 0.04
                          : screenWidth * 0.02,
                      color: Colors.white,
                    ),
                  ),
                ),
              ],
            )
          ],
        ));
  }
}
