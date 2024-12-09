import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';
import 'package:area/pages/home_page.dart';
import 'package:area/pages/appletspage.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:area/tools/userstate.dart';
import 'package:provider/provider.dart';

class DiscordAreaPage extends StatefulWidget {
  final Function(String, String) onActionChanged;

  const DiscordAreaPage({super.key, required this.onActionChanged});

  @override
  State<DiscordAreaPage> createState() => DiscordAreaPageState();
}

class DiscordAreaPageState extends State<DiscordAreaPage> {
  int currentPageIndex = 3;
  final TextEditingController channelIdController = TextEditingController();
  final TextEditingController messageTemplateController =
      TextEditingController();


  Future<void> _sendRequest(String channelId, String message) async {
    final token = Provider.of<UserState>(context).token;
    final Uri uri = Uri.http("localhost:42000", "/api/area");
    final Map<String, String> headers = {
      "Authorization": "Bearer $token",
      "Content-Type": "application/json",
    };

    final Map<String, dynamic> body = {
      "action": {
        "service": "time",
        "name": "in",
        "spices": {"howmuch": 10, "unit": "secondes"}
      },
      "reaction": {
        "service": "discord",
        "name": "send",
        "spices": {
          "channel": int.tryParse(channelId) ?? 0,
          "message": message.isNotEmpty ? message : "Default Message"
        }
      }
    };

    try {
      final response =
          await http.post(uri, headers: headers, body: jsonEncode(body));

      if (response.statusCode == 200) {
        final Map<String, dynamic> data =
            jsonDecode(response.body) as Map<String, dynamic>;
        _showDialog("Success", "Request sent successfully: $data");
      } else {
        final Map<String, dynamic> error =
            jsonDecode(response.body) as Map<String, dynamic>;
        _showDialog("Error", "Failed with status: ${response.statusCode}. ${error['message'] ?? 'Unknown error'}");
      }
    } catch (e) {
      _showDialog("Error", "An exception occurred: $e");
    }
  }

  void _showDialog(String title, String message) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text(title),
          content: Text(message),
          actions: <Widget>[
            TextButton(
              child: const Text("OK"),
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
          ],
        );
      },
    );
  }

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
        selectedIndex: 3,
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
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            const MiniHeaderSection(),
            Applet(
                icon1: Icons.discord,
                icon2: Icons.timer,
                nameService: "Discord",
                nameAREA: "In 10 sec receive message on Discord",
                route: "/discordarea",
                press: () {},
                color: const Color(0xff7289da)),
            const SizedBox(height: 8),
            TextField(
              controller: channelIdController,
              decoration: const InputDecoration(
                labelText: "Discord Channel ID",
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: messageTemplateController,
              decoration: const InputDecoration(
                labelText: "Message Template (use {time} for time)",
                border: OutlineInputBorder(),
              ),
            ),
            ElevatedButton(
                onPressed: () {
                  final String channelId = channelIdController.text.trim();
                  final String messageTemplate =
                      messageTemplateController.text.trim();

                  if (channelId.isNotEmpty && messageTemplate.isNotEmpty) {
                    _sendRequest(channelId, messageTemplate);
                  } else {
                    _showDialog(
                        "Error", "Please fill in all required fields.");
                  }
                },
                child: const Icon(Icons.check_box, color: Colors.green))
          ],
        ),
      ),
    ));
  }
}
