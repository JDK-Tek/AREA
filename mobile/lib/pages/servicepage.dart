import 'package:area/tools/providers.dart';
import 'package:flutter/material.dart';
import 'package:area/pages/home_page.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';
import 'package:http/http.dart' as https;
import 'dart:convert';
import 'package:provider/provider.dart';

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
  final String color;

  const Service({
    super.key,
    required this.serviceName,
    required this.icon,
    required this.onPress,
    required this.color,
  });

  Color _colorFromHex(String hexColor) {
    if (hexColor.isEmpty) return Colors.grey;
    hexColor = hexColor.toUpperCase().replaceAll("#", "");
    if (hexColor.length == 6) {
      hexColor = "FF$hexColor";
    }
    return Color(int.tryParse(hexColor, radix: 16) ?? 0xFF000000);
  }

  @override
  Widget build(BuildContext context) {
    var ip = Provider.of<IPState>(context, listen: false).ip;
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
          backgroundColor: _colorFromHex(color),
          shadowColor: Colors.transparent,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15),
          ),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.network(
              loadingBuilder: (BuildContext context, Widget child,
                  ImageChunkEvent? loadingProgress) {
                if (loadingProgress == null) {
                  return child;
                } else {
                  return Center(
                    child: CircularProgressIndicator(
                      value: loadingProgress.expectedTotalBytes != null
                          ? loadingProgress.cumulativeBytesLoaded /
                              (loadingProgress.expectedTotalBytes ?? 1)
                          : null,
                    ),
                  );
                }
              },
              errorBuilder:
                  (BuildContext context, Object error, StackTrace? stackTrace) {
                print(icon);
                return const Icon(Icons.broken_image, size: 40);
              },
              "https://$ip" + icon,
              width: MediaQuery.of(context).size.width <
                      MediaQuery.of(context).size.height
                  ? MediaQuery.of(context).size.width * 0.2
                  : MediaQuery.of(context).size.width * 0.08,
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
                  color: Colors.white,
                  fontFamily: 'Nunito-Bold'),
            ),
          ],
        ),
      ),
    );
  }
}

class ServiceSection extends StatefulWidget {
  const ServiceSection({super.key});

  @override
  State<ServiceSection> createState() => _ServiceSectionState();
}

class _ServiceSectionState extends State<ServiceSection> {
  List<Map<String, dynamic>> services = [];

  @override
  void initState() {
    super.initState();
    _makeDemand("/about.json");
  }

  Future<void> _makeDemand(String u) async {
    final Uri uri =
        Uri.https(Provider.of<IPState>(context, listen: false).ip, u);
    late final https.Response rep;

    try {
      rep = await https.get(uri);
    } catch (e) {
      if (mounted) {
        _showDialog("Error", "Could not make request: $e");
      }
      return;
    }

    if (rep.statusCode >= 500) {
      if (mounted) {
        _showDialog("Error",
            "Failed with status: ${rep.statusCode}. ${rep.reasonPhrase ?? 'Unknown error'}");
      }
      return;
    }

    Map<String, dynamic> responseBody;
    try {
      responseBody = jsonDecode(rep.body);
    } catch (e) {
      if (mounted) {
        _showDialog("Error", "Invalid JSON format: $e");
      }
      return;
    }

    if (responseBody.containsKey('server') &&
        responseBody['server'].containsKey('services')) {
      final List<dynamic> servicesList = responseBody['server']['services'];
      if (mounted) {
        setState(() {
          services = List<Map<String, dynamic>>.from(servicesList);
        });
      }
    } else {
      if (mounted) {
        _showDialog("Error", "Key 'server.services' not found in response.");
      }
    }
  }

  void _showDialog(String title, String message) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(title),
          content: Text(message),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text("OK"),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
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
            children: services.map(
              (service) {
                final serviceName = service['name'] ?? 'Unknown';
                final icon = service['image'] ?? '';
                final color = service['color'] ?? '#FFFFFF';

                return Service(
                  serviceName: serviceName,
                  icon: icon,
                  color: color,
                  onPress: () {
                    print("Service $serviceName pressed");
                  },
                );
              },
            ).toList(),
          ),
        ],
      ),
    );
  }
}
