import 'package:flutter/material.dart';

class TimeTrigger extends StatelessWidget {
  final Function(Map<String, int>) onTriggerChanged;

  const TimeTrigger({super.key, required this.onTriggerChanged});

  @override
  Widget build(BuildContext context) {
    Map<String, int> timeParams = {
      "Year": 0,
      "Month": 0,
      "Week": 0,
      "Day": 0,
      "Hour": 0,
      "Minute": 0,
      "Second": 0,
    };

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        ...timeParams.entries.map((entry) {
          return Padding(
            padding: const EdgeInsets.all(0),
            child: Row(
              children: [
                Text("${entry.key}:"),
                const SizedBox(width: 8),
                Expanded(
                  child: TextField(
                    decoration: InputDecoration(
                      labelText: "Enter ${entry.key}",
                      border: const OutlineInputBorder(),
                    ),
                    keyboardType: TextInputType.number,
                    onChanged: (value) {
                      timeParams[entry.key] = int.tryParse(value) ?? 0;
                      onTriggerChanged(timeParams);
                    },
                  ),
                ),
              ],
            ),
          );
        }),
      ],
    );
  }
}
