import 'package:flutter/material.dart';
import 'package:footer/footer.dart';
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
                    Image.asset(
                      'assets/fullLogo.png',
                      height: screenWidth < screenHeight
                          ? screenHeight * 0.03
                          : screenHeight * 0.08,
                      width: screenWidth < screenHeight
                          ? screenWidth * 0.2
                          : screenWidth * 0.1,
                      fit: BoxFit.contain,
                    ),
                  ],
                )),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
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
                    style: TextStyle(
                      fontWeight: FontWeight.w900,
                      fontSize: screenWidth < screenHeight
                          ? screenWidth * 0.04
                          : screenWidth * 0.02,
                      color: Colors.white,
                      fontFamily: 'Nunito-Bold'
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
                    style: TextStyle(
                      fontWeight: FontWeight.w900,
                      fontSize: screenWidth < screenHeight
                          ? screenWidth * 0.04
                          : screenWidth * 0.02,
                      color: Colors.white,
                      fontFamily: 'Nunito-Bold'
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
                    style: TextStyle(
                      fontWeight: FontWeight.w900,
                      fontSize: screenWidth < screenHeight
                          ? screenWidth * 0.04
                          : screenWidth * 0.02,
                      color: Colors.white,
                      fontFamily: 'Nunito-Bold'
                    ),
                  ),
                ),
              ],
            )
          ],
        ));
  }
}
