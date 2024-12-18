import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';
import 'package:area/pages/home_page.dart';
import 'package:area/tools/dynamic.dart';

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
        "label": "In...",
        "children": [
          {
            "type": "textfield",
            "extraParams": {
              "labelText": "Saisissez une valeur",
              "keyboardType": "number",
            },
          },
          {
            "type": "dropdown",
            "extraParams": {
              "initialValue": "Seconds",
              "items": [
                "Seconds",
                "Minutes",
                "Hours",
                "Days",
                "Weeks",
                "Months",
                "Years"
              ],
            },
          },
        ],
      },
    ];
    setState(() {});
  }

  void _loadMockReactions() {
    reactions = [
      {
        "label": "Send Message on Discord",
        "icon_url": "https://img.icons8.com/ios/452/discord.png",
      },
    ];
    reactionsConfigurations = [
      {
        "label": "Post message on a channel",
        "children": [
          {
            "type": "textfield",
            "extraParams": {
              "labelText": "Channel ID",
              "keyboardType": "number",
            },
          },
          {
            "type": "textfield",
            "extraParams": {
              "labelText": "Message on channel",
              "keyboardType": "text",
            },
          },
        ],
      },
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
          return config['children']?.map<Widget>((childConfig) {
                return Dynamic(
                  title: childConfig['type'],
                  extraParams:
                      Map<String, dynamic>.from(childConfig['extraParams']),
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
          return config['children']?.map<Widget>((childConfig) {
                return Dynamic(
                  title: childConfig['type'],
                  extraParams:
                      Map<String, dynamic>.from(childConfig['extraParams']),
                );
              }).toList() ??
              [];
        }).toList(),
      ),
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
              const MiniHeaderSection(),
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
                          color: selectedTrigger?['id'] == trigger['id']
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
                          color: selectedReaction?['id'] == reaction['id']
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
                onPressed: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text("Soumis")),
                  );
                },
                child: const Text("Soumettre"),
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
