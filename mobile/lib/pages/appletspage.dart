import 'dart:convert';
import 'package:http/http.dart' as https;
import 'package:area/tools/providers.dart';
import 'package:flutter/material.dart';
import 'package:area/pages/home_page.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';

class AppletsPage extends StatefulWidget {
  const AppletsPage({super.key});
  @override
  State<AppletsPage> createState() => AppletPageState();
}

class AppletPageState extends State<AppletsPage> {
  int currentPageIndex = 0;

  @override
  Widget build(BuildContext context) {
    final List<String> dest = ["/applets", "/create", "/services", "/plus"];
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
      body: const SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [HeaderSection(), AppletSection()],
        ),
      ),
    ));
  }
}

class Applet extends StatelessWidget {
  final String nameService;
  final String icon1;
  final String icon2;
  final String nameAREA;
  final String color;
  final VoidCallback press;

  const Applet({
    super.key,
    required this.icon1,
    required this.icon2,
    required this.nameService,
    required this.nameAREA,
    required this.press,
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
    return Container(
        width: MediaQuery.of(context).size.width <
                MediaQuery.of(context).size.height
            ? MediaQuery.of(context).size.width * 0.65
            : MediaQuery.of(context).size.width * 0.38,
        height: MediaQuery.of(context).size.width <
                MediaQuery.of(context).size.height
            ? MediaQuery.of(context).size.height * 0.25
            : MediaQuery.of(context).size.width * 0.3,
        margin: const EdgeInsets.only(top: 20),
        child: ElevatedButton(
            onPressed: () {
              press();
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: _colorFromHex(color),
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
                            Image.network(
                              loadingBuilder: (BuildContext context,
                                  Widget child,
                                  ImageChunkEvent? loadingProgress) {
                                if (loadingProgress == null) {
                                  return child;
                                } else {
                                  return Center(
                                    child: CircularProgressIndicator(
                                      value:
                                          loadingProgress.expectedTotalBytes !=
                                                  null
                                              ? loadingProgress
                                                      .cumulativeBytesLoaded /
                                                  (loadingProgress
                                                          .expectedTotalBytes ??
                                                      1)
                                              : null,
                                    ),
                                  );
                                }
                              },
                              errorBuilder: (BuildContext context, Object error,
                                  StackTrace? stackTrace) {
                                return const Icon(Icons.broken_image, size: 40);
                              },
                              icon1,
                              color: Colors.white,
                              width: MediaQuery.of(context).size.width <
                                      MediaQuery.of(context).size.height
                                  ? MediaQuery.of(context).size.width * 0.1
                                  : MediaQuery.of(context).size.width * 0.06,
                            ),
                            Image.network(
                              loadingBuilder: (BuildContext context,
                                  Widget child,
                                  ImageChunkEvent? loadingProgress) {
                                if (loadingProgress == null) {
                                  return child;
                                } else {
                                  return Center(
                                    child: CircularProgressIndicator(
                                      value:
                                          loadingProgress.expectedTotalBytes !=
                                                  null
                                              ? loadingProgress
                                                      .cumulativeBytesLoaded /
                                                  (loadingProgress
                                                          .expectedTotalBytes ??
                                                      1)
                                              : null,
                                    ),
                                  );
                                }
                              },
                              errorBuilder: (BuildContext context, Object error,
                                  StackTrace? stackTrace) {
                                return const Icon(Icons.broken_image, size: 40);
                              },
                              icon2,
                              color: Colors.white,
                              width: MediaQuery.of(context).size.width <
                                      MediaQuery.of(context).size.height
                                  ? MediaQuery.of(context).size.width * 0.1
                                  : MediaQuery.of(context).size.width * 0.06,
                            )
                          ],
                        )),
                    Align(
                      alignment: Alignment.topLeft,
                      child: Text(
                        nameAREA,
                        textAlign: TextAlign.start,
                        style: TextStyle(
                            fontSize: MediaQuery.of(context).size.width <
                                    MediaQuery.of(context).size.height
                                ? MediaQuery.of(context).size.width * 0.05
                                : MediaQuery.of(context).size.width * 0.03,
                            fontWeight: FontWeight.w900,
                            color: const Color.fromARGB(255, 255, 255, 255),
                            fontFamily: 'Nunito-Bold'),
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
                          fontSize: MediaQuery.of(context).size.width <
                                  MediaQuery.of(context).size.height
                              ? MediaQuery.of(context).size.width * 0.05
                              : MediaQuery.of(context).size.width * 0.03,
                          fontWeight: FontWeight.w900,
                          color: const Color.fromARGB(255, 255, 255, 255),
                          fontFamily: 'Nunito-Bold'),
                    ),
                  ),
                )
              ],
            )));
  }
}

class AppletSection extends StatefulWidget {
  const AppletSection({super.key});

  @override
  State<AppletSection> createState() => _AppletSectionState();
}

class _AppletSectionState extends State<AppletSection> {
  List<Map<String, dynamic>> applets = [];

  @override
  void initState() {
    super.initState();
    _makeDemand("/api/applets");
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

    if (responseBody.containsKey('res')) {
      final List<dynamic> appletsList = responseBody['res'];
      if (mounted) {
        setState(() {
          applets = List<Map<String, dynamic>>.from(appletsList);
        });
      }
    } else {
      if (mounted) {
        _showDialog("Error", "Key 'server.applets' not found in response.");
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
                  fontSize: MediaQuery.of(context).size.width <
                          MediaQuery.of(context).size.height
                      ? MediaQuery.of(context).size.width * 0.045
                      : MediaQuery.of(context).size.width * 0.02,
                  fontWeight: FontWeight.w900,
                  color: Colors.black,
                  fontFamily: 'Nunito-Bold'),
            ),
          ),
          Wrap(
            spacing: MediaQuery.of(context).size.width <
                    MediaQuery.of(context).size.height
                ? MediaQuery.of(context).size.width * 0.05
                : MediaQuery.of(context).size.width * 0.02,
            alignment: WrapAlignment.center,
            children: applets
                .map((applet) => _buildAppletCard(
                    applet["service"]["name"],
                    applet["name"],
                    applet["service"]["logo"],
                    applet["service"]["logopartner"],
                    applet["service"]["color"]["normal"]))
                .toList(),
          ),
        ],
      ),
    );
  }

  Widget _buildAppletCard(String nameService, String nameAREA, String icon1,
      String icon2, String color) {
    return Applet(
      nameService: nameService,
      nameAREA: nameAREA,
      icon1: icon1,
      icon2: icon2,
      color: color,
      press: () {},
    );
  }
}
