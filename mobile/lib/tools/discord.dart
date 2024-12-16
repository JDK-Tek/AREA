import 'package:flutter/material.dart';
import 'package:area/tools/action_reaction.dart';

class DiscordAction extends StatelessWidget implements ReactionHandler {
  final Function(String, String) onActionChanged;
  DiscordAction({super.key, required this.onActionChanged});

  String message = "";
  String channelId = "";

  @override
  Map<String, dynamic> toJson() {
    return {
      "service": "discord",
      "name": "send",
      "spices": {
        "channel": int.tryParse(channelId) ?? 0,
        "message": message.isNotEmpty ? message : "Default Message"
      }
    };
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TextField(
          decoration: const InputDecoration(
            labelText: "Discord Channel ID",
            border: OutlineInputBorder(),
          ),
          onChanged: (value) {
            channelId = value;
            onActionChanged(channelId, "");
          },
        ),
        const SizedBox(height: 8),
        TextField(
          decoration: const InputDecoration(
            labelText: "Message Template (use {time} for time)",
            border: OutlineInputBorder(),
          ),
          onChanged: (value) {
            message = value;
            onActionChanged("", message);
          },
        ),
      ],
    );
  }
}
