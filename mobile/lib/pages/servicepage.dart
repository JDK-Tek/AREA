import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:area/pages/home_page.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';

class ServicesPage extends StatefulWidget {
  const ServicesPage({super.key});

  @override
  State<ServicesPage> createState() => ServicesPageState();
}

class ServicesPageState extends State<ServicesPage> {
  int currentPageIndex = 2;
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
        selectedIndex: 2,
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
          children: [MiniHeaderSection(), ServiceSection()],
        ),
      ),
    ));
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
