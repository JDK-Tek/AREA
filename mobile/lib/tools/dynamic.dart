import 'package:flutter/material.dart';

class Dynamic extends StatefulWidget {
  const Dynamic({
    super.key,
    required this.title,
    this.extraParams,
    required this.onValueChanged,
  });

  final String title;
  final Map<String, dynamic>? extraParams;
  final Function(String key, String value) onValueChanged;

  @override
  State<Dynamic> createState() => _DynamicState();
}

class _DynamicState extends State<Dynamic> {
  String dropdownValue = 'Option 1';

  @override
  void initState() {
    super.initState();
    final extraParams = widget.extraParams ?? {};
    List<dynamic> items = extraParams['items'] ?? [];
    if (items.isNotEmpty) {
      dropdownValue = items[0];
    } else {
      dropdownValue = '';
    }
  }

  final FocusNode _focusNode = FocusNode();
  final TextEditingController textController = TextEditingController();
  @override
  void dispose() {
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final extraParams = widget.extraParams ?? {};

    final widgetMap = {
      "text": () {
        return SizedBox(
          width: 200,
          child: TextField(
            controller: textController,
            focusNode: _focusNode,
            decoration: InputDecoration(
              labelText: "Enter ${widget.title}",
              border: const OutlineInputBorder(),
            ),
            onChanged: (value) {
              widget.onValueChanged(widget.title, value);
            },
          ),
        );
      },
      "input": () {
        return SizedBox(
          width: 200,
          child: TextField(
            controller: textController,
            focusNode: _focusNode,
            decoration: InputDecoration(
              labelText: "Enter ${widget.title}",
              border: const OutlineInputBorder(),
            ),
            onChanged: (value) {
              widget.onValueChanged(widget.title, value);
            },
          ),
        );
      },
      "number": () {
        return SizedBox(
          width: 200,
          child: TextField(
              decoration: InputDecoration(
                labelText: "Enter ${widget.title}",
                border: const OutlineInputBorder(),
              ),
              keyboardType: TextInputType.number,
              onChanged: (value) {
                widget.onValueChanged(widget.title, value);
              }),
        );
      },
      "email": () {
        return SizedBox(
          width: 200,
          child: TextField(
            decoration: InputDecoration(
              labelText: "Enter ${widget.title}",
              border: const OutlineInputBorder(),
            ),
            keyboardType: TextInputType.emailAddress,
            onChanged: (value) {
              widget.onValueChanged(widget.title, value);
            },
          ),
        );
      },
      "dropdown": () {
        return DropdownButton<String>(
          value: dropdownValue,
          onChanged: (String? newValue) {
            setState(() {
              dropdownValue = newValue ?? dropdownValue;
            });
            widget.onValueChanged(widget.title, dropdownValue);
          },
          items: (extraParams['items'] as List<dynamic>)
              .map<DropdownMenuItem<String>>((dynamic value) {
            return DropdownMenuItem<String>(
              value: value.toString(),
              child: Text(value.toString()),
            );
          }).toList(),
        );
      },
      "date_picker": () {
        return ElevatedButton(
          onPressed: () async {
            DateTime? pickedDate = await showDatePicker(
              context: context,
              initialDate: DateTime.now(),
              firstDate: DateTime(2000),
              lastDate: DateTime(2100),
            );
            if (!context.mounted) return;
            if (pickedDate != null) {
              widget.onValueChanged(widget.title, pickedDate.toIso8601String());
            }
          },
          child: const Text("Select a date"),
        );
      },
      "phonenumber": () {
        return SizedBox(
          width: 200,
          child: TextField(
            decoration: InputDecoration(
              labelText: "Enter ${widget.title}",
              border: const OutlineInputBorder(),
            ),
            keyboardType: TextInputType.phone,
            onChanged: (value) {
              widget.onValueChanged(widget.title, value);
            },
          ),
        );
      },
      "url": () {
        return SizedBox(
          width: 200,
          child: TextField(
            decoration: InputDecoration(
              labelText: "Enter ${widget.title}",
              border: const OutlineInputBorder(),
            ),
            keyboardType: TextInputType.url,
            onChanged: (value) {
              widget.onValueChanged(widget.title, value);
            },
          ),
        );
      },
      "listview": () {
        List<String> items = extraParams['items'] ??
            List.generate(10, (index) => "Element $index");

        return SizedBox(
          height: 200,
          child: ListView.builder(
            physics: const NeverScrollableScrollPhysics(),
            itemCount: items.length,
            itemBuilder: (context, index) {
              return ListTile(
                title: Text(items[index]),
                onTap: () {
                  widget.onValueChanged(widget.title, items[index]);
                },
              );
            },
          ),
        );
      },
    };

    final widgetBuilder = widgetMap[widget.title];
    if (widgetBuilder != null) {
      return widgetBuilder();
    } else {
      return Text("Unsupported widget type: ${widget.title}");
    }
  }
}
