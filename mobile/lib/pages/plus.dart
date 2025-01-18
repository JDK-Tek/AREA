import 'package:area/tools/dynamic.dart';
import 'package:flutter/material.dart';
import 'package:area/pages/home_page.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';
import 'package:area/tools/providers.dart';
import 'package:provider/provider.dart';

class PlusPage extends StatefulWidget {
  const PlusPage({super.key});

  @override
  State<PlusPage> createState() => PlusPageState();
}

class PlusPageState extends State<PlusPage> {
  TextEditingController textcontroller = TextEditingController();
  final FocusNode _focusNode = FocusNode();
  int currentPageIndex = 2;

  @override
  void dispose() {
    _focusNode.dispose();
    super.dispose();
  }

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
          selectedIndex: 3,
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
              SizedBox(
                height: MediaQuery.of(context).size.width <
                        MediaQuery.of(context).size.height
                    ? MediaQuery.of(context).size.height * 0.3
                    : MediaQuery.of(context).size.height * 0.75,
                width: MediaQuery.of(context).size.width,
                child: Dynamic(
                  title: "listview",
                  extraParams: const {
                    "items": [
                      "Profile",
                      "About Us",
                      "Terms Of Services",
                      "Privacy Policy"
                    ],
                  },
                  onValueChanged: (key, value) {
                    setState(() {
                      if (value == "Profile") {
                        context.go("/plus");
                      }
                      if (value == "About Us") {
                        context.go("/developers");
                      }
                      if (value == "Terms Of Services") {
                        context.go("/termsofservices");
                      }
                      if (value == "Privacy Policy") {
                        context.go("/privacypolicy");
                      }
                    });
                  },
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      "IP Address:",
                      style: TextStyle(fontWeight: FontWeight.w500),
                    ),
                    Row(
                      children: [
                        SizedBox(
                          height: MediaQuery.of(context).size.width <
                                  MediaQuery.of(context).size.height
                              ? MediaQuery.of(context).size.height * 0.05
                              : MediaQuery.of(context).size.height * 0.1,
                          width: MediaQuery.of(context).size.width * 0.75,
                          child: TextField(
                            key: const ValueKey("uniqueTextFieldKey"),
                            focusNode: _focusNode,
                            controller: textcontroller,
                            decoration: const InputDecoration(
                              labelText: "Enter your new IP address server...",
                              border: OutlineInputBorder(),
                            ),
                            onTap: () {
                              if (_focusNode.hasFocus) {
                                FocusScope.of(context).unfocus();
                                _focusNode.requestFocus();
                              }
                            },
                          ),
                        ),
                        ElevatedButton(
                            style: ButtonStyle(
                              backgroundColor:
                                  WidgetStateProperty.all(Colors.transparent),
                              shadowColor:
                                  WidgetStateProperty.all(Colors.transparent),
                              foregroundColor:
                                  WidgetStateProperty.all(Colors.transparent),
                            ),
                            onPressed: () {
                              final String ip = textcontroller.text.trim();
                              Provider.of<IPState>(context, listen: false)
                                  .setIP(ip);
                            },
                            child: const Icon(Icons.check_box,
                                color: Colors.green)),
                      ],
                    ),
                  ],
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
