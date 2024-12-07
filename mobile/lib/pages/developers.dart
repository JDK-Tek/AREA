import 'package:flutter/material.dart';
import 'package:footer/footer_view.dart';
import 'package:area/tools/footerarea.dart';
import 'package:area/pages/home_page.dart';



class DevelopersPage extends StatelessWidget {
  const DevelopersPage({super.key});

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
            children: [HeaderSection(), Text("Elise Esteban Greg Paul John", textAlign: TextAlign.center,)],
          ),
        ],
      ),
    ));
  }
}
