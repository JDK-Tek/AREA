import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';
import 'package:area/tools/dynamic.dart';
import 'package:http/http.dart' as https;
import 'dart:convert';
import 'package:area/tools/providers.dart';
import 'package:provider/provider.dart';

class CreateAutomationPage extends StatefulWidget {
  const CreateAutomationPage({super.key});

  @override
  CreateAutomationPageState createState() => CreateAutomationPageState();
}

class CreateAutomationPageState extends State<CreateAutomationPage> {
  List<dynamic> triggers = [];
  List<dynamic> triggersConfigurations = [];
  dynamic selectedTrigger;
  List<dynamic> reactions = [];
  List<dynamic> reactionsConfigurations = [];
  dynamic selectedReaction;

  Map<String, dynamic> triggerValues = {};
  Map<String, dynamic> reactionValues = {};

  @override
  void initState() {
    super.initState();
    _loadMockTriggers();
    _loadMockReactions();
  }

  void _loadMockTriggers() {
    triggers = [
      {"label": "Time", "icon_url": "https://img.icons8.com/ios/452/timer.png"},
    ];
    triggersConfigurations = [
      {
        "type": "action",
        "name": "in",
        "spices": [
          {"name": "howmuch", "type": "number", "extraParams": null},
          {
            "name": "unit",
            "type": "dropdown",
            "extraParams": ["weeks", "days", "hours", "minutes", "seconds"]
          }
        ]
      }
    ];
    setState(() {});
  }

  void _loadMockReactions() {
    reactions = [
      {
        "label": "Discord",
        "icon_url": "https://img.icons8.com/ios/452/discord.png",
      },
    ];
    reactionsConfigurations = [
      {
        "type": "reaction",
        "name": "send",
        "spices": [
          {
            "name": "channel",
            "type": "number",
          },
          {
            "name": "message",
            "type": "text",
          }
        ]
      }
    ];
    setState(() {});
  }

