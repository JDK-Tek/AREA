import 'package:flutter/material.dart';

class DiscordAction extends StatelessWidget {
  final Function(String, String) onActionChanged;

  const DiscordAction({super.key, required this.onActionChanged});

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
            onActionChanged(value, "");
          },
        ),
        const SizedBox(height: 8),
        TextField(
          decoration: const InputDecoration(
            labelText: "Message Template (use {time} for time)",
            border: OutlineInputBorder(),
          ),
          onChanged: (value) {
            onActionChanged("", value);
          },
        ),
      ],
    );
  }
}