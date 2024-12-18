import 'package:flutter/material.dart';

class Dynamic extends StatefulWidget {
  const Dynamic({super.key, required this.title, this.extraParams});
  final String title;
  final Map<String, dynamic>? extraParams;

  @override
  State<Dynamic> createState() => _DynamicState();
}

class _DynamicState extends State<Dynamic> {
  String dropdownValue = 'Option 1';
  bool switchValue = false;

  @override
  Widget build(BuildContext context) {
    final extraParams = widget.extraParams ?? {};

    if (widget.title == "textfield") {
      TextInputType key = ((extraParams["keyboardType"] ?? "text") == "number"
          ? TextInputType.number
          : TextInputType.text);
      return SizedBox(
        width: 200,
          child: TextField(
        decoration: InputDecoration(
          labelText: extraParams['labelText'] ?? "Entrez du texte",
          border: const OutlineInputBorder(),
        ),
        keyboardType: key,
      ));
    } else if (widget.title == "timestamp") {
      return TimePickerDialog(
          initialTime: extraParams['initialTime'] ?? TimeOfDay.now());
    } else if (widget.title == "button") {
      return ElevatedButton(
        onPressed: extraParams['onPressed'] ??
            () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text("Bouton pressé")),
              );
            },
        child: Text(extraParams['buttonText'] ?? "Cliquez Moi"),
      );
    } else if (widget.title == "switch") {
      return SwitchListTile(
        title: Text(extraParams['title'] ?? "Activer/Désactiver"),
        value: switchValue,
        onChanged: (bool value) {
          setState(() {
            switchValue = value;
          });
        },
      );
    } else if (widget.title == "dropdown") {
      dropdownValue = extraParams['initialValue'] ?? dropdownValue;
      List<DropdownMenuItem<String>> items = (extraParams['items']
                  as List<String>?)
              ?.map((item) => DropdownMenuItem(value: item, child: Text(item)))
              .toList() ??
          const [
            DropdownMenuItem(value: 'Option 1', child: Text('Option 1')),
            DropdownMenuItem(value: 'Option 2', child: Text('Option 2')),
            DropdownMenuItem(value: 'Option 3', child: Text('Option 3')),
          ];
      return DropdownButton<String>(
        value: dropdownValue,
        items: items,
        onChanged: (String? newValue) {
          setState(() {
            if (newValue != null) dropdownValue = newValue;
          });
        },
      );
    } else if (widget.title == "slider") {
      double sliderValue = extraParams['initialValue'] ?? 0.5;
      return StatefulBuilder(
        builder: (context, setState) => Slider(
          value: sliderValue,
          min: extraParams['min'] ?? 0.0,
          max: extraParams['max'] ?? 1.0,
          divisions: extraParams['divisions'] ?? 10,
          label: extraParams['label']?.call(sliderValue) ??
              "${(sliderValue * 100).toStringAsFixed(0)}%",
          onChanged: (value) {
            setState(() {
              sliderValue = value;
            });
          },
        ),
      );
    } else if (widget.title == "date_picker") {
      return ElevatedButton(
        onPressed: () async {
          DateTime? pickedDate = await showDatePicker(
            context: context,
            initialDate: extraParams['initialDate'] ?? DateTime.now(),
            firstDate: extraParams['firstDate'] ?? DateTime(2000),
            lastDate: extraParams['lastDate'] ?? DateTime(2100),
          );
          if (!context.mounted) return;
          if (pickedDate != null) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                  content: Text("Date sélectionnée: ${pickedDate.toLocal()}")),
            );
          }
        },
        child: Text(extraParams['buttonText'] ?? "Choisir une date"),
      );
    } else if (widget.title == "checkbox") {
      bool checkboxValue = extraParams['initialValue'] ?? false;
      return StatefulBuilder(
        builder: (context, setState) => CheckboxListTile(
          title: Text(extraParams['title'] ?? "Activer l'option"),
          value: checkboxValue,
          onChanged: (bool? value) {
            if (value != null) {
              setState(() {
                checkboxValue = value;
              });
            }
          },
        ),
      );
    } else if (widget.title == "listview") {
      List<String> items = extraParams['items'] ??
          List.generate(10, (index) => "Élément $index");
      return SizedBox(
        height: 200,
        child: ListView.builder(
          itemCount: items.length,
          itemBuilder: (context, index) {
            return ListTile(
              leading: const Icon(Icons.star),
              title: Text(items[index]),
              onTap: () {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                      content: Text("Vous avez cliqué sur ${items[index]}")),
                );
              },
            );
          },
        ),
      );
    } else if (widget.title == "dialog") {
      return ElevatedButton(
        onPressed: () {
          showDialog(
            context: context,
            builder: (context) => AlertDialog(
              title: Text(extraParams['title'] ?? "Exemple de Dialog"),
              content: Text(extraParams['content'] ??
                  "Ceci est une boîte de dialogue simple"),
              actions: [
                TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: Text(extraParams['cancelText'] ?? "Annuler"),
                ),
                TextButton(
                  onPressed: () {
                    Navigator.of(context).pop();
                    if (extraParams['onConfirm'] != null) {
                      extraParams['onConfirm']();
                    } else {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text("Action confirmée")),
                      );
                    }
                  },
                  child: Text(extraParams['confirmText'] ?? "Confirmer"),
                ),
              ],
            ),
          );
        },
        child: Text(extraParams['buttonText'] ?? "Afficher le Dialog"),
      );
    } else {
      return Text(widget.title);
    }
  }
}
