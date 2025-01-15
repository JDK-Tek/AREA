import 'package:area/pages/home_page.dart';
import 'package:area/tools/dynamic.dart';
import 'package:area/tools/providers.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as https;
import 'dart:convert';
import 'package:provider/provider.dart';

class CreateAutomationPage extends StatefulWidget {
  const CreateAutomationPage({super.key});

  @override
  CreateAutomationPageState createState() => CreateAutomationPageState();
}

class CreateAutomationPageState extends State<CreateAutomationPage> {
  Map<String, dynamic> services = {};
  List<dynamic> selectedTriggers = [];
  List<dynamic> selectedReactions = [];
  List<Map<String, dynamic>> triggerValues = [];
  List<Map<String, dynamic>> reactionValues = [];

  @override
  void initState() {
    super.initState();
    _loadServices();
  }

  void _loadServices() {
    services = {
      "time": {
        "name": "time",
        "icon": "https://img.icons8.com/ios/452/timer.png",
        "color": "#ffffff",
        "actions": [
          {
            "name": "in",
            "description": "Triggers in some amount of time.",
            "spices": [
              {
                "name": "howmuch",
                "type": "number",
                "title": "How much time to wait.",
                "extraParams": null,
              },
              {
                "name": "unit",
                "type": "dropdown",
                "title": "The unit to wait.",
                "extraParams": ["weeks", "days", "hours", "minutes", "seconds"]
              }
            ]
          }
        ],
        "reactions": []
      },
      "discord": {
        "name": "discord",
        "icon":
            "https://cdn.prod.website-files.com/6257adef93867e50d84d30e2/636e0a6cc3c481a15a141738_icon_clyde_white_RGB.png",
        "color": "#5865F2",
        "actions": [],
        "reactions": [
          {
            "name": "send",
            "description": "Sends a message in a channel.",
            "spices": [
              {
                "name": "channel",
                "type": "text",
                "title": "The Discord channel ID.",
                "extraParams": null,
              },
              {
                "name": "message",
                "type": "text",
                "title": "The message to send.",
                "extraParams": null,
              }
            ]
          }
        ]
      }
    };
    setState(() {});
  }

  Map<String, dynamic> _buildDynamicConfig(
      List<dynamic> selectedItems, List<dynamic> spices, int index) {
    Map<String, dynamic> tempValues = {};
    showModalBottomSheet(
      context: context,
      builder: (BuildContext context) {
        return Padding(
          padding: const EdgeInsets.all(16.0),
          child: ListView(
            children: spices.map((spice) {
              return Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    spice['title'],
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 8),
                  Dynamic(
                    title: spice['type'],
                    extraParams: {'items': spice['extraParams']},
                    onValueChanged: (key, value) {
                      setState(() {
                        if (spice['type'] == "number") {
                          tempValues[spice['name']] = int.parse(value);
                        } else {
                          tempValues[spice['name']] = value;
                        }
                      });
                    },
                  ),
                  const SizedBox(height: 16),
                ],
              );
            }).toList(),
          ),
        );
      },
    );
    return tempValues;
  }

  void _addService(String type) {
    final items = services.values
        .where((service) => service[type].isNotEmpty)
        .map((service) => {
              "name": service['name'],
              "icon": service['icon'],
              "actionsOrReactions": service[type]
            })
        .toList();

    showModalBottomSheet(
      context: context,
      builder: (context) {
        return ListView.builder(
          itemCount: items.length,
          itemBuilder: (context, index) {
            final item = items[index];
            return ListTile(
              leading: Image.network(
                item['icon'],
                width: 40,
                height: 40,
              ),
              title: Text(item['name']),
              onTap: () {
                _selectService(type, item['actionsOrReactions'], item['name']);
              },
            );
          },
        );
      },
    );
  }

  void _selectService(String type, List<dynamic> options, String serviceName) {
    showModalBottomSheet(
      context: context,
      builder: (context) {
        return ListView.builder(
          itemCount: options.length,
          itemBuilder: (context, index) {
            final option = options[index];
            return ListTile(
              title: Text(option['name']),
              subtitle: Text(option['description']),
              onTap: () {
                setState(() {
                  if (type == "actions") {
                    selectedTriggers.add({
                      "service": serviceName,
                      "name": option['name'],
                      "icon": services[serviceName]['icon'],
                    });
                    triggerValues.add(_buildDynamicConfig(
                        selectedTriggers, option['spices'], index));
                  } else {
                    selectedReactions.add({
                      "service": serviceName,
                      "name": option['name'],
                      "icon": services[serviceName]['icon'],
                    });
                    reactionValues.add(_buildDynamicConfig(
                        selectedReactions, option['spices'], index));
                  }
                });
              },
            );
          },
        );
      },
    );
  }

  void _submitAutomation() {
    if (selectedTriggers.isEmpty || selectedReactions.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text("Please select at least one trigger and one reaction."),
        ),
      );
      return;
    }

    final automation = {
      "actions": selectedTriggers.map((trigger) {
        return {
          "service": trigger['service'],
          "name": trigger['name'],
          "spices": triggerValues[selectedTriggers.indexOf(trigger)],
        };
      }).toList(),
      "reactions": selectedReactions.map((reaction) {
        return {
          "service": reaction['service'],
          "name": reaction['name'],
          "spices": reactionValues[selectedReactions.indexOf(reaction)],
        };
      }).toList(),
    };

    print(jsonEncode(automation));
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

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Scaffold(
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
                      onPressed: () => _addService("actions"),
                    ),
                  ],
                ),
              ),
              ...selectedTriggers.map((trigger) => ListTile(
                    title: Text(trigger['name']),
                    leading: Image.network(trigger['icon']),
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
              const Divider(),
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
                      onPressed: () => _addService("reactions"),
                    ),
                  ],
                ),
              ),
              ...selectedReactions.map((reaction) => ListTile(
                    title: Text(reaction['name']),
                    leading: Image.network(reaction['icon']),
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
              Center(
                child: ElevatedButton(
                  onPressed: _submitAutomation,
                  child: const Text("Submit"),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
