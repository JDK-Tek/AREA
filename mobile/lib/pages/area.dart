import 'package:area/pages/home_page.dart';
import 'package:area/tools/dynamic.dart';
import 'package:area/tools/providers.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:http/http.dart' as https;
import 'dart:convert';
import 'package:provider/provider.dart';

class CreateAutomationPage extends StatefulWidget {
  const CreateAutomationPage({super.key});

  @override
  CreateAutomationPageState createState() => CreateAutomationPageState();
}

class CreateAutomationPageState extends State<CreateAutomationPage> {
  List<Map<String, dynamic>> services = [];
  List<dynamic> selectedTriggers = [];
  List<dynamic> selectedReactions = [];
  List<Map<String, dynamic>> triggerValues = [];
  List<Map<String, dynamic>> reactionValues = [];

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
                    extraParams: {
                      'items': (spice['extra'] is List<String>)
                          ? spice['extra']
                          : (spice['extra']
                                  ?.map((e) => e.toString())
                                  .toList() ??
                              [])
                    },
                    onValueChanged: (key, value) {
                      setState(() {
                        if (spice['type'] == "number") {
                          tempValues[spice['name']] = int.tryParse(value) ?? 0;
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

  Color _colorFromHex(String hexColor) {
    if (hexColor.isEmpty) return Colors.grey;
    hexColor = hexColor.toUpperCase().replaceAll("#", "");
    if (hexColor.length == 6) {
      hexColor = "FF$hexColor";
    }
    return Color(int.tryParse(hexColor, radix: 16) ?? 0xFF000000);
  }

  void _addService(String type) {
    final items = services
        .where((service) => service[type] != null && service[type].isNotEmpty)
        .map((service) => {
              "name": service['name'],
              "icon": service['icon'],
              "color": service['color'],
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
              tileColor: _colorFromHex(item['color']),
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
              onTap: () async {
                final config = _buildDynamicConfig(
                    selectedTriggers, option['spices'], index);
                setState(() {
                  if (type == "actions") {
                    selectedTriggers.add({
                      "service": serviceName,
                      "name": option['name'],
                      "icon": services
                          .firstWhere((s) => s['name'] == serviceName)['icon'],
                      "color": services
                          .firstWhere((s) => s['name'] == serviceName)['color']
                    });
                    triggerValues.add(config);
                  } else {
                    selectedReactions.add({
                      "service": serviceName,
                      "name": option['name'],
                      "icon": services
                          .firstWhere((s) => s['name'] == serviceName)['icon'],
                      "color": services
                          .firstWhere((s) => s['name'] == serviceName)['color']
                    });
                    reactionValues.add(config);
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
      "action": selectedTriggers
          .map((trigger) {
            return {
              "service": trigger['service'],
              "name": trigger['name'],
              "spices": triggerValues[selectedTriggers.indexOf(trigger)],
            };
          })
          .toList()
          .first,
      "reaction": selectedReactions
          .map((reaction) {
            return {
              "service": reaction['service'],
              "name": reaction['name'],
              "spices": reactionValues[selectedReactions.indexOf(reaction)],
            };
          })
          .toList()
          .first,
    };
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
    final List<String> dest = ["/applets", "/create", "/services", "/plus"];
    return SafeArea(
      child: Scaffold(
        bottomNavigationBar: NavigationBar(
          labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
          backgroundColor: Colors.black,
          indicatorColor: Colors.grey,
          shadowColor: Colors.transparent,
          surfaceTintColor: Colors.transparent,
          selectedIndex: 1,
          onDestinationSelected: (int index) {
            setState(() {
              currentPageIndex = index;
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
        backgroundColor: Colors.white,
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
                    tileColor: _colorFromHex(trigger['color']),
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
                    tileColor: _colorFromHex(reaction['color']),
                    trailing: IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () {
                        setState(() {
                          reactionValues
                              .removeAt(selectedReactions.indexOf(reaction));
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
