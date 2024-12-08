import 'package:area/pages/home_page.dart';
import 'package:area/tools/footerarea.dart';
import 'package:flutter/material.dart';
import 'package:area/tools/timer.dart';
import 'package:area/tools/discord.dart';
import 'package:footer/footer_view.dart';

class CreateAutomationPage extends StatefulWidget {
  const CreateAutomationPage({super.key});

  @override
  CreateAutomationPageState createState() => CreateAutomationPageState();
}

class CreateAutomationPageState extends State<CreateAutomationPage> {
  String? selectedAction;
  String? selectedReaction;

  Map<String, int> timeTriggerData = {};
  String discordChannelId = "";
  String discordMessageTemplate = "";

  final List<String> triggers = ["Time"];
  final List<IconData> triggerIcons = [Icons.access_time];
  final List<Widget Function(Function(Map<String, int>) callback)>
      triggersBuilder = [
    (callback) => TimeTrigger(onTriggerChanged: callback),
  ];

  final List<String> actions = ["Send a Discord message"];
  final List<IconData> actionIcons = [Icons.discord];
  final List<Widget Function(Function(String, String) callback)>
      actionsBuilder = [
    (callback) => DiscordAction(onActionChanged: callback),
  ];

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
            Icon(icon,
                color: isSelected ? Colors.blueAccent : Colors.grey[600],
                size: 32),
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

  Widget _buildSummary(double screenHeight) {
    return SizedBox(
      height: screenHeight * 0.2,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text("Summary",
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
          const SizedBox(height: 8),
          Text("Action: $selectedAction"),
          if (selectedAction == "Time")
            ...timeTriggerData.entries.map((entry) {
              return Text("${entry.key}: ${entry.value}");
            }),
          const SizedBox(height: 8),
          Text("Reaction: $selectedReaction"),
          if (selectedReaction == "Send a Discord message")
            Text(
                "Channel ID: $discordChannelId\nMessage: $discordMessageTemplate"),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final screenHeight = MediaQuery.of(context).size.height;
    final screenWidth = MediaQuery.of(context).size.width;

    return SafeArea(
      child: Scaffold(
        body: FooterView(footer: const Footerarea().build(context), children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const HeaderSection(),
              const Text("Choose a Reaction",
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
              SizedBox(
                height: screenHeight * 0.35,
                child: GridView.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    crossAxisSpacing: 0,
                    mainAxisSpacing: 0,
                  ),
                  itemCount: triggers.length,
                  itemBuilder: (context, index) {
                    return _buildGridItem(
                      label: triggers[index],
                      icon: triggerIcons[index],
                      isSelected: selectedAction == triggers[index],
                      onTap: () {
                        setState(() {
                          selectedAction = triggers[index];
                        });
                      },
                    );
                  },
                ),
              ),
              if (selectedAction != null)
                triggersBuilder[triggers.indexOf(selectedAction!)](
                  (data) {
                    setState(() {
                      timeTriggerData = data;
                    });
                  },
                ),
              SizedBox(height: screenHeight * 0.02),
              const Text("Choose an Action",
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
              SizedBox(
                height: screenHeight * 0.35,
                child: GridView.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    crossAxisSpacing: 0,
                    mainAxisSpacing: 0,
                  ),
                  itemCount: actions.length,
                  itemBuilder: (context, index) {
                    return _buildGridItem(
                      label: actions[index],
                      icon: actionIcons[index],
                      isSelected: selectedReaction == actions[index],
                      onTap: () {
                        setState(() {
                          selectedReaction = actions[index];
                        });
                      },
                    );
                  },
                ),
              ),
              if (selectedReaction != null)
                actionsBuilder[actions.indexOf(selectedReaction!)](
                  (channelId, messageTemplate) {
                    setState(() {
                      discordChannelId = channelId;
                      discordMessageTemplate = messageTemplate;
                    });
                  },
                ),
              SizedBox(height: screenHeight * 0.02),
              if (selectedAction != null && selectedReaction != null)
                _buildSummary(screenHeight),
              SizedBox(height: screenHeight * 0.02),
              if (selectedAction != null && selectedReaction != null)
                Center(
                  child: SizedBox(
                    width: screenWidth * 0.5,
                    height: 50,
                    child: ElevatedButton(
                      onPressed: () {
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(
                              content:
                                  Text("Automation created successfully!")),
                        );
                      },
                      child: const Text("Create Automation"),
                    ),
                  ),
                ),
              SizedBox(height: screenHeight * 0.02),
            ],
          ),
        ]),
      ),
    );
  }
}
