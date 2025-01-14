import 'package:flutter/material.dart';
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
    final List<String> dest = ["/applets", "/create", "/services", "/plus"];
    return SafeArea(
        child: Scaffold(
      bottomNavigationBar: NavigationBar(
        labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
        backgroundColor: Colors.black,
        indicatorColor: Colors.grey,
        shadowColor: Colors.transparent,
        surfaceTintColor: Colors.transparent,
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
  final String icon;
  final VoidCallback onPress;

  const Service({
    super.key,
    required this.serviceName,
    required this.icon,
    required this.onPress,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width:
          MediaQuery.of(context).size.width < MediaQuery.of(context).size.height
              ? MediaQuery.of(context).size.width * 0.25
              : MediaQuery.of(context).size.width * 0.2,
      height:
          MediaQuery.of(context).size.width < MediaQuery.of(context).size.height
              ? MediaQuery.of(context).size.height * 0.15
              : MediaQuery.of(context).size.width * 0.15,
      margin: const EdgeInsets.only(top: 20),
      child: ElevatedButton(
        onPressed: onPress,
        style: ElevatedButton.styleFrom(
          padding: const EdgeInsets.all(10),
          backgroundColor: Colors.transparent,
          shadowColor: Colors.transparent,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15),
          ),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.network(
              icon,
              width: MediaQuery.of(context).size.width <
                        MediaQuery.of(context).size.height ? MediaQuery.of(context).size.width * 0.2 : MediaQuery.of(context).size.width * 0.08,
            ),
            SizedBox(
                height: MediaQuery.of(context).size.width <
                        MediaQuery.of(context).size.height
                    ? 10
                    : 5),
            Text(
              serviceName,
              textAlign: TextAlign.center,
              style: TextStyle(
                  fontSize: MediaQuery.of(context).size.width <
                          MediaQuery.of(context).size.height
                      ? MediaQuery.of(context).size.width * 0.038
                      : MediaQuery.of(context).size.width * 0.025,
                  fontWeight: FontWeight.w700,
                  color: Colors.black,
                  fontFamily: 'Nunito-Bold'),
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
        "label": "Discord",
        "icon_url": "https://upload.wikimedia.org/wikipedia/fr/8/80/Logo_Discord_2015.png",
      },
      {
        "label": "Discord",
        "icon_url": "https://img.icons8.com/ios/452/discord.png",
      },
      {
        "label": "Discord",
        "icon_url": "https://img.icons8.com/ios/452/discord.png",
      },
    ];

    return Container(
      alignment: Alignment.center,
      padding: const EdgeInsets.symmetric(vertical: 20),
      child: Column(
        children: [
          const Text(
            "Services Available",
            style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: Colors.black,
                fontFamily: 'Nunito-Bold'),
          ),
          const SizedBox(height: 20),
          Wrap(
            spacing: 10.0,
            alignment: WrapAlignment.center,
            children: services
                .map(
                  (service) => Service(
                    serviceName: service['label'],
                    icon: service['icon_url'],
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
