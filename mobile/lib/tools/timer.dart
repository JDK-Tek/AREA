import 'package:flutter/material.dart';
import 'package:area/tools/action_reaction.dart';

class TimeTrigger extends StatefulWidget implements ActionHandler {
  final Function(Map<String, int>) onTriggerChanged;

  const TimeTrigger({super.key, required this.onTriggerChanged});

  /// fonction implementé plus tard
  @override
  Map<String, dynamic> toJson() {
    return {
      "service": "time",
      "name": "in",
      "spices": {"howmuch": 0, "unit": "..."}
    };
  }
  @override
  _TimeTriggerState createState() => _TimeTriggerState();
}

class _TimeTriggerState extends State<TimeTrigger> implements ActionHandler {
  Map<String, int> timeParams = {
    "Year": 0,
    "Month": 0,
    "Week": 0,
    "Day": 0,
    "Hour": 0,
    "Minute": 0,
    "Second": 0,
  };

  String selectedUnit = "Year";

  @override
  Map<String, dynamic> toJson() {
    return {
      "service": "time",
      "name": "in",
      "spices": {"howmuch": timeParams[selectedUnit], "unit": selectedUnit}
    };
  }

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
                    widget.onTriggerChanged(timeParams);
                  });
                }
              },
            ),
            const SizedBox(width: 8),
            Expanded(
              child: TextField(
                keyboardType: TextInputType.number,
                decoration: InputDecoration(
                  labelText: "Enter $selectedUnit",
                  border: const OutlineInputBorder(),
                ),
                onChanged: (value) {
                  timeParams[selectedUnit] = int.tryParse(value) ?? 0;
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
