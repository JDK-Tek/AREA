import 'package:area/tools/action_reaction.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:area/pages/home_page.dart';
import 'package:area/tools/timer.dart';
import 'package:area/tools/discord.dart';
import 'package:flutter/cupertino.dart';
import 'package:area/tools/userstate.dart';
import 'package:provider/provider.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class CreateAutomationPage extends StatefulWidget {
  const CreateAutomationPage({super.key});

  @override
  CreateAutomationPageState createState() => CreateAutomationPageState();
}

class CreateAutomationPageState extends State<CreateAutomationPage> {
  String? selectedAction;
  String? selectedReaction;
  int indexAction = -1;
  int indexReaction = -1;
  ActionHandler? selectedActionHandler;
  ReactionHandler? selectedReactionHandler;

  Future<void> _sendRequest() async {
    final token = context.read<UserState>().token;
    final Uri uri = Uri.http("172.20.10.3:42000", "/api/area");
    final Map<String, String> headers = {
      "Authorization": "Bearer $token",
      "Content-Type": "application/json",
    };

    final Map<String, dynamic> body = {
      "action": selectedActionHandler!.toJson(),
      "reaction": selectedReactionHandler!.toJson(),
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
        _showDialog("Error",
            "Failed with status: ${response.statusCode}. ${error['message'] ?? 'Unknown error'}");
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

  final List<Map<String, dynamic>> triggers = [
    {
      "label": "Time",
      "icon": Icons.access_time,
      "builder": (Function(Map<String, int>) callback) =>
          TimeTrigger(onTriggerChanged: callback),
      "description": "In"
    },
  ];

  final List<Map<String, dynamic>> reactions = [
    {
      "label": "Send a Discord message",
      "icon": Icons.discord,
      "builder": (Function(String, String) callback) =>
          DiscordAction(onActionChanged: callback),
      "description": "post a message to a channel"
    },
  ];

  Map<String, dynamic> triggerData = {};
  Map<String, dynamic> reactionData = {};

  int getGridColumnCount(BuildContext context) {
    return MediaQuery.of(context).size.width <
            MediaQuery.of(context).size.height
        ? 2
        : 3;
  }

  Widget _buildGridItem({
    required String label,
    required IconData icon,
    required bool isSelected,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          color: isSelected
              ? Colors.blueAccent.withOpacity(0.1)
              : Colors.grey[200],
          borderRadius: BorderRadius.circular(12),
          border: isSelected
              ? Border.all(color: Colors.blueAccent, width: 2)
              : Border.all(color: Colors.grey[300]!),
        ),
        padding: const EdgeInsets.all(16),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              color: isSelected ? Colors.blueAccent : Colors.grey[600],
              size: MediaQuery.of(context).size.width <
                      MediaQuery.of(context).size.height
                  ? MediaQuery.of(context).size.height * 0.03
                  : MediaQuery.of(context).size.height * 0.05,
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: TextStyle(
                color: isSelected ? Colors.blueAccent : Colors.grey[700],
                fontWeight: FontWeight.bold,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSummary(BuildContext context, String descriptionAction,
      String descriptionReaction) {
    return SizedBox(
      height: MediaQuery.of(context).size.height * 0.2,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text("Summary",
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: MediaQuery.of(context).size.width <
                        MediaQuery.of(context).size.height
                    ? MediaQuery.of(context).size.height * 0.02
                    : MediaQuery.of(context).size.height * 0.03,
              )),
          const SizedBox(height: 8),
          if (selectedAction != null)
            Text("Trigger: $selectedAction, $descriptionAction $triggerData"),
          const SizedBox(height: 8),
          if (selectedReaction != null)
            Text(
                "Reaction: $selectedReaction, $descriptionReaction $reactionData"),
          Text(
              "Name Applet: $descriptionAction $triggerData, $descriptionReaction")
        ],
      ),
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
          selectedIndex: 1,
          onDestinationSelected: (int index) {
            setState(() {
              context.go(dest[index]);
            });
          },
          destinations: const [
            NavigationDestination(
                icon: Icon(Icons.folder, color: Colors.white),
                label: 'Applets'),
            NavigationDestination(
                icon: Icon(Icons.add_circle_outline, color: Colors.white),
                label: 'Create'),
            NavigationDestination(
                icon: Icon(Icons.cloud, color: Colors.white),
                label: 'Services'),
            NavigationDestination(
                icon: Icon(CupertinoIcons.ellipsis, color: Colors.white),
                label: 'Developers'),
          ],
        ),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              const MiniHeaderSection(),
              Text("If this ...",
                  style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: MediaQuery.of(context).size.width <
                              MediaQuery.of(context).size.height
                          ? MediaQuery.of(context).size.height * 0.03
                          : MediaQuery.of(context).size.height * 0.04)),
              SizedBox(
                height: MediaQuery.of(context).size.height * 0.35,
                width: MediaQuery.of(context).size.width * 0.9,
                child: GridView.builder(
                  gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: getGridColumnCount(context),
                  ),
                  itemCount: triggers.length,
                  itemBuilder: (context, index) {
                    final trigger = triggers[index];
                    indexAction = index;
                    return _buildGridItem(
                      label: trigger['label'],
                      icon: trigger['icon'],
                      isSelected: selectedAction == trigger['label'],
                      onTap: () {
                        setState(() {
                          selectedAction = trigger['label'];
                          indexAction = index;
                          selectedActionHandler = trigger['builder']((data) {
                            setState(() {
                              triggerData = data;
                            });
                          });
                        });
                      },
                    );
                  },
                ),
              ),
              if (selectedAction != null && indexAction != -1)
                triggers[indexAction]['builder'](
                  (data) {
                    setState(() {
                      triggerData = data;
                    });
                  },
                ),
              Text("Then that ...",
                  style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: MediaQuery.of(context).size.width <
                              MediaQuery.of(context).size.height
                          ? MediaQuery.of(context).size.height * 0.03
                          : MediaQuery.of(context).size.height * 0.04)),
              SizedBox(
                height: MediaQuery.of(context).size.height * 0.35,
                width: MediaQuery.of(context).size.width * 0.9,
                child: GridView.builder(
                  gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: getGridColumnCount(context),
                  ),
                  itemCount: reactions.length,
                  itemBuilder: (context, index) {
                    final reaction = reactions[index];
                    return _buildGridItem(
                      label: reaction['label'],
                      icon: reaction['icon'],
                      isSelected: selectedReaction == reaction['label'],
                      onTap: () {
                        setState(() {
                          indexReaction = index;
                          selectedReaction = reaction['label'];
                          selectedReactionHandler = reaction['builder'](
                              (String channel, String message) {
                            setState(() {
                              reactionData = {
                                "channel": channel,
                                "message": message
                              };
                            });
                          });
                        });
                      },
                    );
                  },
                ),
              ),
              if (selectedReaction != null && indexReaction != -1)
                reactions[indexReaction]['builder'](
                  (String channel, String message) {
                    setState(() {
                      reactionData = {"channel": channel, "message": message};
                    });
                  },
                ),
              SizedBox(height: MediaQuery.of(context).size.height * 0.02),
              if (selectedAction != null && selectedReaction != null)
                _buildSummary(context, triggers[indexAction]['description'],
                    reactions[indexReaction]['description']),
              SizedBox(height: MediaQuery.of(context).size.height * 0.02),
              if (selectedAction != null && selectedReaction != null)
                Center(
                  child: SizedBox(
                    width: MediaQuery.of(context).size.width * 0.5,
                    height: 50,
                    child: ElevatedButton(
                      onPressed: () {
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(
                              content:
                                  Text("Automation created successfully!")),
                        );
                        _sendRequest();
                      },
                      child: const Text("Create Automation"),
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