  Widget _buildDynamicTriggerConfig() {
    if (selectedTrigger == null || triggersConfigurations.isEmpty) {
      return const SizedBox.shrink();
    }

    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.start,
        children: triggersConfigurations.expand<Widget>((config) {
          return config['spices']?.map<Widget>((childConfig) {
                return Dynamic(
                  title: childConfig['type'],
                  extraParams: {
                    'items': childConfig['extraParams'],
                  },
                  onValueChanged: (key, value) {
                    setState(() {
                      if (childConfig['name'] == 'howmuch') {
                        triggerValues[childConfig['name']] = int.parse(value);
                      } else {
                        triggerValues[childConfig['name']] = value;
                      }
                    });
                  },
                );
              }).toList() ??
              [];
        }).toList(),
      ),
    );
  }

  Widget _buildDynamicReactionConfig() {
    if (selectedReaction == null || reactionsConfigurations.isEmpty) {
      return const SizedBox.shrink();
    }

    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        children: reactionsConfigurations.expand<Widget>((config) {
          return config['spices']?.map<Widget>((childConfig) {
                return Dynamic(
                  title: childConfig['type'],
                  extraParams: {
                    'items': childConfig['extraParams'],
                  },
                  onValueChanged: (key, value) {
                    setState(() {
                      reactionValues[childConfig['name']] = value;
                    });
                  },
                );
              }).toList() ??
              [];
        }).toList(),
      ),
    );
  }

  void _submitAutomation() {
    if (selectedTrigger == null || selectedReaction == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
            content: Text("Please select both a trigger and a reaction.")),
      );
      return;
    }

    final automation = {
      "action": {
        "service": selectedTrigger['label'].toLowerCase(),
        "name": triggersConfigurations[0]['name'],
        "spices": triggerValues,
      },
      "reaction": {
        "service": selectedReaction['label'].toLowerCase(),
        "name": reactionsConfigurations[0]['name'],
        "spices": reactionValues,
      },
    };

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
          content: Text("Automation Submitted: \n${automation.toString()}")),
    );
    _sendRequest(automation);
  }

  Future<void> _sendRequest(body) async {
    final token = Provider.of<UserState>(context, listen: false).token;
    print("$token");
    final Uri uri =
        Uri.https(Provider.of<IPState>(context, listen: false).ip, "/api/area");
    print("Request URL: $uri");
    final Map<String, String> headers = {
      "Authorization": "Bearer $token",
      "Content-Type": "application/json",
    };
   print("Request body: ${jsonEncode(body)}");

    try {
      final response =
          await https.post(uri, headers: headers, body: jsonEncode(body));

      if (response.statusCode == 200) {
        final Map<String, dynamic> data =
            jsonDecode(response.body) as Map<String, dynamic>;
        _showDialog("Success", "Request sent successfully: $data");
      } else {
        _showDialog("Error",
            "Failed with status: ${response.statusCode}. ${response.reasonPhrase ?? 'Unknown error'}");
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
    int currentPageIndex = 1;
    final List<String> routes = [
      "/applets",
      "/create",
      "/services",
      "/developers"
    ];

    return SafeArea(
      child: Scaffold(
        body: SingleChildScrollView(
          child: Column(
            children: [
              const SizedBox(height: 20),
              Text(
                "If this ...",
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              const SizedBox(height: 10),
              SizedBox(
                height: 200,
                child: GridView.builder(
                  scrollDirection: Axis.vertical,
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    mainAxisSpacing: 10,
                  ),
                  itemCount: triggers.length,
                  itemBuilder: (context, index) {
                    final trigger = triggers[index];
                    return GestureDetector(
                      onTap: () {
                        setState(() {
                          selectedTrigger = trigger;
                        });
                      },
                      child: Container(
                        decoration: BoxDecoration(
                          color: selectedTrigger == trigger
                              ? Colors.white
                              : Colors.grey,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        padding: const EdgeInsets.all(8),
                        child: Column(
                          children: [
                            Expanded(
                              child: Image.network(
                                trigger['icon_url'],
                                fit: BoxFit.contain,
                                errorBuilder: (context, error, stackTrace) =>
                                    const Icon(Icons.error),
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(trigger['label']),
                          ],
                        ),
                      ),
                    );
                  },
                ),
              ),
              _buildDynamicTriggerConfig(),
              const SizedBox(height: 20),
              Text(
                "Then that...",
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              const SizedBox(height: 10),
              SizedBox(
                height: 200,
                child: GridView.builder(
                  scrollDirection: Axis.vertical,
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    mainAxisSpacing: 10,
                  ),
                  itemCount: reactions.length,
                  itemBuilder: (context, index) {
                    final reaction = reactions[index];
                    return GestureDetector(
                      onTap: () {
                        setState(() {
                          selectedReaction = reaction;
                        });
                      },
                      child: Container(
                        decoration: BoxDecoration(
                          color: selectedReaction == reaction
                              ? Colors.white
                              : Colors.grey,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        padding: const EdgeInsets.all(8),
                        child: Column(
                          children: [
                            Expanded(
                              child: Image.network(
                                reaction['icon_url'],
                                fit: BoxFit.contain,
                                errorBuilder: (context, error, stackTrace) =>
                                    const Icon(Icons.error),
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(reaction['label']),
                          ],
                        ),
                      ),
                    );
                  },
                ),
              ),
              _buildDynamicReactionConfig(),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: _submitAutomation,
                child: const Text("Submit"),
              ),
            ],
          ),
        ),
        backgroundColor: Colors.white,
        bottomNavigationBar: NavigationBar(
          labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
          backgroundColor: Colors.black,
          indicatorColor: Colors.grey,
          selectedIndex: currentPageIndex,
          onDestinationSelected: (index) {
            setState(() {
              context.go(routes[index]);
            });
          },
          destinations: const [
            NavigationDestination(
              icon: Icon(Icons.folder, color: Colors.white),
              label: 'Applets',
            ),
            NavigationDestination(
              icon: Icon(Icons.add_circle_outline, color: Colors.white),
              label: 'Create',
            ),
            NavigationDestination(
              icon: Icon(Icons.cloud, color: Colors.white),
              label: 'Services',
            ),
            NavigationDestination(
              icon: Icon(CupertinoIcons.ellipsis, color: Colors.white),
              label: 'Developers',
            ),
          ],
        ),
      ),
    );
  }
}
