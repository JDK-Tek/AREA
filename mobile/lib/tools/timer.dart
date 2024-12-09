import 'package:flutter/material.dart';

class TimeTrigger extends StatefulWidget {
  final Function(Map<String, int>) onTriggerChanged;

  const TimeTrigger({super.key, required this.onTriggerChanged});

  @override
  State<TimeTrigger> createState() => _TimeTriggerState();
}

class _TimeTriggerState extends State<TimeTrigger> {
  final Map<String, int> timeParams = {
    "Year": 0,
    "Month": 0,
    "Week": 0,
    "Day": 0,
    "Hour": 0,
    "Minute": 0,
    "Second": 0,
  };

  String selectedUnit = "Year"; // L'unité par défaut

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            DropdownButton<String>(
              value: selectedUnit,
              items: timeParams.keys.map((unit) {
                return DropdownMenuItem<String>(
                  value: unit,
                  child: Text(unit),
                );
              }).toList(),
              onChanged: (value) {
                if (value != null) {
                  setState(() {
                    selectedUnit = value;
                  });
                }
              },
            ),
            const SizedBox(width: 16),
            Expanded(
              child: TextField(
                decoration: InputDecoration(
                  labelText: "Enter $selectedUnit",
                  border: const OutlineInputBorder(),
                ),
                keyboardType: TextInputType.number,
                onChanged: (value) {
                  setState(() {
                    timeParams[selectedUnit] = int.tryParse(value) ?? 0;
                  });
                  widget.onTriggerChanged(timeParams);
                },
              ),
            ),
          ],
        ),
      ],
    );
  }
}
