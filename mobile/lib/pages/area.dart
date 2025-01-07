import 'package:area/pages/home_page.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:area/tools/dynamic.dart';
import 'package:go_router/go_router.dart';
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
  List<dynamic> selectedTriggers = [];
  List<dynamic> reactions = [];
  List<dynamic> reactionsConfigurations = [];
  List<dynamic> selectedReactions = [];

  List<Map<String, dynamic>> triggerValues = [];
  List<Map<String, dynamic>> reactionValues = [];

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
          {"name": "channel", "type": "number"},
          {"name": "message", "type": "text"}
        ]
      }
    ];
    setState(() {});
  }

  void _addService(String type) {
    showModalBottomSheet(
      context: context,
      builder: (context) {
        final services = type == "Action" ? triggers : reactions;

        return ListView.builder(
          itemCount: services.length,
          itemBuilder: (context, index) {
            return ListTile(
              title: Row(children: [
                Image.network(
                  services[index]['icon_url'],
                  width: MediaQuery.of(context).size.width * 0.05,
                ),
                SizedBox(
                  width: MediaQuery.of(context).size.width * 0.01,
                ),
                Text(services[index]['label'])
              ]),
              onTap: () {
                setState(() {
                  if (type == "Action") {
                    selectedTriggers.add(services[index]);
                    if (index < triggersConfigurations.length) {
                      _buildDynamicConfig(selectedTriggers,
                          triggersConfigurations, triggerValues, index);
                    }
                  } else {
                    selectedReactions.add(services[index]);
                    if (index < reactionsConfigurations.length) {
                      _buildDynamicConfig(selectedReactions,
                          reactionsConfigurations, reactionValues, index);
                    }
                  }
                });
              },
            );
          },
        );
      },
    );
  }

  void _buildDynamicConfig(
      List<dynamic> selectedItems,
      List<dynamic> configurations,
      List<Map<String, dynamic>> values,
      int index) {
    showModalBottomSheet(
      useSafeArea: true,
      context: context,
      builder: (BuildContext context) {
        return Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: configurations.asMap().entries.expand<Widget>((entry) {
              var config = entry.value;
              return config['spices']?.map<Widget>((childConfig) {
                    return Dynamic(
                      title: childConfig['type'],
                      extraParams: {
                        'items': childConfig['extraParams'],
                      },
                      onValueChanged: (key, value) {
                        setState(() {
                          if (childConfig['type'] == "number") {
                            values.add({childConfig['name']: int.parse(value)});
                          } else {
                            values.add({childConfig['name']: value});
                          }
                        });
                      },
                    );
                  }).toList() ??
                  [];
            }).toList(),
          ),
        );
      },
    );
  }

  void _submitAutomation() {
    if (selectedTriggers.isEmpty || selectedReactions.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
            content:
                Text("Please select at least one trigger and one reaction.")),
      );
      return;
    }

    final automation = {
      "action": List.generate(selectedTriggers.length, (index) {
        return {
          "service": selectedTriggers[index]['label'].toLowerCase(),
          "name": triggersConfigurations[index]['name'],
          "spices": triggerValues[index],
        };
      }),
      "reaction": List.generate(selectedReactions.length, (index) {
        return {
          "service": selectedReactions[index]['label'].toLowerCase(),
          "name": reactionsConfigurations[index]['name'],
          "spices": reactionValues[index],
        };
      }),
    };

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
          content: Text("Automation Submitted: \n${jsonEncode(automation)}")),
    );
    _sendRequest(automation);
  }

  Future<void> _sendRequest(Map<String, dynamic> body) async {
    final token = Provider.of<UserState>(context, listen: false).token;
    final Uri uri =
        Uri.https(Provider.of<IPState>(context, listen: false).ip, "/api/area");

    final headers = {
      "Authorization": "Bearer $token",
      "Content-Type": "application/json",
    };

    try {
      final response =
          await https.post(uri, headers: headers, body: jsonEncode(body));

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
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

  int currentPageIndex = 1;
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
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const MiniHeaderSection(),
            Padding(
                padding: const EdgeInsets.all(16.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      'If this ...',
                      style:
                          TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                    ),
                    IconButton(
                      icon: const Icon(Icons.add),
                      onPressed: () => _addService("Action"),
                    ),
                  ],
                )),
            const Divider(),
            ...selectedTriggers.map((trigger) => ListTile(
                  key: ValueKey(trigger['label']),
                  title: Text(trigger['label']),
                  leading: Image.network(
                    trigger['icon_url'],
                    width: 40,
                    height: 40,
                    errorBuilder: (context, error, stackTrace) =>
                        const Icon(Icons.error),
                  ),
                  trailing: IconButton(
                    icon: const Icon(Icons.delete),
                    onPressed: () {
                      setState(() {
                        triggerValues
                            .removeAt(selectedTriggers.indexOf(trigger));
                        selectedTriggers.remove(trigger);
                      });
                    },
                  ),
                )),
            const SizedBox(height: 32),
            Padding(
                padding: const EdgeInsets.all(16.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      'then that...',
                      style:
                          TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                    ),
                    IconButton(
                      icon: const Icon(Icons.add),
                      onPressed: () => _addService("Reaction"),
                    ),
                  ],
                )),
            const Divider(),
            ...selectedReactions.map((reaction) => ListTile(
                  key: ValueKey(reaction['label']),
                  title: Text(reaction['label']),
                  leading: Image.network(
                    reaction['icon_url'],
                    width: 40,
                    height: 40,
                    errorBuilder: (context, error, stackTrace) =>
                        const Icon(Icons.error),
                  ),
                  trailing: IconButton(
                    icon: const Icon(Icons.delete),
                    onPressed: () {
                      setState(() {
                        reactionValues
                            .removeAt(selectedTriggers.indexOf(reaction));
                        selectedReactions.remove(reaction);
                      });
                    },
                  ),
                )),
            const SizedBox(height: 32),
            Center(
              child: ElevatedButton(
                onPressed: _submitAutomation,
                child: const Text("Submit"),
              ),
            ),
          ],
        ),
      ),
    ));
  }
}
